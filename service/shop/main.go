package main

import (
	"context"
	"log"
	"net"
	"os"
	"proto_buffer/shop"
	"shared"
	"shared/telemetry"
	v1 "shop/delivery/v1"
	"shop/infra"
	"shop/registry"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	tracerProvider := infra.LoadTraceProvider()

	srv := echo.New()
	srv.Use(telemetry.HttpOtel("shop"))

	shopV1 := v1.NewShop(registry.LoadShopUsecase())
	shopV1.Mount(srv.Group("/v1/shop"))

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			shared.GrpcUnaryParser(),
			otelgrpc.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			shared.GrpcStreamParser(),
			otelgrpc.StreamServerInterceptor(),
		),
	)

	shop.RegisterShopServiceServer(grpcSrv, registry.LoadShopGrpc())

	reflection.Register(grpcSrv)

	go func() {
		log.Printf("start server at %v\n", infra.LoadConfig().App.Address)

		err := srv.Start(infra.LoadConfig().App.Address)
		if err != nil {
			log.Fatalf("failed to start the server: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		log.Printf("start grpc server at %v\n", infra.LoadConfig().App.GrpcAddress)

		netServer, err := net.Listen("tcp", infra.LoadConfig().App.GrpcAddress)
		if err != nil {
			log.Fatalf("failed to start the tcp server: %v", err)
			os.Exit(1)
		}

		err = grpcSrv.Serve(netServer)
		if err != nil {
			log.Fatalf("failed to start the grpc server: %v", err)
			os.Exit(1)
		}
	}()

	shutdown.Wait()

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Printf("failed to shutdown the server: %v\n", err)
	}

	grpcSrv.GracefulStop()

	err = tracerProvider.Close()
	if err != nil {
		log.Fatalf("failed to shutdown tracer provider: %v", err)
	}

}
