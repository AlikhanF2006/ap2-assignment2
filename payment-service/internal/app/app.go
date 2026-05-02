package app

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"payment-service/internal/publisher"
	"payment-service/internal/repository"
	httpTransport "payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

func NewApp() (*gin.Engine, *usecase.PaymentUsecase, error) {
	dbURL := os.Getenv("DB_URL")

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	repo := repository.NewPostgresRepository(db)

	rabbitPublisher, err := publisher.NewRabbitMQPublisher(rabbitURL)
	if err != nil {
		return nil, nil, err
	}

	paymentUsecase := usecase.NewPaymentUsecase(repo, rabbitPublisher)
	handler := httpTransport.NewHandler(paymentUsecase)

	router := gin.Default()
	httpTransport.RegisterRoutes(router, handler)

	return router, paymentUsecase, nil
}
