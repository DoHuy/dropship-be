package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsLogic {
	return &GetProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Products ---
func (l *GetProductsLogic) GetProducts(in *dropshipbe.EmptyRequest) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ProductListResponse{}, nil
}
