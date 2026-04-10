package app

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"payment-service/internal/repository"
	httpTransport "payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

func NewApp() (*gin.Engine, *usecase.PaymentUsecase, error) {
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	repo := repository.NewPostgresRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(repo)
	handler := httpTransport.NewHandler(paymentUsecase)

	router := gin.Default()
	httpTransport.RegisterRoutes(router, handler)

	return router, paymentUsecase, nil
}
