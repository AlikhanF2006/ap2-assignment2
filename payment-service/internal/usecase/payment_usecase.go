package usecase

import (
	"errors"

	"github.com/google/uuid"

	"payment-service/internal/domain"
)

var (
	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	GetByOrderID(orderID string) (*domain.Payment, error)
}

type PaymentUsecase struct {
	repo PaymentRepository
}

func NewPaymentUsecase(repo PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{
		repo: repo,
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
