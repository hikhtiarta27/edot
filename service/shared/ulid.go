package shared

import (
	"github.com/oklog/ulid/v2"
)

func IsZero(id ulid.ULID) bool {
	return id.String() == "00000000000000000000000000"
}
