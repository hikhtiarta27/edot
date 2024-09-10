package repository

import (
	"auth/model"
	"context"
	"shared"

	"gorm.io/gorm"
)

type Account interface {
	Get(ctx context.Context, param *model.GetAccount) (*model.Account, error)
	Create(ctx context.Context, account *model.Account) error
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccount(
	db *gorm.DB,
) Account {
	return &accountRepo{
		db: db,
	}
}

func (r accountRepo) Get(ctx context.Context, param *model.GetAccount) (*model.Account, error) {

	var account *model.Account

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ID) {
		q = q.Where("id = ?", param.ID)
	}

	if param.Username != "" {
		q = q.Where("username = ?", param.Username)
	}

	err := q.
		First(&account).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	return account, nil
}

func (r accountRepo) Create(ctx context.Context, account *model.Account) error {
	err := r.db.
		WithContext(ctx).
		Create(&account).
		Error

	if err != nil {
		return err
	}

	return nil
}
