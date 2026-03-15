package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/delivery-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetStatusHandler(ctx *gin.Context) {
	userIDAny, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID, ok := userIDAny.(int64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	orderIDStr := ctx.Query("order_id")
	if orderIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "miss order id param"})
		return
	}
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	orderStatus, err := h.service.GetDeliveryStatus(ctx, userID, orderID)
	if err != nil {
		if errors.Is(err, service.ErrDeliveryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orderStatus)
}

func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	router.GET("/deliveries", h.GetStatusHandler)
}
