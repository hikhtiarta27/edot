package registry

import (
	"product/infra"
	"product/repository"
	"sync"
)

var (
	productRepoOnce sync.Once
	productRepo     repository.Product
)

func LoadProductRepo() repository.Product {
	productRepoOnce.Do(func() {
		productRepo = repository.NewProduct(
			infra.LoadDB(),
			infra.LoadMeilisearch(),
		)
	})

	return productRepo
}
