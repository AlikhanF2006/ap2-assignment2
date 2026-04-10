package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

type PaymentUsecase interface {
	CreatePayment(orderID string, amount int64) (*domain.Payment, error)
	GetPaymentByOrderID(orderID string) (*domain.Payment, error)
}

type Handler struct {
	usecase PaymentUsecase
}

func NewHandler(usecase PaymentUsecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

type CreatePaymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	payment, err := h.usecase.CreatePayment(req.OrderID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create payment",
			})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *Handler) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")

	payment, err := h.usecase.GetPaymentByOrderID(orderID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get payment",
			})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}
