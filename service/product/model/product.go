package model

import (
	"regexp"
	"shared"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var ErrProductNotFound = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.NotFound,
	Message:        "product not found",
}

var ErrInvalidID = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.InvalidArgument,
	Message:        "invalid id",
}

type Product struct {
	ID        ulid.ULID
	Slug      string
	Name      string
	Price     uint64
	ShopID    ulid.ULID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (Product) TableName() string {
	return "products"
}

func NewProduct(
	name string,
	price uint64,
	shopID ulid.ULID,
) (*Product, error) {
	instance := &Product{
		ID:     ulid.Make(),
		Name:   strings.TrimSpace(name),
		Price:  price,
		ShopID: shopID,
	}

	if err := instance.SetSlug(name); err != nil {
		return nil, err
	}

	return instance, nil
}

func (m *Product) SetSlug(str string) error {

	str = strings.ToLower(strings.TrimSpace(str))

	str = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(str, "-")

	str = strings.Trim(str, "-")

	m.Slug = str
	return nil
}
