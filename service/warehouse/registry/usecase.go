package registry

import (
	"sync"
	"warehouse/usecase"
	"warehouse/usecase/grpc"
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

var (
	stockUsecaseOnce sync.Once
	stockUsecase     usecase.Stock
)

func LoadStockUsecase() usecase.Stock {
	stockUsecaseOnce.Do(func() {
		stockUsecase = usecase.NewStock(
			LoadStockRepo(),
			LoadProductRepo(),
			LoadWarehouseRepo(),
		)
	})

	return stockUsecase
}

var (
	stockGrpcOnce sync.Once
	stockGrpc     *grpc.StockGrpc
)

func LoadStockGrpc() *grpc.StockGrpc {
	stockGrpcOnce.Do(func() {
		stockGrpc = grpc.NewStock(
			LoadStockRepo(),
			LoadProductRepo(),
			LoadWarehouseRepo(),
		)
	})

	return stockGrpc
}

var (
	warehouseGrpcOnce sync.Once
	warehouseGrpc     *grpc.WarehouseGrpc
)

func LoadWarehouseGrpc() *grpc.WarehouseGrpc {
	warehouseGrpcOnce.Do(func() {
		warehouseGrpc = grpc.NewWarehouse(
			LoadWarehouseRepo(),
		)
	})

	return warehouseGrpc
}
