package main

import (
	"context"
	"log"
	"shared"
	v1 "shop/delivery/v1"
	"shop/infra"
	"shop/registry"

	"github.com/labstack/echo/v4"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	srv := echo.New()

	shopV1 := v1.NewShop(registry.LoadShopUsecase())
	shopV1.Mount(srv.Group("/v1/shop"))

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
