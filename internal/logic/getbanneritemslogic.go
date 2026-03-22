package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBannerItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBannerItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBannerItemsLogic {
	return &GetBannerItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBannerItemsLogic) GetBannerItems(in *dropshipbe.DefaultRequest) (*dropshipbe.BannerListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.BannerListResponse{}, nil
}
