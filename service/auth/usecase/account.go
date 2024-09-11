package usecase

import (
	"auth/model"
	"auth/repository"
	"auth/usecase/account"
	"context"
	"errors"
	"shared"
	"time"
)

type Account interface {
	Login(ctx context.Context, param *account.LoginRequest) (*account.LoginResponse, error)
	Register(ctx context.Context, param *account.RegisterRequest) (*account.LoginResponse, error)
}

type accountUsecase struct {
	accountRepo repository.Account
	jwt         *shared.Jwt
}

func NewAccount(
	accountRepo repository.Account,
	jwt *shared.Jwt,
) Account {
	return &accountUsecase{
		accountRepo: accountRepo,
		jwt:         jwt,
	}
}

func (s *accountUsecase) Login(ctx context.Context, param *account.LoginRequest) (*account.LoginResponse, error) {

	ctx, parentSpan := tracer.Start(ctx, "accountUsecase.login")
	defer parentSpan.End()

	if err := param.Validate(); err != nil {
		return nil, err
	}

	accountModel, err := s.accountRepo.Get(ctx, &model.GetAccount{
		Username: param.Username,
	})

	if err != nil {
		return nil, err
	}

	if accountModel == nil {
		return nil, model.ErrUserNotFound
	}

	if err := accountModel.ComparePassword(param.Password); err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(time.Hour*24, accountModel.ID.String())
	if err != nil {
		return nil, err
	}

	return &account.LoginResponse{
		Type:  "Bearer",
		Token: token,
	}, nil
}

func (s *accountUsecase) Register(ctx context.Context, param *account.RegisterRequest) (*account.LoginResponse, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	existingAccount, err := s.accountRepo.Get(ctx, &model.GetAccount{
		Username: param.Username,
	})

	if err != nil {
		return nil, err
	}

	if existingAccount != nil {
		return nil, errors.New("user already exist. try another username")
	}

	newAccount, err := model.NewAccount(
		param.Name,
		param.Username,
		param.Password,
	)
	if err != nil {
		return nil, err
	}

	err = s.accountRepo.Create(ctx, newAccount)
	if err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(time.Hour*24, newAccount.ID.String())
	if err != nil {
		return nil, err
	}

	return &account.LoginResponse{
		Type:  "Bearer",
		Token: token,
	}, nil
}
