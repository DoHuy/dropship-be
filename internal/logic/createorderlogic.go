package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"dropshipbe/common/constant"
	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Orders & Checkout ---
func (l *CreateOrderLogic) CreateOrder(in *dropshipbe.CreateOrderRequest) (*dropshipbe.CreateOrderResponse, error) {
	// todo: add your logic here and delete this line

	var variantIDs []uint64
	for _, item := range in.Items {
		variantIDs = append(variantIDs, item.VariantId)
	}

	// (Giả sử bạn có bảng variants chứa giá)
	var variants []model.Variant
	if err := l.svcCtx.DB.Preload("Product").Where("id IN ?", variantIDs).Find(&variants).Error; err != nil {
		l.Logger.Errorf("Lỗi truy vấn DB: %v", err)
		return nil, fmt.Errorf("lỗi hệ thống khi kiểm tra sản phẩm")
	}

	// 1. Khởi tạo Map để lưu trữ VariantID -> Price
	variantPriceMap := make(map[uint64]model.Variant)

	// 2. Đổ dữ liệu từ mảng variants vào Map
	for _, v := range variants {
		variantPriceMap[v.ID] = v
	}

	// 3. Tính tổng tiền thực tế (Anti-Hack) và chuẩn bị mảng OrderItem
	var totalAmount float64
	var orderItems []model.OrderItem

	for _, item := range in.Items {
		v, exists := variantPriceMap[item.VariantId]
		if !exists || v.IsActive == nil || !*v.IsActive {
			return nil, fmt.Errorf("sản phẩm %d không tồn tại hoặc ngừng kinh doanh", item.VariantId)
		}

		price := v.Price
		quantity := int(item.Quantity)
		lineTotal := price * float64(quantity)
		totalAmount += lineTotal

		// Tạo Snapshot cho dòng sản phẩm này
		orderItems = append(orderItems, model.OrderItem{
			VariantID:   v.ID,
			ProductID:   v.ProductID,
			ProductName: v.Product.Name,
			VariantName: v.Sku,
			Sku:         v.Sku,
			Quantity:    quantity,
			Price:       price,
			Total:       lineTotal,
		})
	}
	l.Logger.Infof("Tổng tiền đơn hàng: %.2f USD", totalAmount)

	// ==========================================
	// BƯỚC 3: TRỪ TỒN KHO REDIS BẰNG LUA SCRIPT
	// ==========================================

	var successfullyDeducted []*dropshipbe.CartItem

	for _, item := range in.Items {
		res, err := utils.DeductInventory(l.ctx, l.svcCtx.Redis, item.VariantId, item.Quantity)
		if err != nil || res != 1 {
			// Nếu có 1 món hết hàng -> Lập tức ROLLBACK những món đã trừ trước đó
			l.Logger.Errorf("Sản phẩm %s hết hàng hoặc lỗi kho", item.VariantId)
			utils.RollbackInventory(l.ctx, l.svcCtx.Redis, item.VariantId, item.Quantity)
			return nil, fmt.Errorf("sản phẩm %d hiện không đủ số lượng trong kho", item.VariantId)
		}
		// Lưu lại danh sách đã trừ thành công để lỡ có lỗi phía sau còn biết đường Rollback
		successfullyDeducted = append(successfullyDeducted, item)
	}

	// ==========================================
	// BƯỚC 4: LƯU DATABASE BẰNG TRANSACTION
	// ==========================================

	// ==========================================
	// BƯỚC 4: XỬ LÝ OPTIONAL FIELDS & JSON ADDRESS
	// ==========================================
	// Bóc tách các giá trị từ con trỏ Optional
	email, phone, name, addressStr := "", "", "", ""
	if in.CustomerEmail != nil {
		email = *in.CustomerEmail
	}
	if in.CustomerPhone != nil {
		phone = *in.CustomerPhone
	}
	if in.CustomerName != nil {
		name = *in.CustomerName
	}
	if in.ShippingAddress != nil {
		addressStr = *in.ShippingAddress
	}

	// Gói chuỗi địa chỉ vào định dạng JSON an toàn cho cột JSONB của Postgres
	addressMap := map[string]string{
		"recipient_name": name,
		"full_address":   addressStr,
		"method":         in.ShippingMethod,
	}
	addressBytes, _ := json.Marshal(addressMap)

	// Tạo Order Code dạng: ORD-1681234567
	orderCode := constant.CreateOrderNumber()

	newOrder := model.Order{
		OrderNumber:       orderCode,
		CustomerEmail:     email,
		CustomerPhone:     phone,
		ShippingAddress:   addressBytes,
		TotalPrice:        totalAmount,
		SubtotalPrice:     totalAmount,
		Currency:          "USD",
		FinancialStatus:   "pending",
		FulfillmentStatus: "unfulfilled",
	}

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 4.1 Lưu bảng Order
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		// 4.2 Gắn OrderID vào từng OrderItem và lưu (Bulk Insert)
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}

		if len(orderItems) > 0 {
			if err := tx.Create(&orderItems).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		// NẾU DB SẬP -> Phải Rollback kho Redis
		l.Logger.Errorf("Lỗi Transaction DB: %v", err)
		l.rollbackAllInventory(successfullyDeducted)
		return nil, fmt.Errorf("lỗi hệ thống khi tạo đơn hàng, vui lòng thử lại")
	}

	// ==========================================
	// BƯỚC 5: GỌI PAYPAL TẠO PHIÊN THANH TOÁN
	// ==========================================
	l.Logger.Infof("Gọi PayPal cho Order: %s, Số tiền: %.2f", orderCode, totalAmount)

	paypalOrderID, _, err := utils.CreatePayPalOrder(
		l.svcCtx.Config.PayPal.PaypalBaseURL,
		l.svcCtx.Config.PayPal.ClientID,
		l.svcCtx.Config.PayPal.Secret,
		l.svcCtx.Config.PayPal.Mode,
		totalAmount,
	)

	if err != nil {
		// NẾU PAYPAL LỖI -> Đánh dấu DB là CANCELED và Rollback kho Redis
		l.Logger.Errorf("Lỗi gọi PayPal: %v", err)
		l.svcCtx.DB.Model(&newOrder).Update("financial_status", "canceled")
		l.rollbackAllInventory(successfullyDeducted)
		return nil, fmt.Errorf("cổng thanh toán đang bảo trì, vui lòng thử lại sau")
	}

	// ==========================================
	// BƯỚC 7: LƯU VẾT TRANSACTION ĐỂ PHỤC VỤ API CAPTURE
	// ==========================================
	newTransaction := model.Transaction{
		OrderID:              newOrder.ID,
		Gateway:              "paypal",
		PaymentMethod:        "paypal_checkout",
		TransactionReference: paypalOrderID,
		Amount:               totalAmount,
		Currency:             "USD",
		Status:               "pending",
	}

	// Dùng Create bình thường (không cần bọc Transaction riêng)
	if err := l.svcCtx.DB.Create(&newTransaction).Error; err != nil {
		l.Logger.Errorf("Lưu Transaction thất bại cho Order %s: %v", orderCode, err)

		// [SỬA LỖI]: Bắt buộc phải Rollback nếu không lưu được vết
		// 1. Đánh dấu đơn hàng là Canceled
		l.svcCtx.DB.Model(&newOrder).Update("financial_status", "canceled")

		// 2. Nhả lại tồn kho Redis cho khách khác mua
		l.rollbackAllInventory(successfullyDeducted)

		return nil, fmt.Errorf("lỗi hệ thống khi khởi tạo phiên thanh toán, vui lòng thử lại")
	}

	// ==========================================
	// BƯỚC 6: TRẢ KẾT QUẢ VỀ FRONTEND
	// ==========================================
	return &dropshipbe.CreateOrderResponse{
		LocalOrderId:  orderCode,
		PaypalOrderId: paypalOrderID,
	}, nil

}

func (l *CreateOrderLogic) rollbackAllInventory(items []*dropshipbe.CartItem) {
	for _, item := range items {
		err := utils.RollbackInventory(context.Background(), l.svcCtx.Redis, item.VariantId, item.Quantity)
		if err != nil {
			l.Logger.Errorf("CRITICAL: Lỗi Rollback kho cho VariantID %s: %v", item.VariantId, err)
			// Trong thực tế, nếu lỗi ở đây cần bắn Alert ra Telegram/Slack cho Admin xử lý
		}
	}
}
