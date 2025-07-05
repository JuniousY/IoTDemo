package repo

import (
	"api/internal/model"
	"context"
	"gorm.io/gorm"
)

type deviceRepo struct {
	db *gorm.DB
}

type DeviceFilter struct {
	ID         *int64
	IDs        []int64
	ProductID  *int
	ProductIDs []int
	DeviceName *string
	Status     *int
}

type DeviceRepo interface {
	Create(ctx context.Context, device *model.Device) error
	FindOneByFilter(ctx context.Context, f DeviceFilter) (*model.Device, error)
	CountByFilter(ctx context.Context, f DeviceFilter) (int64, error)
	FindByFilter(ctx context.Context, f DeviceFilter, offset, limit int) ([]*model.Device, error)
	UpdateOnlineStatus(ctx context.Context, id int64, isOnline int) error
}

func NewDeviceRepo(db *gorm.DB) DeviceRepo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Model(&model.Device{}).Create(device).Error
}

func (r *deviceRepo) fmtFilter(ctx context.Context, f DeviceFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&model.Device{})
	if f.ID != nil {
		db = db.Where("id = ?", f.ID)
	}
	if f.ProductID != nil {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if len(f.ProductIDs) > 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if f.DeviceName != nil {
		db = db.Where("name = ?", *f.DeviceName)
	}
	if f.Status != nil {
		db = db.Where("status = ?", f.Status)
	} else {
		db = db.Where("status != ?", model.DeviceStatusDeleted)
	}
	return db
}

func (r *deviceRepo) FindOneByFilter(ctx context.Context, f DeviceFilter) (*model.Device, error) {
	var result model.Device
	db := r.fmtFilter(ctx, f)
	err := db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *deviceRepo) CountByFilter(ctx context.Context, f DeviceFilter) (int64, error) {
	var total int64
	if err := r.fmtFilter(ctx, f).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *deviceRepo) FindByFilter(ctx context.Context, f DeviceFilter, offset, limit int) ([]*model.Device, error) {
	var devices []*model.Device
	if err := r.fmtFilter(ctx, f).
		Offset(offset).Limit(limit).
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepo) UpdateOnlineStatus(ctx context.Context, id int64, isOnline int) error {
	return r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ?", id).
		Update("is_online", isOnline).Error
}
