package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeaturedProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeaturedProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeaturedProductsLogic {
	return &GetFeaturedProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeaturedProductsLogic) GetFeaturedProducts(in *dropshipbe.EmptyRequest) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ProductListResponse{}, nil
}
