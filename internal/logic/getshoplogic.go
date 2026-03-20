package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopLogic {
	return &GetShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Shop Search ---
func (l *GetShopLogic) GetShop(in *dropshipbe.ShopSearchParams) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ProductListResponse{}, nil
}
