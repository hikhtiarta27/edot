package model

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Shop struct {
	ID        ulid.ULID
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (Shop) TableName() string {
	return "shops"
}

func NewShop(
	name string,
) (*Shop, error) {
	instance := &Shop{
		ID:   ulid.Make(),
		Name: strings.TrimSpace(name),
	}

	return instance, nil
}
