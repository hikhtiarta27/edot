package repository

import (
	"context"
	"database/sql"
	"shared"
	"warehouse/model"

	"gorm.io/gorm"
)

type Stock interface {
	Get(ctx context.Context, param *model.GetStock) (*model.Stock, error)
	Select(ctx context.Context, param *model.SelectStock) ([]model.Stock, error)
	Create(ctx context.Context, param *model.CreateStock) error
}

type stockRepo struct {
	db *gorm.DB
}

func NewStock(db *gorm.DB) Stock {
	return &stockRepo{db: db}
}

func (r *stockRepo) Get(ctx context.Context, param *model.GetStock) (*model.Stock, error) {
	stock := new(model.Stock)

	query := r.db.WithContext(ctx)

	if !shared.IsZero(param.ID) {
		query = query.Where("id = ?", param.ID)
	}

	if !shared.IsZero(param.ProductID) {
		query = query.Where("product_id = ?", param.ProductID)
	}

	err := query.First(stock).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (r *stockRepo) Select(ctx context.Context, param *model.SelectStock) ([]model.Stock, error) {
	var stocks []model.Stock

	q := r.db.WithContext(ctx)

	if len(param.ProductIDs) > 0 {
		q = q.Where("product_id IN ?", param.ProductIDs)
	}

	err := q.Find(&stocks).Error
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *stockRepo) Create(ctx context.Context, param *model.CreateStock) error {

	tx := r.db.Begin(&sql.TxOptions{})

	err := tx.Create(&param.Stock).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(&param.WarehouseTransfer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
