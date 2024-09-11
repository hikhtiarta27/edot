package registry

import (
	"proto_buffer/warehouse"
	"shop/infra"
	"shop/repository"
	"sync"
)

var (
	shopRepoOnce sync.Once
	shopRepo     repository.Shop
)

func LoadShopRepo() repository.Shop {
	shopRepoOnce.Do(func() {
		shopRepo = repository.NewShop(
			infra.LoadDB(),
		)
	})

	return shopRepo
}

var (
	warehouseRepoOnce sync.Once
	warehouseRepo     repository.Warehouse
)

func LoadWarehouseRepo() repository.Warehouse {
	warehouseRepoOnce.Do(func() {
		warehouseRepo = repository.NewWarehouse(
			warehouse.NewWarehouseServiceClient(
				infra.LoadWarehouseService(),
			),
		)
	})

	return warehouseRepo
}

var (
	shopWarehouseRepoOnce sync.Once
	shopWarehouseRepo     repository.ShopWarehouse
)

func LoadShopWarehouseRepo() repository.ShopWarehouse {
	shopWarehouseRepoOnce.Do(func() {
		shopWarehouseRepo = repository.NewShopWarehouse(
			infra.LoadDB(),
		)
	})

	return shopWarehouseRepo
}
