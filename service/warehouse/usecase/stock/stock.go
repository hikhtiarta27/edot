package stock

import (
	"shared"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type CreateRequest struct {
	ProductIDStr   string `json:"product_id"`
	ProductID      ulid.ULID
	Stock          uint64 `json:"stock"`
	WarehouseIDStr string `json:"warehouse_id"`
	WarehouseID    ulid.ULID
}

func (r *CreateRequest) Validate() (err error) {
	if err = validation.Validate(r.ProductIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "product id required",
		}
	}

	r.ProductID, err = ulid.Parse(r.ProductIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid product id",
		}
	}

	if err = validation.Validate(r.WarehouseIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "warehouse id required",
		}
	}

	r.WarehouseID, err = ulid.Parse(r.WarehouseIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid warehouse id",
		}
	}

	if err = validation.Validate(r.Stock, validation.NotNil); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "stock required",
		}
	}

	return nil
}

type GetRequest struct {
	ProductIDStr string `json:"product_id"`
	ProductID    ulid.ULID
}

func (r *GetRequest) Validate() (err error) {
	if err = validation.Validate(r.ProductIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "product id required",
		}
	}

	r.ProductID, err = ulid.Parse(r.ProductIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid product id",
		}
	}
	return nil
}

type Stock struct {
	ID             ulid.ULID `json:"id"`
	AvailableStock uint64    `json:"available_stock"`
	ReservedStock  uint64    `json:"reserved_stock"`
	ProductID      ulid.ULID `json:"product_id"`
}
