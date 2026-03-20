package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductBySlugLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductBySlugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductBySlugLogic {
	return &GetProductBySlugLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductBySlugLogic) GetProductBySlug(in *dropshipbe.GetProductBySlugRequest) (*dropshipbe.Product, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.Product{}, nil
}
