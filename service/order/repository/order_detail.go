package repository

import (
	"context"
	"order/model"
	"shared"

	"gorm.io/gorm"
)

type OrderDetail interface {
	Select(ctx context.Context, param *model.SelectOrderDetail) ([]model.OrderDetail, error)
}

type orderDetailRepo struct {
	db *gorm.DB
}

func NewOrderDetail(
	db *gorm.DB,
) OrderDetail {
	return &orderDetailRepo{
		db: db,
	}
}

func (r orderDetailRepo) Select(ctx context.Context, param *model.SelectOrderDetail) ([]model.OrderDetail, error) {

	var orderDetails []model.OrderDetail

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.OrderID) {
		q = q.Where("order_id = ?", param.OrderID)
	}

	err := q.
		Find(&orderDetails).
		Error

	if err != nil {
		return nil, err
	}

	return orderDetails, nil
}
