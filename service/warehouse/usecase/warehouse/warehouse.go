package warehouse

import (
	"shared"
	"time"
	"warehouse/model"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type Warehouse struct {
	ID        ulid.ULID             `json:"id"`
	Name      string                `json:"name"`
	Status    model.WarehouseStatus `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt *time.Time            `json:"updated_at"`
}

type WarehouseTransfer struct {
	ID              ulid.ULID `json:"id"`
	FromWarehouseID ulid.ULID `json:"from_warehouse_id"`
	ToWarehouseID   ulid.ULID `json:"to_warehouse_Id"`
	ProductID       ulid.ULID `json:"product_id"`
	Stock           uint64    `json:"stock"`
	CreatedAt       time.Time `json:"created_at"`
}

type ListRequest struct {
}

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

	if err := validation.Validate(r.Name, validation.RuneLength(3, 0)); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "name minimum length is 3",
		}
	}

	return nil
}

type UpdateRequest struct {
	ID     string                `param:"id"`
	Name   string                `json:"name"`
	Status model.WarehouseStatus `json:"status"`
}

func (r UpdateRequest) Validate() error {

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

	if err := validation.Validate(r.Status, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "status required",
		}
	}

	if err := validation.Validate(r.Status, validation.In(model.WarehouseActive, model.WarehouseInactive)); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "status invalid format. must be ACTIVE or INACTIVE",
		}
	}

	return nil
}

type TransferRequest struct {
	FromWarehouseIDStr string `json:"from_warehouse_id"`
	ToWarehouseIDStr   string `json:"to_warehouse_id"`
	ProductIDStr       string `json:"product_id"`
	Stock              uint64 `json:"stock"`
	FromWarehouseID    ulid.ULID
	ToWarehouseID      ulid.ULID
	ProductID          ulid.ULID
}

func (r *TransferRequest) Validate() (err error) {
	if err = validation.Validate(r.FromWarehouseIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "from warehouse id required",
		}
	}

	r.FromWarehouseID, err = ulid.Parse(r.FromWarehouseIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid from warehouse id",
		}
	}

	if err = validation.Validate(r.ToWarehouseIDStr, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "to warehouse id required",
		}
	}

	r.ToWarehouseID, err = ulid.Parse(r.ToWarehouseIDStr)
	if err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "invalid to warehouse id",
		}
	}

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

	if err = validation.Validate(r.Stock, validation.NotNil); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "stock required",
		}
	}

	return nil
}
