package registry

import (
	"auth/infra"
	"auth/repository"
	"sync"
)

var (
	accountRepoOnce sync.Once
	accountRepo     repository.Account
)

func LoadAccountRepo() repository.Account {
	accountRepoOnce.Do(func() {
		accountRepo = repository.NewAccount(
			infra.LoadDB(),
		)
	})

	return accountRepo
}
