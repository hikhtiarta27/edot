package v1

import (
	"auth/usecase"
	"auth/usecase/account"
	"shared"

	"github.com/labstack/echo/v4"
)

type Account struct {
	accountUsecase usecase.Account
}

func NewAccount(
	accountUsecase usecase.Account,
) *Account {
	return &Account{
		accountUsecase: accountUsecase,
	}
}

func (d Account) Mount(group *echo.Group) {
	group.POST("/login", d.login)
	group.POST("/register", d.register)
}

func (d Account) login(c echo.Context) error {

	req := &account.LoginRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.accountUsecase.Login(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success login", res)
}

func (d Account) register(c echo.Context) error {

	req := &account.RegisterRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.accountUsecase.Register(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success register", res)
}
