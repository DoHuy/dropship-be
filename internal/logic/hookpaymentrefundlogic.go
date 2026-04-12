package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HookPaymentRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHookPaymentRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HookPaymentRefundLogic {
	return &HookPaymentRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HookPaymentRefundLogic) HookPaymentRefund(in *dropshipbe.PayPalWebhookRequest) (*dropshipbe.WebhookResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.WebhookResponse{}, nil
}
