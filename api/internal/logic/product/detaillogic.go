package product

import (
	"api/internal/logic"
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ProductDetailReq) (resp *types.ProductDetailResp, err error) {
	data, err := l.svcCtx.ProductCache.GetData(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &types.ProductDetailResp{Info: *logic.ToProductInfo(data)}, nil
}
