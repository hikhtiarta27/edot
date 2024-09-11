package infra

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

var (
	shopServiceOnce sync.Once
	shopService     *grpc.ClientConn
)

func LoadShopService() *grpc.ClientConn {
	shopServiceOnce.Do(func() {

		grpcOpt := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		shopServiceConn, err := grpc.Dial(LoadConfig().Service.Shop, grpcOpt...)
		if err != nil {
			log.Fatalf("failed to connect with shop service: %v", err)
		}

		shopService = shopServiceConn
	})

	return shopService
}
