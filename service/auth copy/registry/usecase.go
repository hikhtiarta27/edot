package registry

import (
	"auth/infra"
	"auth/usecase"
	"sync"
)

var (
	accountUsecaseOnce sync.Once
	accountUsecase     usecase.Account
)

func LoadAccountUsecase() usecase.Account {
	accountUsecaseOnce.Do(func() {
		accountUsecase = usecase.NewAccount(
			LoadAccountRepo(),
			infra.LoadJWT(),
		)
	})

	return accountUsecase
}
