package repo

import (
	"api/internal/model"
	"context"
	"gorm.io/gorm"
)

type deviceInfoRepo struct {
	db *gorm.DB
}

type DeviceFilter struct {
	ID         int64
	IDs        []int64
	ProductIDs []string
}

type DeviceInfoRepo interface {
	Create(ctx context.Context, device *model.Device) error
	FindOneByFilter(ctx context.Context, f DeviceFilter) (*model.Device, error)
	CountByFilter(ctx context.Context, f DeviceFilter) (int64, error)
	FindByFilter(ctx context.Context, f DeviceFilter, offset, limit int) ([]*model.Device, error)
	UpdateOnlineStatus(ctx context.Context, id int64, isOnline int) error
}

func NewDeviceInfoRepo(db *gorm.DB) DeviceInfoRepo {
	return &deviceInfoRepo{db: db}
}

func (r *deviceInfoRepo) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Model(&model.Device{}).Create(device).Error
}

func (r *deviceInfoRepo) fmtFilter(ctx context.Context, f DeviceFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&model.Device{})
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if len(f.ProductIDs) > 0 {
		db = db.Where("id in ?", f.ProductIDs)
	}
	return db
}

func (r *deviceInfoRepo) FindOneByFilter(ctx context.Context, f DeviceFilter) (*model.Device, error) {
	var result model.Device
	db := r.fmtFilter(ctx, f)
	err := db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *deviceInfoRepo) CountByFilter(ctx context.Context, f DeviceFilter) (int64, error) {
	var total int64
	if err := r.fmtFilter(ctx, f).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *deviceInfoRepo) FindByFilter(ctx context.Context, f DeviceFilter, offset, limit int) ([]*model.Device, error) {
	var devices []*model.Device
	if err := r.fmtFilter(ctx, f).
		Offset(offset).Limit(limit).
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceInfoRepo) UpdateOnlineStatus(ctx context.Context, id int64, isOnline int) error {
	return r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ?", id).
		Update("is_online", isOnline).Error
}
