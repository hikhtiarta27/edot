package v1

import (
	"order/usecase"
	"order/usecase/order"
	"shared"

	"github.com/labstack/echo/v4"
)

type Order struct {
	orderUsecase usecase.Order
}

func NewOrder(
	orderUsecase usecase.Order,
) *Order {
	return &Order{
		orderUsecase: orderUsecase,
	}
}

func (d Order) Mount(group *echo.Group) {
	group.POST("", d.create)
	group.GET("/release", d.release)
}

func (d Order) create(c echo.Context) error {

	req := &order.CreateRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.orderUsecase.Create(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success create order", res)
}

func (d Order) release(c echo.Context) error {

	err := d.orderUsecase.Release(c.Request().Context())
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success release order", nil)
}
