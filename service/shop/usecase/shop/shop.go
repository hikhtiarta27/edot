package shop

import (
	"shared"
	"shop/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type CreateRequest struct {
	Name string `json:"name"`
}

func (r CreateRequest) Validate() error {
	if err := validation.Validate(r.Name, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "name required",
		}
	}

	return nil
}

type Shop struct {
	ID        ulid.ULID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AssignWarehouseRequest struct {
	IDStr          string `param:"id"`
	ID             ulid.ULID
	WarehouseIDStr []string `json:"warehouse_id"`
	WarehouseID    []ulid.ULID
}

func (r *AssignWarehouseRequest) Validate() (err error) {
	if err := validation.Validate(r.WarehouseIDStr, validation.Required, validation.Min(1)); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "warehouse id required",
		}
	}

	for _, w := range r.WarehouseIDStr {
		id, err := ulid.Parse(w)
		if err != nil {
			return err
		}

		r.WarehouseID = append(r.WarehouseID, id)
	}

	r.ID, err = ulid.Parse(r.IDStr)
	if err != nil {
		return model.ErrInvalidShopID
	}

	return nil
}

type ShopWarehouse struct {
	*Shop
	Warehouse []ulid.ULID `json:"warehouse"`
}
