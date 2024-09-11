package infra

import (
	"shared"
	"sync"
)

var (
	jwtOnce sync.Once
	jwt     *shared.Jwt
)

func LoadJWT() *shared.Jwt {
	jwtOnce.Do(func() {
		jwt = shared.New(LoadConfig().JWT.Secret)
	})

	return jwt
}
