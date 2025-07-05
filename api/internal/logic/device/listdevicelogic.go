package device

import (
	"api/internal/logic"
	"api/internal/repo"
	"context"
	"errors"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDeviceLogic {
	return &ListDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDeviceLogic) ListDevice(req *types.ListDeviceReq) (resp *types.ListDeviceResp, err error) {
	if req.Size <= 0 || req.Page <= 0 {
		return &types.ListDeviceResp{Total: 0}, errors.New("invalid page or size")
	}
	f := repo.DeviceFilter{}

	total, err := l.svcCtx.DeviceRepo.CountByFilter(l.ctx, f)
	if err != nil {
		return &types.ListDeviceResp{Total: 0}, err
	}
	if total == 0 {
		return &types.ListDeviceResp{Total: 0}, nil
	}

	devices, err := l.svcCtx.DeviceRepo.FindByFilter(l.ctx, f, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		return &types.ListDeviceResp{Total: 0}, err
	}

	var list []types.DeviceInfo
	for _, d := range devices {
		list = append(list, *logic.ToDeviceInfo(d))
	}
	return &types.ListDeviceResp{
		Total: total,
		List:  list,
	}, nil
}
