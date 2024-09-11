package v1

import (
	"shared"
	"warehouse/usecase"

	"github.com/labstack/echo/v4"
)

type Stock struct {
	stockUsecase usecase.Stock
}

func NewStock(
	stockUsecase usecase.Stock,
) *Stock {
	return &Stock{
		stockUsecase: stockUsecase,
	}
}

func (d Stock) Mount(group *echo.Group) {
	group.POST("", d.create)
}

func (d Stock) create(c echo.Context) error {

	// req := &stock.CreateRequest{}

	// if err := c.Bind(req); err != nil {
	// 	return shared.FailResponseFromCustomError(c, err)
	// }

	// res, err := d.stockUsecase.Create(c.Request().Context(), req)
	// if err != nil {

	// 	return shared.FailResponseFromCustomError(c, err)
	// }

	return shared.SuccessResponse(c, "deprecated: success create stock", nil)
}
