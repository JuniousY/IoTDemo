package repo

import (
	"api/internal/model"
	"context"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

type ProductFilter struct {
	ID     int
	Status *int
}

func (r *ProductRepo) fmtFilter(ctx context.Context, f ProductFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&model.Product{})
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if f.Status != nil {
		db = db.Where("status = ?", f.Status)
	}
	return db
}

func (r *ProductRepo) FindOneByFilter(ctx context.Context, f ProductFilter) (*model.Product, error) {
	var result model.Product
	db := r.fmtFilter(ctx, f)
	err := db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
