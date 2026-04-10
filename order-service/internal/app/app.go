package app

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"order-service/internal/client"
	"order-service/internal/repository"
	grpcTransport "order-service/internal/transport/grpc"
	httpTransport "order-service/internal/transport/http"
	"order-service/internal/usecase"
)

func NewApp() (*gin.Engine, *usecase.OrderUsecase, error) {
	dbURL := os.Getenv("DB_URL")
	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	repo := repository.NewPostgresRepository(db)

	paymentClient, err := client.NewPaymentClient(paymentGRPCAddr)
	if err != nil {
		return nil, nil, err
	}

	orderUsecase := usecase.NewOrderUsecase(repo, paymentClient)
	handler := httpTransport.NewHandler(orderUsecase)

	router := gin.Default()
	httpTransport.RegisterRoutes(router, handler)

	return router, orderUsecase, nil
}

func NewGRPCHandler(u *usecase.OrderUsecase) *grpcTransport.OrderHandler {
	return grpcTransport.NewOrderHandler(u)
}
