package registry

import (
	"product/infra"
	"product/repository"
	"proto_buffer/shop"
	"proto_buffer/stock"
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

var (
	stockRepoOnce sync.Once
	stockRepo     repository.Stock
)

func LoadStockRepo() repository.Stock {
	stockRepoOnce.Do(func() {
		stockRepo = repository.NewStock(
			stock.NewStockServiceClient(
				infra.LoadWarehouseService(),
			),
		)
	})

	return stockRepo
}

var (
	shopRepoOnce sync.Once
	shopRepo     repository.Shop
)

func LoadShopRepo() repository.Shop {
	shopRepoOnce.Do(func() {
		shopRepo = repository.NewShop(
			shop.NewShopServiceClient(
				infra.LoadShopService(),
			),
		)
	})

	return shopRepo
}
