package registry

import (
	"sync"
	"warehouse/usecase"
)

var (
	warehouseUsecaseOnce sync.Once
	warehouseUsecase     usecase.Warehouse
)

func LoadWarehouseUsecase() usecase.Warehouse {
	warehouseUsecaseOnce.Do(func() {
		warehouseUsecase = usecase.NewWarehouse(
			LoadWarehouseRepo(),
			LoadProductRepo(),
			LoadWarehouseTransferRepo(),
		)
	})

	return warehouseUsecase
}
