package registry

import (
	"shop/usecase"
	"shop/usecase/grpc"
	"sync"
)

var (
	shopUsecaseOnce sync.Once
	shopUsecase     usecase.Shop
)

func LoadShopUsecase() usecase.Shop {
	shopUsecaseOnce.Do(func() {
		shopUsecase = usecase.NewShop(
			LoadShopRepo(),
			LoadShopWarehouseRepo(),
			LoadWarehouseRepo(),
		)
	})

	return shopUsecase
}

var (
	shopGrpcOnce sync.Once
	shopGrpc     *grpc.ShopGrpc
)

func LoadShopGrpc() *grpc.ShopGrpc {
	shopGrpcOnce.Do(func() {
		shopGrpc = grpc.NewShop(
			LoadShopRepo(),
			LoadShopWarehouseRepo(),
		)
	})

	return shopGrpc
}
