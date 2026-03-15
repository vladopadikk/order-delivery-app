package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/models"
	"github.com/vladopadikk/order-delivery-app/auth-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var userIn models.UserInput

	if err := ctx.ShouldBindJSON(&userIn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	user, err := h.service.Register(ctx.Request.Context(), userIn)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *Handler) LoginHandler(ctx *gin.Context) {
	var loginIn models.LoginInput

	if err := ctx.ShouldBindJSON(&loginIn); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	tokens, err := h.service.Login(ctx.Request.Context(), loginIn)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/register", handler.RegisterHandler)
	router.POST("/login", handler.LoginHandler)
}
