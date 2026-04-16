package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"dropshipbe/common/kafka_events"
	"dropshipbe/common/notify"
	"dropshipbe/mq/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationConsumer struct {
	Config    config.Config
	ctx       context.Context
	sesClient *notify.SESClient
	logx.Logger
}

func NewNotificationConsumer(ctx context.Context, c config.Config) *NotificationConsumer {
	// Khởi tạo AWS SES Client
	// ses, err := notify.NewSESClient(
	// 	context.Background(),
	// 	c.Email.Region,
	// 	c.Email.AccessKey,
	// 	c.Email.SecretKey,
	// 	c.Email.FromAddress,
	// )

	// if err != nil {
	// 	// Dùng logx.Must để nếu cấu hình AWS sai, app sẽ crash ngay lập tức lúc khởi động
	// 	// thay vì im lặng chạy rồi lỗi lúc gửi mail.
	// 	logx.Must(err)
	// }

	return &NotificationConsumer{
		ctx:       ctx,
		Config:    c,
		sesClient: nil, // chưa có key nên tạm thời đóng
		Logger:    logx.WithContext(ctx),
	}
}

// Consume là hàm bắt buộc phải có để implement kq.ConsumeHandler interface
// Kafka sẽ tự động gọi hàm này mỗi khi có event mới chui vào Topic
func (c *NotificationConsumer) Consume(ctx context.Context, key, val string) error {
	c.Logger.Infof("📥 Nhận được event từ Kafka: %s", val)

	// 1. Ép chuỗi JSON thành Struct
	var payload kafka_events.OrderEventPayload
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		c.Logger.Errorf("❌ Lỗi parse JSON event: %v", err)
		// Return nil để báo Kafka "bỏ qua" event lỗi định dạng này, tránh kẹt hàng đợi (Dead-letter)
		return nil
	}

	// 2. Phân loại Event để xử lý
	switch payload.EventType {
	case "ORDER_PAID":
		return c.handleOrderPaid(ctx, payload)
	default:
		c.Logger.Infof("⚠️ Bỏ qua event không thuộc chức năng thông báo: %s", payload.EventType)
		return nil
	}
}

// handleOrderPaid chứa logic nghiệp vụ chính khi có đơn hàng được thanh toán
func (c *NotificationConsumer) handleOrderPaid(ctx context.Context, payload kafka_events.OrderEventPayload) error {

	// ==========================================
	// 1. GỬI TELEGRAM CHO ADMIN
	// ==========================================
	teleMsg := fmt.Sprintf(
		"💸 <b>ĐƠN HÀNG MỚI ĐÃ THANH TOÁN</b>\n"+
			"Mã đơn: %s\nKhách hàng: %s\nDoanh thu: %.2f %s",
		payload.OrderNumber, payload.CustomerName, payload.TotalAmount, payload.Currency,
	)

	err := notify.SendTelegramMessage(c.Config.Telegram.BotToken, c.Config.Telegram.ChatID, teleMsg)
	if err != nil {
		c.Logger.Errorf("❌ Lỗi gửi Telegram cho đơn %s: %v", payload.OrderNumber, err)
		// Không return err ngay vì ta vẫn muốn đi tiếp để gửi Email cho khách
	}

	// ==========================================
	// 2. GỬI EMAIL AWS SES CHO KHÁCH HÀNG
	// ==========================================
	if payload.CustomerEmail != "" {
		subject := fmt.Sprintf("Thank you for your order %s!", payload.OrderNumber)

		body := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; color: #333;">
				<h2>Hi %s,</h2>
				<p>We have successfully received your payment of <b>%.2f %s</b>.</p>
				<p>Your order <strong>%s</strong> is currently being prepared for shipment.</p>
				<br>
				<p>Thank you for shopping with us!</p>
			</div>`,
			payload.CustomerName, payload.TotalAmount, payload.Currency, payload.OrderNumber,
		)

		err := c.sesClient.SendEmail(ctx, payload.CustomerEmail, subject, body)
		if err != nil {
			c.Logger.Errorf("❌ Lỗi gửi Email qua SES cho %s: %v", payload.CustomerEmail, err)
		} else {
			c.Logger.Infof("📧 Đã gửi Email hóa đơn thành công cho khách hàng: %s", payload.CustomerEmail)
		}
	} else {
		c.Logger.Infof("⚠️ Đơn hàng %s không có Email khách hàng, bỏ qua bước gửi mail.", payload.OrderNumber)
	}

	c.Logger.Infof("✅ Hoàn tất xử lý Notification cho đơn hàng: %s", payload.OrderNumber)

	return nil
}
