package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"notification-service/internal/consumer"
	"notification-service/internal/provider"
	"notification-service/internal/store"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	maxRetries, _ := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if maxRetries == 0 {
		maxRetries = 3
	}

	backoffBaseSeconds, _ := strconv.Atoi(os.Getenv("BACKOFF_BASE_SECONDS"))
	if backoffBaseSeconds == 0 {
		backoffBaseSeconds = 2
	}

	providerLatencyMS, _ := strconv.Atoi(os.Getenv("PROVIDER_LATENCY_MS"))
	if providerLatencyMS == 0 {
		providerLatencyMS = 1000
	}

	providerFailureRate, _ := strconv.Atoi(os.Getenv("PROVIDER_FAILURE_RATE"))
	if providerFailureRate == 0 {
		providerFailureRate = 30
	}

	providerMode := os.Getenv("PROVIDER_MODE")
	if providerMode == "" {
		providerMode = "SIMULATED"
	}

	var emailSender provider.EmailSender

	switch providerMode {
	case "SIMULATED":
		emailSender = provider.NewSimulatedEmailSender(providerLatencyMS, providerFailureRate)
	default:
		log.Printf("Unknown PROVIDER_MODE=%s, using SIMULATED provider", providerMode)
		emailSender = provider.NewSimulatedEmailSender(providerLatencyMS, providerFailureRate)
	}

	idempotencyStore := store.NewRedisIdempotencyStore(
		redisAddr,
		redisPassword,
		redisDB,
		24*time.Hour,
	)

	notificationUsecase := usecase.NewNotificationUsecase(
		emailSender,
		idempotencyStore,
		maxRetries,
		backoffBaseSeconds,
	)

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

	log.Printf(
		"Notification service started provider_mode=%s max_retries=%d backoff_base_seconds=%d",
		providerMode,
		maxRetries,
		backoffBaseSeconds,
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Notification service shutting down...")
}
