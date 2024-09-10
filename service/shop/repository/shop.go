package repository

import (
	"context"
	"shared"
	"shop/model"

	"gorm.io/gorm"
)

type Shop interface {
	Select(ctx context.Context) ([]model.Shop, error)
	Create(ctx context.Context, shop *model.Shop) error
	Get(ctx context.Context, param *model.GetShop) (*model.Shop, error)
}

type shopRepo struct {
	db *gorm.DB
}

func NewShop(
	db *gorm.DB,
) Shop {
	return &shopRepo{
		db: db,
	}
}

func (r shopRepo) Select(ctx context.Context) ([]model.Shop, error) {

	var shops []model.Shop

	q := r.db.
		WithContext(ctx)

	err := q.
		Find(&shops).
		Error

	if err != nil {
		return nil, err
	}

	return shops, nil
}

func (r shopRepo) Create(ctx context.Context, shop *model.Shop) error {
	return r.db.
		WithContext(ctx).
		Create(&shop).
		Error
}

func (r shopRepo) Get(ctx context.Context, param *model.GetShop) (*model.Shop, error) {

	var shop *model.Shop

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ID) {
		q = q.Where("id = ?", param.ID)
	}

	err := q.
		First(&shop).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return shop, nil
}
