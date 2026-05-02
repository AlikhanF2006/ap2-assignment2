package usecase

import (
	"errors"

	"github.com/google/uuid"

	"payment-service/internal/domain"
	"payment-service/internal/publisher"
)

var (
	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	GetByOrderID(orderID string) (*domain.Payment, error)
	ListByStatus(status string) ([]*domain.Payment, error)
}

type PaymentUsecase struct {
	repo      PaymentRepository
	publisher *publisher.RabbitMQPublisher
}

func NewPaymentUsecase(repo PaymentRepository, publisher *publisher.RabbitMQPublisher) *PaymentUsecase {
	return &PaymentUsecase{
		repo:      repo,
		publisher: publisher,
	}
}

func (u *PaymentUsecase) CreatePayment(orderID string, amount int64) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}

	payment := &domain.Payment{
		ID:            uuid.NewString(),
		OrderID:       orderID,
		TransactionID: uuid.NewString(),
		Amount:        amount,
		Status:        status,
	}

	if err := u.repo.Create(payment); err != nil {
		return nil, err
	}

	if u.publisher != nil {
		event := publisher.PaymentEvent{
			EventID:       uuid.NewString(),
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			CustomerEmail: "user@example.com",
			Status:        payment.Status,
		}

		if err := u.publisher.Publish(event); err != nil {
			return nil, err
		}
	}

	return payment, nil
}

func (u *PaymentUsecase) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	payment, err := u.repo.GetByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

func (u *PaymentUsecase) ListPayments(status string) ([]*domain.Payment, error) {
	return u.repo.ListByStatus(status)
}
