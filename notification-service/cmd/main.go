package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"notification-service/internal/consumer"
	"notification-service/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	notificationUsecase := usecase.NewNotificationUsecase()

	rabbitConsumer, err := consumer.NewRabbitMQConsumer(rabbitURL, notificationUsecase)
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}
	defer rabbitConsumer.Close()

	go func() {
		if err := rabbitConsumer.Start(); err != nil {
			log.Println("consumer stopped:", err)
		}
	}()

	log.Println("Notification service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Notification service shutting down...")
}
