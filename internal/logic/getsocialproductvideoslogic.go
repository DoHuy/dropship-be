package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSocialProductVideosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSocialProductVideosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSocialProductVideosLogic {
	return &GetSocialProductVideosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- Media/Galleries ---
func (l *GetSocialProductVideosLogic) GetSocialProductVideos(in *dropshipbe.GetSocialProductVideoRequest) (*dropshipbe.GalleryListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.GalleryListResponse{}, nil
}
