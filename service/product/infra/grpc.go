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
