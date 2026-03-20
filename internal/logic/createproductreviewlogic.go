package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

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

	return &dropshipbe.ReviewItem{}, nil
}
