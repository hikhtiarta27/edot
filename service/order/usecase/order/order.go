package order

import (
	"order/model"
	"shared"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
)

type OrderDetail struct {
	ID          ulid.ULID `json:"id"`
	ProductID   ulid.ULID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Qty         uint64    `json:"qty"`
	Price       uint64    `json:"price"`
	TotalPrice  uint64    `json:"total_price"`
}

type Order struct {
	ID         ulid.ULID         `json:"id"`
	TotalItem  uint64            `json:"total_item"`
	TotalPrice uint64            `json:"total_price"`
	Status     model.OrderStatus `json:"status"`
	ExpiredAt  time.Time         `json:"expired_at"`
	CreatedAt  time.Time         `json:"created_at"`

	Detail []OrderDetail `json:"detail"`
}

type OrderProductRequest struct {
	ProductIDStr string `json:"product_id"`
	ProductID    ulid.ULID
	Qty          uint64
}

type CreateRequest struct {
	Product []OrderProductRequest `json:"product"`
}

func (r *CreateRequest) Validate() error {
	if err := validation.Validate(r.Product, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "product required",
		}
	}

	for i, prd := range r.Product {
		if err := validation.Validate(prd.ProductID, validation.Required); err != nil {
			return &shared.Error{
				HttpStatusCode: 400,
				Message:        "product id required",
			}
		}

		if err := validation.Validate(prd.Qty, validation.Required); err != nil {
			return &shared.Error{
				HttpStatusCode: 400,
				Message:        "qty required",
			}
		}

		prdID, err := ulid.Parse(prd.ProductIDStr)
		if err != nil {
			return &shared.Error{
				HttpStatusCode: 400,
				GrpcStatusCode: codes.InvalidArgument,
				Message:        "invalid product id",
			}
		}

		r.Product[i].ProductID = prdID
	}

	return nil
}
