package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"
	model "dropshipbe/model/schema"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CaptureOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptureOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptureOrderLogic {
	return &CaptureOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *CaptureOrderLogic) CaptureOrder(in *dropshipbe.CaptureOrderRequest) (*dropshipbe.CaptureOrderResponse, error) {
	// ==========================================
	// 1. FIND TRANSACTION AND ORDER
	// ==========================================
	var transaction model.Transaction

	// Preload Order and OrderItems (crucial for inventory rollback if capture fails)
	err := l.svcCtx.DB.Preload("Order.OrderItems").
		Where("transaction_reference = ?", in.PaypalOrderId).
		First(&transaction).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid paypal order id: transaction not found")
		}
		return nil, fmt.Errorf("system error retrieving transaction: %v", err)
	}

	order := transaction.Order
	if order == nil {
		return nil, fmt.Errorf("data corruption: order details missing")
	}

	// ==========================================
	// 2. IDEMPOTENCY CHECK (ANTI-SPAM / DOUBLE CHARGE)
	// ==========================================
	if transaction.Status == "completed" || order.FinancialStatus == "paid" {
		return &dropshipbe.CaptureOrderResponse{
			Success: true,
			Status:  "COMPLETED",
			Message: "This order has already been successfully paid.",
		}, nil
	}

	// ==========================================
	// 3. CALL PAYPAL API TO CAPTURE FUNDS
	// ==========================================
	l.Logger.Infof("Processing capture for PayPal Order ID: %s", in.PaypalOrderId)

	status, rawResponse, err := utils.CapturePayPalOrder(
		l.svcCtx.Config.PayPal.PaypalBaseURL,
		l.svcCtx.Config.PayPal.ClientID,
		l.svcCtx.Config.PayPal.Secret,
		in.PaypalOrderId,
	)

	// Marshal rawResponse to byte array for Postgres JSONB storage
	rawJSON, _ := json.Marshal(rawResponse)

	// ==========================================
	// 4. HANDLE PAYMENT FAILURE OR DECLINE
	// ==========================================
	if err != nil || status != "COMPLETED" {
		l.Logger.Errorf("Capture failed. Status: %s, Error: %v", status, err)

		l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			tx.Model(&transaction).Updates(map[string]interface{}{
				"status":       "failed",
				"raw_response": datatypes.JSON(rawJSON),
			})
			tx.Model(order).Update("financial_status", "failed")
			return nil
		})

		l.rollbackInventory(order.OrderItems)

		return &dropshipbe.CaptureOrderResponse{
			Success: false,
			Status:  status,
			Message: "Transaction declined. Please check your payment method or PayPal account.",
		}, nil
	}

	// ==========================================
	// 5. EXTRACT CUSTOMER INFO FROM PAYPAL
	// ==========================================
	// Bóc tách Email, Tên, và Địa chỉ giao hàng từ kết quả trả về của PayPal
	paypalEmail := ""
	if payer, ok := rawResponse["payer"].(map[string]interface{}); ok {
		if email, ok := payer["email_address"].(string); ok {
			paypalEmail = email
		}
	}

	var shippingName, shippingAddressStr string
	if purchaseUnits, ok := rawResponse["purchase_units"].([]interface{}); ok && len(purchaseUnits) > 0 {
		if firstUnit, ok := purchaseUnits[0].(map[string]interface{}); ok {
			if shipping, ok := firstUnit["shipping"].(map[string]interface{}); ok {
				// Lấy Tên
				if nameObj, ok := shipping["name"].(map[string]interface{}); ok {
					shippingName, _ = nameObj["full_name"].(string)
				}
				// Lấy Địa chỉ
				if addr, ok := shipping["address"].(map[string]interface{}); ok {
					shippingAddressStr = fmt.Sprintf("%v, %v, %v, %v, %v",
						addr["address_line_1"],
						addr["admin_area_2"], // City
						addr["admin_area_1"], // State/Province
						addr["postal_code"],
						addr["country_code"],
					)
				}
			}
		}
	}

	// Đóng gói lại thành JSON an toàn cho DB
	shippingInfo := map[string]string{
		"recipient_name": shippingName,
		"full_address":   shippingAddressStr,
		"email":          paypalEmail,
		"source":         "paypal_capture",
	}
	shippingJSON, _ := json.Marshal(shippingInfo)

	// ==========================================
	// 6. HANDLE PAYMENT SUCCESS & UPDATE DB
	// ==========================================
	l.Logger.Infof("💰 Successfully captured payment for order: %s", order.OrderNumber)

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 6.1 Update Transaction
		if err := tx.Model(&transaction).Updates(map[string]interface{}{
			"status":       "completed",
			"raw_response": datatypes.JSON(rawJSON),
		}).Error; err != nil {
			return err
		}

		// 6.2 Update Order Status AND Customer Info
		if err := tx.Model(order).Updates(map[string]interface{}{
			"financial_status": "paid",
			"customer_email":   paypalEmail,
			"shipping_address": datatypes.JSON(shippingJSON),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		l.Logger.Errorf("CRITICAL DB ERROR: Payment collected but DB update failed for order %s: %v", order.OrderNumber, err)
	}

	// ==========================================
	// 7. PUBLISH KAFKA EVENT (ORDER_PAID)
	// ==========================================
	eventPayload := map[string]interface{}{
		"event_id":       fmt.Sprintf("evt_%s_%d", order.OrderNumber, time.Now().UnixNano()),
		"event_type":     "ORDER_PAID",
		"order_id":       order.ID,
		"order_number":   order.OrderNumber,
		"customer_email": paypalEmail,
		"customer_name":  shippingName,
		"total_amount":   order.TotalPrice,
		"currency":       order.Currency,
		"timestamp":      time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(eventPayload)
	if err != nil {
		l.Logger.Errorf("Failed to Marshal Kafka Event for order %s: %v", order.OrderNumber, err)
	} else {
		if l.svcCtx.KqNotificationPusherClient != nil {
			err = l.svcCtx.KqNotificationPusherClient.Push(l.ctx, string(msgBytes))
			if err != nil {
				l.Logger.Errorf("Failed to push Kafka event for order %s: %v", order.OrderNumber, err)
			} else {
				l.Logger.Infof("🚀 Successfully published Kafka Event [ORDER_PAID] for order: %s", order.OrderNumber)
			}
		} else {
			l.Logger.Errorf("Kafka Client is not initialized in ServiceContext. Skipped sending event for order: %s", order.OrderNumber)
		}
	}

	// ==========================================
	// 8. RETURN SUCCESS TO FRONTEND
	// ==========================================
	return &dropshipbe.CaptureOrderResponse{
		Success: true,
		Status:  "COMPLETED",
		Message: "Payment successful. Your order is being prepared.",
	}, nil
}

// Helper function to release inventory upon payment failure
func (l *CaptureOrderLogic) rollbackInventory(items []model.OrderItem) {
	for _, item := range items {

		err := utils.RollbackInventory(context.Background(), l.svcCtx.Redis, item.VariantID, int32(item.Quantity))
		if err != nil {
			l.Logger.Errorf("CRITICAL: Failed to rollback inventory for VariantID %d: %v", item.VariantID, err)
		}
	}
}
