package repository

import (
	"context"

	"gorm.io/gorm"
)

type Tx interface {
	DoInTransaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) (err error)
}

type tx struct {
	db *gorm.DB
}

func NewTx(
	db *gorm.DB,
) Tx {
	return &tx{
		db: db,
	}
}

func (t *tx) DoInTransaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) (err error) {

	tx := t.db.
		WithContext(ctx).
		Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}
