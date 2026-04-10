package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"order-service/internal/domain"
	"order-service/internal/usecase"
)

type OrderUsecase interface {
	CreateOrder(customerID, itemName string, amount int64) (*domain.Order, error)
	GetOrderByID(id string) (*domain.Order, error)
	CancelOrder(id string) (*domain.Order, error)
}

type Handler struct {
	usecase OrderUsecase
}

func NewHandler(usecase OrderUsecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	ItemName   string `json:"item_name"`
	Amount     int64  `json:"amount"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	order, err := h.usecase.CreateOrder(req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, usecase.ErrPaymentUnavailable):
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create order",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.usecase.GetOrderByID(id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get order",
			})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.usecase.CancelOrder(id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, usecase.ErrCannotCancelOrder):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to cancel order",
			})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}
