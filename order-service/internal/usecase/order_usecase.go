package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"order-service/internal/domain"
)

var (
	ErrInvalidAmount      = errors.New("amount must be greater than 0")
	ErrOrderNotFound      = errors.New("order not found")
	ErrCannotCancelOrder  = errors.New("only pending orders can be cancelled")
	ErrPaymentUnavailable = errors.New("payment service unavailable")
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status string) error
}

type PaymentClient interface {
	CreatePayment(orderID string, amount int64) (string, error)
}

type OrderCache interface {
	Get(ctx context.Context, id string) (*domain.Order, error)
	Set(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id string) error
}

type OrderUsecase struct {
	repo          OrderRepository
	paymentClient PaymentClient
	cache         OrderCache
}

func NewOrderUsecase(repo OrderRepository, paymentClient PaymentClient, cache OrderCache) *OrderUsecase {
	return &OrderUsecase{
		repo:          repo,
		paymentClient: paymentClient,
		cache:         cache,
	}
}

func (u *OrderUsecase) CreateOrder(customerID, itemName string, amount int64) (*domain.Order, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	order := &domain.Order{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Create(order); err != nil {
		return nil, err
	}

	paymentStatus, err := u.paymentClient.CreatePayment(order.ID, order.Amount)
	if err != nil {
		_ = u.repo.UpdateStatus(order.ID, "Failed")
		order.Status = "Failed"

		if u.cache != nil {
			if err := u.cache.Delete(context.Background(), order.ID); err != nil {
				log.Printf("failed to invalidate order cache id=%s error=%v", order.ID, err)
			}
		}

		return nil, ErrPaymentUnavailable
	}

	if paymentStatus == "Authorized" {
		order.Status = "Paid"
	} else {
		order.Status = "Failed"
	}

	if err := u.repo.UpdateStatus(order.ID, order.Status); err != nil {
		return nil, err
	}

	if u.cache != nil {
		if err := u.cache.Delete(context.Background(), order.ID); err != nil {
			log.Printf("failed to invalidate order cache id=%s error=%v", order.ID, err)
		}
	}

	return order, nil
}

func (u *OrderUsecase) GetOrderByID(id string) (*domain.Order, error) {
	ctx := context.Background()

	if u.cache != nil {
		cachedOrder, err := u.cache.Get(ctx, id)
		if err == nil {
			return cachedOrder, nil
		}
	}

	order, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	if u.cache != nil {
		if err := u.cache.Set(ctx, order); err != nil {
			log.Printf("failed to set order cache id=%s error=%v", id, err)
		}
	}

	return order, nil
}

func (u *OrderUsecase) CancelOrder(id string) (*domain.Order, error) {
	order, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != "Pending" {
		return nil, ErrCannotCancelOrder
	}

	order.Status = "Cancelled"

	if err := u.repo.UpdateStatus(id, "Cancelled"); err != nil {
		return nil, err
	}

	if u.cache != nil {
		if err := u.cache.Delete(context.Background(), id); err != nil {
			log.Printf("failed to invalidate order cache id=%s error=%v", id, err)
		}
	}

	return order, nil
}
