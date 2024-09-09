package model

import "github.com/oklog/ulid/v2"

type GetAccount struct {
	ID       ulid.ULID
	Username string
}
