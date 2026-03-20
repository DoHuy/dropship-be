package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductReviewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductReviewsLogic {
	return &GetProductReviewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Reviews ---
func (l *GetProductReviewsLogic) GetProductReviews(in *dropshipbe.GetProductReviewsRequest) (*dropshipbe.ReviewSummary, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ReviewSummary{}, nil
}
