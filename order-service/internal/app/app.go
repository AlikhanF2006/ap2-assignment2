package app

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"order-service/internal/client"
	"order-service/internal/repository"
	httpTransport "order-service/internal/transport/http"
	"order-service/internal/usecase"
)

func NewApp() (*gin.Engine, error) {
	dbURL := os.Getenv("DB_URL")
	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := repository.NewPostgresRepository(db)
	paymentClient := client.NewPaymentClient(paymentServiceURL)
	orderUsecase := usecase.NewOrderUsecase(repo, paymentClient)
	handler := httpTransport.NewHandler(orderUsecase)

	router := gin.Default()
	httpTransport.RegisterRoutes(router, handler)

	return router, nil
}
