package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	paymentpb "github.com/AlikhanF2006/ap2-protos-gen/payment"
	"payment-service/internal/app"
	grpcTransport "payment-service/internal/transport/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	router, paymentUsecase, err := app.NewApp()
	if err != nil {
		log.Fatal("failed to initialize app: ", err)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("failed to listen for gRPC: ", err)
		}

		grpcServer := grpc.NewServer()
		grpcHandler := grpcTransport.NewPaymentHandler(paymentUsecase)

		paymentpb.RegisterPaymentServiceServer(grpcServer, grpcHandler)

		log.Println("gRPC payment-service running on port:", grpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to run gRPC server: ", err)
		}
	}()

	log.Println("HTTP payment-service running on port:", httpPort)

	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal("failed to run HTTP server: ", err)
	}
}
