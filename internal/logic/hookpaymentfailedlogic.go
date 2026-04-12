package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HookPaymentFailedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHookPaymentFailedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HookPaymentFailedLogic {
	return &HookPaymentFailedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HookPaymentFailedLogic) HookPaymentFailed(in *dropshipbe.PayPalWebhookRequest) (*dropshipbe.WebhookResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.WebhookResponse{}, nil
}
