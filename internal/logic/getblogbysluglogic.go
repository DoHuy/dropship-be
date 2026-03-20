package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlogBySlugLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlogBySlugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlogBySlugLogic {
	return &GetBlogBySlugLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBlogBySlugLogic) GetBlogBySlug(in *dropshipbe.GetBlogBySlugRequest) (*dropshipbe.BlogDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.BlogDetailResponse{}, nil
}
