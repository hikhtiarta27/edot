package model

import "github.com/oklog/ulid/v2"

type SelectStock struct {
	ProductIDs []ulid.ULID
}
