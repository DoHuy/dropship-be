package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESClient quản lý kết nối đến AWS
type SESClient struct {
	client *ses.Client
	sender string
}

// Khởi tạo Client SES 1 lần duy nhất để tái sử dụng
func NewSESClient(ctx context.Context, region, accessKey, secretKey, senderEmail string) (*SESClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		// Nạp trực tiếp AccessKey và SecretKey
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("không thể load AWS config: %v", err)
	}

	return &SESClient{
		client: ses.NewFromConfig(cfg),
		sender: senderEmail,
	}, nil
}

// SendEmail gửi mail qua AWS SES
func (s *SESClient) SendEmail(ctx context.Context, toEmail, subject, htmlBody string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("lỗi gửi email qua AWS SES: %v", err)
	}
	return nil
}
