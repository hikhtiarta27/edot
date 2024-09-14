package repository

import (
	"context"
	"database/sql"
	"order/model"
	"shared"

	"gorm.io/gorm"
)

type Order interface {
	Get(ctx context.Context, param *model.GetOrder) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, tx *gorm.DB, order *model.Order) error
}

type orderRepo struct {
	db        *gorm.DB
	stockRepo Stock
}

func NewOrder(
	db *gorm.DB,
	stockRepo Stock,
) Order {
	return &orderRepo{
		db:        db,
		stockRepo: stockRepo,
	}
}

func (r orderRepo) Get(ctx context.Context, param *model.GetOrder) (*model.Order, error) {

	var order *model.Order

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ID) {
		q = q.Where("id = ?", param.ID)
	}

	err := q.
		First(&order).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r orderRepo) Create(ctx context.Context, order *model.Order) (err error) {

	tx := r.db.Begin(&sql.TxOptions{})

	var successReserved []model.OrderDetail

	defer func() {

		// release reserved stock if there's an error when creating an order
		if err != nil {
			for _, rsv := range successReserved {
				_, _ = r.stockRepo.ReserveRelease(ctx, &model.ReserveReleaseStock{
					ProductID: rsv.ProductID,
					Qty:       rsv.Qty,
					Action:    model.StockReserve,
				})
			}
		}
	}()

	err = tx.
		WithContext(ctx).
		Create(&order).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, prd := range order.Detail {
		err = tx.WithContext(ctx).
			Create(&prd).
			Error

		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = r.stockRepo.ReserveRelease(ctx, &model.ReserveReleaseStock{
			ProductID: prd.ProductID,
			Qty:       prd.Qty,
			Action:    model.StockReserve,
		})

		if err != nil {
			tx.Rollback()
			return err
		}

		successReserved = append(successReserved, prd)
	}

	tx.Commit()
	return nil
}

func (r orderRepo) Update(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	return tx.
		WithContext(ctx).
		Updates(&order).
		Error
}
