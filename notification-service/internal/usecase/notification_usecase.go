package usecase

import (
	"fmt"
	"sync"

	"notification-service/internal/domain"
)

type NotificationUsecase struct {
	processed map[string]bool
	mu        sync.Mutex
}

func NewNotificationUsecase() *NotificationUsecase {
	return &NotificationUsecase{
		processed: make(map[string]bool),
	}
}

func (u *NotificationUsecase) SendNotification(event domain.PaymentEvent) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.processed[event.EventID] {
		fmt.Println("[Notification] Duplicate event skipped:", event.EventID)
		return nil
	}

	fmt.Printf(
		"[Notification] Sent email to %s for Order #%s. Amount: %d. Status: %s\n",
		event.CustomerEmail,
		event.OrderID,
		event.Amount,
		event.Status,
	)

	u.processed[event.EventID] = true

	return nil
}
