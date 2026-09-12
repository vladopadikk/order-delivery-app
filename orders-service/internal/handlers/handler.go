package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/orders-service/internal/models"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID int64, orderIn models.OrderInput) (models.OrderResponse, error)
	GetOrderList(ctx context.Context, userID int64) (models.OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
}

type Handler struct {
	service OrderService
}

func NewHandler(service OrderService) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateOrderHandler(ctx *gin.Context) {
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

	var input models.OrderInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	order, err := h.service.CreateOrder(ctx, userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) GetOrderListHandler(ctx *gin.Context) {
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

	orders, err := h.service.GetOrderList(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func RegisterRouter(router *gin.RouterGroup, handler *Handler) {
	router.POST("/orders", handler.CreateOrderHandler)
	router.GET("/orders", handler.GetOrderListHandler)
}
