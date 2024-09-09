package registry

import (
	"proto_buffer/product"
	"sync"
	"warehouse/infra"
	"warehouse/repository"
)

var (
	warehouseRepoOnce sync.Once
	warehouseRepo     repository.Warehouse
)

func LoadWarehouseRepo() repository.Warehouse {
	warehouseRepoOnce.Do(func() {
		warehouseRepo = repository.NewWarehouse(
			infra.LoadDB(),
		)
	})

	return warehouseRepo
}

var (
	warehouseTransferRepoOnce sync.Once
	warehouseTransferRepo     repository.WarehouseTransfer
)

func LoadWarehouseTransferRepo() repository.WarehouseTransfer {
	warehouseTransferRepoOnce.Do(func() {
		warehouseTransferRepo = repository.NewWarehouseTransfer(
			infra.LoadDB(),
			LoadProductRepo(),
		)
	})

	return warehouseTransferRepo
}

var (
	productRepoOnce sync.Once
	productRepo     repository.Product
)

func LoadProductRepo() repository.Product {
	productRepoOnce.Do(func() {
		productRepo = repository.NewProduct(
			product.NewProductServiceClient(
				infra.LoadProductService(),
			),
		)
	})

	return productRepo
}
