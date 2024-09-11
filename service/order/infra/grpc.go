package infra

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	productServiceOnce sync.Once
	productService     *grpc.ClientConn
)

func LoadProductService() *grpc.ClientConn {
	productServiceOnce.Do(func() {

		grpcOpt := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		productServiceConn, err := grpc.Dial(LoadConfig().Service.Product, grpcOpt...)
		if err != nil {
			log.Fatalf("failed to connect with product service: %v", productServiceConn)
		}

		productService = productServiceConn
	})

	return productService
}

var (
	warehouseServiceOnce sync.Once
	warehouseService     *grpc.ClientConn
)

func LoadWarehouseService() *grpc.ClientConn {
	warehouseServiceOnce.Do(func() {

		grpcOpt := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		warehouseServiceConn, err := grpc.Dial(LoadConfig().Service.Warehouse, grpcOpt...)
		if err != nil {
			log.Fatalf("failed to connect with warehouse service: %v", err)
		}

		warehouseService = warehouseServiceConn
	})

	return warehouseService
}
