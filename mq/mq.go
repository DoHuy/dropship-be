package main

import (
	"flag"
	"log"
	"os"

	"dropshipbe/mq/internal/config"
	"dropshipbe/mq/internal/consumer"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/mq.yaml", "the config file")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Hãy đảm bảo toàn bộ các biến môi trường trong .env được export: %v", err)
	}

	logMode := os.Getenv("LOG_MODE")
	logx.Infof("🔍 DEBUG: Chương trình đọc được biến LOG_MODE = '%s'\n", logMode)

	if logMode == "" {
		log.Fatalf("❌ Hãy đảm bảo toàn bộ các biến môi trường trong .env được export")
	}

	// 3. Load config với tuỳ chọn conf.UseEnv() để go-zero tự map biến
	var c config.Config
	err := conf.Load(*configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("❌ LỖI: Load file YAML thất bại: %v", err)
	}

	// Thiết lập cấu hình log
	logx.MustSetup(c.Log)

	sg := service.NewServiceGroup()
	defer sg.Stop()

	logx.Info("✅ Hệ thống Log đã sẵn sàng")

	notificationLogic := consumer.NewNotificationConsumer(c)
	notificationQueue := kq.MustNewQueue(c.NotificationConsumer, notificationLogic)
	sg.Add(notificationQueue)

	// ==========================================
	// 4. KHỞI CHẠY SERVICE
	// ==========================================
	logx.Info("🚀 Khởi động MQ Service thành công! Đang rình nghe các event từ Kafka...")
	sg.Start()
}
