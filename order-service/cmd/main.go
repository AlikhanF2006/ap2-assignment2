package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"order-service/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	router, err := app.NewApp()
	if err != nil {
		log.Fatal("failed to initialize app:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("failed to run server:", err)
	}
}ы
