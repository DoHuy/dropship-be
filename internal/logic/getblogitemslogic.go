package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlogItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlogItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlogItemsLogic {
	return &GetBlogItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Blogs ---
func (l *GetBlogItemsLogic) GetBlogItems(in *dropshipbe.EmptyRequest) (*dropshipbe.BlogListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.BlogListResponse{}, nil
}
