package auth

import (
	"errors"
	"gogo/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Age      *uint  `json:"age"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	auth := router.Group("/auth")
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
	auth.POST("/refresh", h.refresh)
}

func (h *Handler) register(c *gin.Context) {
	payload, ok := common.BindJSON[registerRequest](c)
	if !ok {
		return
	}

	result, err := h.service.Register(payload.Email, payload.Password, payload.Name, payload.Age)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyRegistered) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) login(c *gin.Context) {
	payload, ok := common.BindJSON[loginRequest](c)
	if !ok {
		return
	}

	result, err := h.service.Login(payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login user"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) refresh(c *gin.Context) {
	payload, ok := common.BindJSON[refreshRequest](c)
	if !ok {
		return
	}

	result, err := h.service.Refresh(payload.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, result)
}
