package v1

import (
	"shared"
	"warehouse/usecase"
	"warehouse/usecase/warehouse"

	"github.com/labstack/echo/v4"
)

type Warehouse struct {
	warehouseUsecase usecase.Warehouse
}

func NewWarehouse(
	warehouseUsecase usecase.Warehouse,
) *Warehouse {
	return &Warehouse{
		warehouseUsecase: warehouseUsecase,
	}
}

func (d Warehouse) Mount(group *echo.Group) {
	group.GET("", d.list)
	group.POST("", d.create)
	group.PUT("/:id", d.update)
	group.POST("/transfer-stock", d.transferStock)
}

func (d Warehouse) list(c echo.Context) error {

	req := &warehouse.ListRequest{}

	_ = c.Bind(req)

	res, err := d.warehouseUsecase.List(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success get warehouse", res)
}

func (d Warehouse) create(c echo.Context) error {

	req := &warehouse.CreateRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.warehouseUsecase.Create(c.Request().Context(), req)
	if err != nil {

		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success create warehouse", res)
}

func (d Warehouse) update(c echo.Context) error {

	req := &warehouse.UpdateRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.warehouseUsecase.Update(c.Request().Context(), req)
	if err != nil {

		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success update warehouse", res)
}

func (d Warehouse) transferStock(c echo.Context) error {

	req := &warehouse.TransferRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.warehouseUsecase.TransferStock(c.Request().Context(), req)
	if err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success transfer stock", res)
}
