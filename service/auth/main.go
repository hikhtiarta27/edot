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

	tracerProvider := infra.LoadTraceProvider()

	srv := echo.New()
	srv.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Everything's fine",
		})
	})

	accountV1 := v1.NewAccount(registry.LoadAccountUsecase())
	accountV1.Mount(srv.Group("/v1/account"))

	go func() {
		shutdown.Wait()

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("failed to shutdown the server: %v", err)
		}

		err = tracerProvider.Close()
		if err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()

	log.Printf("start server at %v\n", infra.LoadConfig().App.Address)

	err := srv.Start(infra.LoadConfig().App.Address)
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
