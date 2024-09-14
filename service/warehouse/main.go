package main

import (
	"context"
	"log"
	"net"
	"os"
	"proto_buffer/stock"
	"proto_buffer/warehouse"
	"shared"
	"shared/telemetry"
	v1 "warehouse/delivery/v1"
	"warehouse/infra"
	"warehouse/registry"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	shutdown := shared.NewGracefullShutdown()

	tracerProvider := infra.LoadTraceProvider()

	srv := echo.New()
	srv.Use(telemetry.HttpOtel("warehouse"))

	warehouseV1 := v1.NewWarehouse(registry.LoadWarehouseUsecase())
	warehouseV1.Mount(srv.Group("/v1/warehouse", infra.LoadJWT().Validate()))

	stockV1 := v1.NewStock(registry.LoadStockUsecase())
	stockV1.Mount(srv.Group("/v1/stock", infra.LoadJWT().Validate()))

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

	stock.RegisterStockServiceServer(grpcSrv, registry.LoadStockGrpc())
	warehouse.RegisterWarehouseServiceServer(grpcSrv, registry.LoadWarehouseGrpc())

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
