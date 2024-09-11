package product

import (
	"shared"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type ListRequest struct {
	Keyword string `query:"q"`
}

type Product struct {
	ID             ulid.ULID `json:"id"`
	Slug           string    `json:"slug"`
	Name           string    `json:"name"`
	Price          uint64    `json:"price"`
	ShopID         ulid.ULID `json:"shop_id"`
	AvailableStock uint64    `json:"available_stock"`
	ReservedStock  uint64    `json:"reserved_stock"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name           string `json:"name"`
	Price          uint64 `json:"price"`
	Stock          uint64 `json:"stock"`
	ShopIDStr      string `json:"shop_id"`
	ShopID         ulid.ULID
	WarehouseIDStr string `json:"warehouse_id"`
	WarehouseID    ulid.ULID
}

func (r *CreateRequest) Validate() (err error) {
	if err := validation.Validate(r.Name, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "name required",
		}
	}

	if err := validation.Validate(r.Name, validation.RuneLength(3, 0)); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "name minimum length is 3",
		}
	}

	if err := validation.Validate(r.Price, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "price required",
		}
	}

	if err := validation.Validate(r.Stock, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "stock required",
		}
	}

	if err := validation.Validate(r.ShopIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "shop id required",
		}
	}

	r.ShopID, err = ulid.Parse(r.ShopIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid shop id",
		}
	}

	if err := validation.Validate(r.WarehouseIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "shop id required",
		}
	}

	r.WarehouseID, err = ulid.Parse(r.WarehouseIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid warehouse id",
		}
	}

	return nil
}
