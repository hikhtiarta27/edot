package main

import (
	"context"
	"log"
	v1 "order/delivery/v1"
	"order/infra"
	"order/registry"
	"shared"

	"github.com/labstack/echo/v4"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	srv := echo.New()

	orderV1 := v1.NewOrder(registry.LoadOrderUsecase())
	orderV1.Mount(srv.Group("/v1/order", infra.LoadJWT().Validate()))

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
