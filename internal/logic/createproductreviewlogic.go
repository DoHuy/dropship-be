package logic

import (
	"context"
	"time"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductReviewLogic {
	return &CreateProductReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductReviewLogic) CreateProductReview(in *dropshipbe.CreateProductReviewRequest) (*dropshipbe.ReviewItem, error) {
	// todo: add your logic here and delete this line
	review, err := l.svcCtx.EcommerceRepo.CreateProductReview(l.ctx, in)
	if err != nil {
		logx.Errorf("Lỗi khi tạo review sản phẩm: %v", err)
		return nil, err
	}

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// presigned media URLs
	if len(review.Media.Images) > 0 {
		for i, imgKey := range review.Media.Images {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(imgKey),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho ảnh review: %v", err)
				review.Media.Images[i] = "" // fallback nếu có lỗi
			} else {
				review.Media.Images[i] = presignedURL.URL
			}
		}
	}
	if len(review.Media.Videos) > 0 {
		for i, videoKey := range review.Media.Videos {
			presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
				Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
				Key:    aws.String(videoKey),
			}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
			if err != nil {
				logx.Errorf("Lỗi khi tạo presigned URL cho video review: %v", err)
				review.Media.Videos[i] = "" // fallback nếu có lỗi
			} else {
				review.Media.Videos[i] = presignedURL.URL
			}
		}
	}

	if review.AuthorAvatar != "" {
		presignedURL, err := l.svcCtx.PresignClient.PresignGetObject(contextWithTimeout, &s3.GetObjectInput{
			Bucket: aws.String(l.svcCtx.Config.R2.BucketName),
			Key:    aws.String(review.AuthorAvatar),
		}, s3.WithPresignExpires(time.Duration(l.svcCtx.Config.R2.LinkExpiration)*time.Minute))
		if err != nil {
			logx.Errorf("Lỗi khi tạo presigned URL cho avatar review: %v", err)
			review.AuthorAvatar = "" // fallback nếu có lỗi
		} else {
			review.AuthorAvatar = presignedURL.URL
		}
	}

	return &dropshipbe.ReviewItem{
		Id:       review.ID,
		Name:     review.AuthorName,
		Avatar:   review.AuthorAvatar,
		Rating:   int32(review.Rating),
		Comment:  review.Content,
		Verified: review.IsVerified,
		Images:   review.Media.Images,
		Videos:   review.Media.Videos,
		Date:     review.CreatedAt.Format(time.RFC3339),
	}, nil
}
