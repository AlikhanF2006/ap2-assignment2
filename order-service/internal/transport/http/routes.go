package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders/:id", handler.GetOrderByID)
	router.PATCH("/orders/:id/cancel", handler.CancelOrder)
}
