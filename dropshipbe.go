package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/internal/config"
	"dropshipbe/internal/server"
	"dropshipbe/internal/svc"
)

var configFile = flag.String("f", "etc/dropshipbe.yaml", "the config file")

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ LỖI: Không thể đọc file .env! Hãy đảm bảo file tên là '.env' nằm cùng thư mục với dropshipbe.go. Chi tiết: %v", err)
	}

	// 2. [DEBUG] Kiểm tra thực tế giá trị biến LOG_MODE
	logMode := os.Getenv("LOG_MODE")
	fmt.Printf("🔍 DEBUG: Chương trình đọc được biến LOG_MODE = '%s'\n", logMode)

	if logMode == "" {
		log.Fatalf("❌ LỖI: File .env đã được đọc, nhưng LOG_MODE bị trống. Hãy kiểm tra lại nội dung file .env!")
	}

	// 3. Load config với tuỳ chọn conf.UseEnv() để go-zero tự map biến
	var c config.Config
	err := conf.Load(*configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("❌ LỖI: Load file YAML thất bại: %v", err)
	}

	// Thiết lập cấu hình log
	logx.MustSetup(c.Log)
	ctx := svc.NewServiceContext(c)
	logx.Info("✅ Hệ thống Log đã sẵn sàng")

	// Pre-warming tồn kho
	logx.Info("🔄 Đang đồng bộ tồn kho từ Database lên Redis...")
	if err := utils.LoadAllInventoryToRedis(context.Background(), ctx.Redis, ctx.DB); err != nil {
		logx.WithContext(context.Background()).Errorf("❌ Lỗi nạp tồn kho lên Redis: %v", err)
		panic(fmt.Sprintf("❌ Lỗi nạp tồn kho lên Redis: %v", err))
	}
	logx.Info("✅ Nạp tồn kho thành công. Hệ thống sẵn sàng!")

	// Khởi chạy server
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		dropshipbe.RegisterDropshipbeServer(grpcServer, server.NewDropshipbeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
