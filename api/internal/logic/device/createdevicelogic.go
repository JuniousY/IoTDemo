package device

import (
	"context"
	"regexp"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDeviceLogic {
	return &CreateDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDeviceLogic) CreateDevice(req *types.CreateDeviceReq) (resp *types.CreateDeviceResp, err error) {
	// todo: add your logic here and delete this line

	return
}

/*
CheckDevice
发现返回true 没有返回false
*/
//func (l *CreateDeviceLogic) CheckDevice(in *dm.DeviceInfo) (bool, error) {
//	result, _ := regexp.MatchString(`^[a-zA-Z0-9-_]+$`, in.DeviceName)
//	if !result {
//		return false, errors.Parameter.AddMsg("设备ID支持英文、数字、-,_ 的组合，最多不超过48个字符")
//	}
//	_, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
//	if err == nil {
//		return true, nil
//	}
//	if errors.Cmp(err, errors.NotFind) {
//		return false, nil
//	}
//	return false, err
//}
//
///*
//CheckProduct
//发现返回true 没有返回false
//*/
//func (l *CreateDeviceLogic) CheckProduct(in *dm.DeviceInfo) (*dm.ProductInfo, error) {
//	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.ProductID)
//	if err == nil {
//		return pi, nil
//	}
//	if errors.Cmp(err, errors.NotFind) {
//		return nil, nil
//	}
//	return nil, err
//}
