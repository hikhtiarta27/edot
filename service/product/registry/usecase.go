package registry

import (
	"product/usecase"
	"product/usecase/grpc"
	"sync"
)

var (
	productUsecaseOnce sync.Once
	productUsecase     usecase.Product
)

func LoadProductUsecase() usecase.Product {
	productUsecaseOnce.Do(func() {
		productUsecase = usecase.NewProduct(
			LoadProductRepo(),
			LoadStockRepo(),
		)
	})

	return productUsecase
}

var (
	productGrpcOnce sync.Once
	productGrpc     *grpc.ProductGrpc
)

func LoadProductGrpc() *grpc.ProductGrpc {
	productGrpcOnce.Do(func() {
		productGrpc = grpc.NewProduct(
			LoadProductRepo(),
		)
	})

	return productGrpc
}
