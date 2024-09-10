package v1

import (
	"shared"
	"shop/usecase"
	"shop/usecase/shop"

	"github.com/labstack/echo/v4"
)

type Shop struct {
	shopUsecase usecase.Shop
}

func NewShop(
	shopUsecase usecase.Shop,
) *Shop {
	return &Shop{
		shopUsecase: shopUsecase,
	}
}

func (d Shop) Mount(group *echo.Group) {
	group.POST("", d.create)
	group.GET("", d.list)
}

func (d Shop) list(c echo.Context) error {

	res, err := d.shopUsecase.List(c.Request().Context())
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success get list shop", res)
}

func (d Shop) create(c echo.Context) error {

	req := &shop.CreateRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.shopUsecase.Create(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success create shop", res)
}
