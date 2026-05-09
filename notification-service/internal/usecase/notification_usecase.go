package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"notification-service/internal/domain"
	"notification-service/internal/provider"
	"notification-service/internal/store"
)

type NotificationUsecase struct {
	sender             provider.EmailSender
	idempotencyStore   store.IdempotencyStore
	maxRetries         int
	backoffBaseSeconds int
}

func NewNotificationUsecase(
	sender provider.EmailSender,
	idempotencyStore store.IdempotencyStore,
	maxRetries int,
	backoffBaseSeconds int,
) *NotificationUsecase {
	return &NotificationUsecase{
		sender:             sender,
		idempotencyStore:   idempotencyStore,
		maxRetries:         maxRetries,
		backoffBaseSeconds: backoffBaseSeconds,
	}
}

func (u *NotificationUsecase) SendNotification(event domain.PaymentEvent) error {
	ctx := context.Background()

	paymentID := event.PaymentID
	if paymentID == "" {
		paymentID = event.EventID
	}

	processed, err := u.idempotencyStore.IsProcessed(ctx, paymentID)
	if err != nil {
		return err
	}

	if processed {
		return nil
	}

	subject := fmt.Sprintf("Payment %s for Order %s", event.Status, event.OrderID)
	body := fmt.Sprintf(
		"Your payment for order %s was processed. Amount: %d. Status: %s.",
		event.OrderID,
		event.Amount,
		event.Status,
	)

	var lastErr error

	for attempt := 1; attempt <= u.maxRetries; attempt++ {
		log.Printf("[Worker] Sending notification payment_id=%s attempt=%d/%d", paymentID, attempt, u.maxRetries)

		err := u.sender.Send(ctx, event.CustomerEmail, subject, body)
		if err == nil {
			if err := u.idempotencyStore.MarkProcessed(ctx, paymentID); err != nil {
				return err
			}

			log.Printf("[Worker] Notification completed payment_id=%s", paymentID)
			return nil
		}

		lastErr = err

		if attempt < u.maxRetries {
			backoff := time.Duration(u.backoffBaseSeconds*(1<<(attempt-1))) * time.Second
			log.Printf("[Worker] Provider failed payment_id=%s error=%v retry_after=%s", paymentID, err, backoff.String())
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("notification failed after %d attempts: %w", u.maxRetries, lastErr)
}
