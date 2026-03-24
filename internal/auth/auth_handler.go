package auth

import (
	"gogo/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password *uint  `json:"password" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	auth := router.Group("/auth")
	auth.POST("login", h.login)
}

func (h *Handler) login(c *gin.Context) {
	payload, ok := common.BindJSON[loginRequest](c)
	if !ok {
		return
	}

	user, err := h.service.Login(payload.Email, *payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
