package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoBannerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoBannerLogic {
	return &GetVideoBannerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoBannerLogic) GetVideoBanner(in *dropshipbe.DefaultRequest) (*dropshipbe.Banner, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.Banner{}, nil
}
