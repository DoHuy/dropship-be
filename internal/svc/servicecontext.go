package svc

import (
	"context"
	"dropshipbe/internal/config"
	"dropshipbe/model/repository"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config        config.Config
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	EcommerceRepo repository.EcommerceRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. Tự động tạo Endpoint từ Account ID của bạn
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.R2.AccountID)

	// 2. Khởi tạo cấu hình AWS với thông tin xác thực từ config
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.R2.AccessKey,
			c.R2.SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Không thể tải cấu hình R2: %v", err)
	}

	// 3. Tạo S3 Client chính để upload
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// 4. Tạo Presign Client để sinh link có thời hạn
	presignClient := s3.NewPresignClient(s3Client)

	// 5. Kết nối Database với GORM

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          //  color
		},
	)

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Ho_Chi_Minh", c.DB.Host, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.Port, c.DB.SSLMode)), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v", err)
	}

	// Tùy chọn: Cấu hình Connection Pool cho GORM để tăng hiệu suất
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(c.DB.MaxOpenConns)
		sqlDB.SetMaxIdleConns(c.DB.MaxIdleConns)
	}

	// Khởi tạo Repository với kết nối DB vừa tạo
	ecomRepo := repository.NewEcommerceRepository(db)

	return &ServiceContext{
		Config:        c,
		S3Client:      s3Client,
		PresignClient: presignClient,
		EcommerceRepo: ecomRepo,
	}
}
