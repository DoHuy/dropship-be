package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNewProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNewProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNewProductsLogic {
	return &GetNewProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNewProductsLogic) GetNewProducts(in *dropshipbe.EmptyRequest) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ProductListResponse{}, nil
}
