package repository

import (
	"context"
	"database/sql"
	"errors"
	"shared"
	"warehouse/model"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Stock interface {
	Get(ctx context.Context, param *model.GetStock) (*model.Stock, error)
	Select(ctx context.Context, param *model.SelectStock) ([]model.Stock, error)
	Create(ctx context.Context, param *model.CreateStock) error
	ReserveRelease(ctx context.Context, param *model.ReserveReleaseStock) error
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

func (r *stockRepo) ReserveRelease(ctx context.Context, param *model.ReserveReleaseStock) error {

	var (
		rowsAffected int64
		err          error
	)

	if param.Action == model.StockRelease {
		res := r.db.
			WithContext(ctx).
			Raw("UPDATE stocks SET available_stock = available_stock + ?, reserved_stock = reserved_stock - ? WHERE id = ? AND reserved_stock - ? >= 0",
				param.Qty, param.Stock.ID, param.Qty,
			)

		err = res.Error
		rowsAffected = res.RowsAffected
	} else {
		res := r.db.
			WithContext(ctx).
			Raw("UPDATE stocks SET available_stock = available_stock - ?, reserved_stock = reserved_stock + ? WHERE id = ? AND available_stock - ? >= 0",
				param.Qty, param.Stock.ID, param.Qty,
			)

		err = res.Error
		rowsAffected = res.RowsAffected
	}

	if err != nil {

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1690:
				return errors.New("insuffient stock")
			}
		}

		return err
	}

	if rowsAffected == 0 {
		return errors.New("insuffient stock")
	}

	return nil
}
