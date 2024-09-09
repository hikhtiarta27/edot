package v1

import (
	"log"
	"product/usecase"
	"product/usecase/product"
	"shared"

	"github.com/labstack/echo/v4"
)

type Product struct {
	productUsecase usecase.Product
}

func NewProduct(
	productUsecase usecase.Product,
) *Product {
	return &Product{
		productUsecase: productUsecase,
	}
}

func (d Product) Mount(group *echo.Group) {
	group.GET("", d.list)
	group.POST("", d.create)
}

func (d Product) list(c echo.Context) error {

	req := &product.ListRequest{}

	_ = c.Bind(req)

	res, err := d.productUsecase.List(c.Request().Context(), req)
	if err != nil {
		log.Println(err)
		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success get product", res)
}

func (d Product) create(c echo.Context) error {

	req := &product.CreateRequest{}

	if err := c.Bind(req); err != nil {
		return shared.FailResponseFromCustomError(c, err)
	}

	res, err := d.productUsecase.Create(c.Request().Context(), req)
	if err != nil {

		return shared.FailResponseFromCustomError(c, err)
	}

	return shared.SuccessResponse(c, "success create product", res)
}
