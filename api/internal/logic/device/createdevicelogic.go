package device

import (
	"api/internal/logic"
	"api/internal/model"
	"api/internal/repo"
	"api/internal/utils"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func (l *CreateDeviceLogic) CreateDevice(req *types.CreateDeviceReq) (resp *types.DeviceResp, err error) {
	find, err := l.CheckDevice(req)
	if err != nil {
		l.Errorf("【创建设备】失败 req=%v err=%v\n", req, err)
		return nil, err
	} else if find == true {
		return nil, errors.New(fmt.Sprintf("【创建设备】设备ID重复:%s", req.DeviceName))
	}

	product, err := l.CheckProduct(req)
	if err != nil {
		l.Errorf("【创建设备】失败 req=%v err=%v\n", req, err)
		return nil, err
	} else if product == nil {
		return nil, errors.New(fmt.Sprintf("【创建设备】未查询到产品 产品id:%d", req.ProductId))
	}

	d := &model.Device{
		ProductID: req.ProductId,
		Name:      req.DeviceName,
		Status:    model.DeviceStatusInactive,
		IsOnline:  model.DeviceOffline,
	}
	d.Secret, err = utils.SecureRandomString(20)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("【创建设备】err:%d", err))
	}
	if req.Info != "" {
		d.Info = req.Info
	}

	err = l.svcCtx.DeviceRepo.Create(l.ctx, d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("【创建设备】err:%d", err))
	}

	return &types.DeviceResp{Info: *logic.ToDeviceInfo(d)}, nil
}

/*
CheckDevice
发现返回true 没有返回false
*/
func (l *CreateDeviceLogic) CheckDevice(req *types.CreateDeviceReq) (bool, error) {
	// 校验格式
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9-_]+$`, req.DeviceName); !ok {
		return false, errors.New("设备ID支持英文、数字、-,_ 的组合，最多不超过48个字符")
	}

	_, err := l.svcCtx.DeviceRepo.FindOneByFilter(l.ctx, repo.DeviceFilter{
		ProductID:  &req.ProductId,
		DeviceName: &req.DeviceName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

/*
CheckProduct
发现返回true 没有返回false
*/
func (l *CreateDeviceLogic) CheckProduct(req *types.CreateDeviceReq) (*model.Product, error) {
	p, err := l.svcCtx.ProductCache.GetData(l.ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return p, nil
}
