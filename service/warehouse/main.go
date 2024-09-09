package main

import (
	"context"
	"log"
	"shared"
	v1 "warehouse/delivery/v1"
	"warehouse/infra"
	"warehouse/registry"

	"github.com/labstack/echo/v4"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	srv := echo.New()

	warehouseV1 := v1.NewWarehouse(registry.LoadWarehouseUsecase())
	warehouseV1.Mount(srv.Group("/v1/warehouse", infra.LoadJWT().Validate()))

	go func() {
		shutdown.Wait()

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("failed to shutdown the server: %v", err)
		}
	}()

	log.Printf("start server at %v\n", infra.LoadConfig().App.Address)

	err := srv.Start(infra.LoadConfig().App.Address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
