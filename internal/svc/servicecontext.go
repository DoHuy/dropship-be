package svc

import (
	"context"
	"dropshipbe/internal/config"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ServiceContext struct {
	Config        config.Config
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
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

	return &ServiceContext{
		Config:        c,
		S3Client:      s3Client,
		PresignClient: presignClient,
	}
}
