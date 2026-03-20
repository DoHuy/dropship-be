package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductFaqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductFaqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductFaqsLogic {
	return &GetProductFaqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- FAQs ---
func (l *GetProductFaqsLogic) GetProductFaqs(in *dropshipbe.GetProductFaqsRequest) (*dropshipbe.FaqListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.FaqListResponse{}, nil
}
