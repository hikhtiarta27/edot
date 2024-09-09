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
	AvailableStock uint64    `json:"available_stock"`
	ReservedStock  uint64    `json:"reserved_stock"`
	CreatedAt      time.Time `json:"created_at"`
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
