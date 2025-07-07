package auth

import (
	"api/internal/svc"
	"api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设备操作认证
func NewAccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessLogic {
	return &AccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccessLogic) Access(req *types.DeviceAccessReq) (resp *types.DeviceAccessResp, err error) {
	l.Logger.Info("AccessLogic", req)

	return &types.DeviceAccessResp{Result: "allow"}, nil
}
