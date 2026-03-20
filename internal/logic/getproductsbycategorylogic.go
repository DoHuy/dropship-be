package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductsByCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductsByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsByCategoryLogic {
	return &GetProductsByCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductsByCategoryLogic) GetProductsByCategory(in *dropshipbe.GetProductsByCategoryRequest) (*dropshipbe.ProductListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.ProductListResponse{}, nil
}
