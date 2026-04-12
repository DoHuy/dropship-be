package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptureOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptureOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptureOrderLogic {
	return &CaptureOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptureOrderLogic) CaptureOrder(in *dropshipbe.CaptureOrderRequest) (*dropshipbe.CaptureOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.CaptureOrderResponse{}, nil
}
