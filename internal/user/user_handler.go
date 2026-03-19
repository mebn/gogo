package user

import (
	"errors"
	"gogo/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

type userRequest struct {
	Name string `json:"name" binding:"required"`
	Age  *uint  `json:"age" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	users := router.Group("/users")
	users.POST("", h.createUser)
	users.GET("/:id", h.getUser)
	users.PUT("/:id", h.updateUser)
}

func (h *Handler) createUser(c *gin.Context) {
	payload, ok := common.BindJSON[userRequest](c)
	if !ok {
		return
	}

	user, err := h.service.Create(payload.Name, *payload.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) getUser(c *gin.Context) {
	id, ok := common.ParseID(c)
	if !ok {
		return
	}

	user, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, ok := common.ParseID(c)
	if !ok {
		return
	}

	payload, ok := common.BindJSON[userRequest](c)
	if !ok {
		return
	}

	user, err := h.service.Update(id, payload.Name, *payload.Age)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
