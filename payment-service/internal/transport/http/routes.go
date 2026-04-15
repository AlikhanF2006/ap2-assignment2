package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.POST("/payments", handler.CreatePayment)
	router.GET("/payments/:order_id", handler.GetPaymentByOrderID)
	router.GET("/payments", handler.ListPayments)
}
