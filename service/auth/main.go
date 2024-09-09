package main

import (
	v1 "auth/delivery/v1"
	"auth/infra"
	"auth/registry"
	"context"
	"log"
	"shared"

	"github.com/labstack/echo/v4"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	srv := echo.New()

	accountV1 := v1.NewAccount(registry.LoadAccountUsecase())
	accountV1.Mount(srv.Group("/v1/account"))

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
