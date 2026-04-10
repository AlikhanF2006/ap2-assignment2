package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"order-service/internal/app"

	orderpb "github.com/AlikhanF2006/ap2-protos-gen/order"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	router, usecase, err := app.NewApp()
	if err != nil {
		log.Fatal("failed to initialize app:", err)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	grpcHandler := app.NewGRPCHandler(usecase)

	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("failed to listen:", err)
		}

		grpcServer := grpc.NewServer()
		orderpb.RegisterOrderServiceServer(grpcServer, grpcHandler)

		log.Println("gRPC server running on port:", grpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC:", err)
		}
	}()

	log.Println("HTTP server running on port:", httpPort)

	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
