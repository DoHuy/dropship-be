package logic

import (
	"context"

	"dropshipbe/dropshipbe"
	"dropshipbe/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSliderItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSliderItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSliderItemsLogic {
	return &GetSliderItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// --- UI Items (Sliders, Categories, Banners) ---
func (l *GetSliderItemsLogic) GetSliderItems(in *dropshipbe.EmptyRequest) (*dropshipbe.SliderListResponse, error) {
	// todo: add your logic here and delete this line

	return &dropshipbe.SliderListResponse{}, nil
}
