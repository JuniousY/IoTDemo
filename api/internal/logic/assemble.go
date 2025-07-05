package logic

import (
	"api/internal/model"
	"api/internal/types"
)

func ToDeviceInfo(d *model.Device) *types.DeviceInfo {
	return &types.DeviceInfo{
		Id:         d.ID,
		ProductId:  d.ProductID,
		DeviceName: d.Name,
		Info:       d.Info,
		Status:     d.Status,
		IsOnline:   d.IsOnline,
		CreatedAt:  d.CreateTime.Format("2006-01-02 15:04:05"),
		UpdatedAt:  d.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}

func ToProductInfo(product *model.Product) *types.ProductInfo {
	return &types.ProductInfo{
		Id:   product.ID,
		Name: product.Name,
	}
}
