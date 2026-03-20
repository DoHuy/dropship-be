package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryItemsLogic {
	return &GetCategoryItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoryItemsLogic) GetCategoryItems(in *dropshipbe.EmptyRequest) (*dropshipbe.CategoryListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.CategoryListResponse{}, nil
}
