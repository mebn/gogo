package pet

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

type petRequest struct {
	OwnerID *uint  `json:"ownerId" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Age     *uint  `json:"age" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	pets := router.Group("/pets")
	pets.POST("", h.createPet)
	pets.GET("/:id", h.getPet)
	pets.PUT("/:id", h.updatePet)

	users := router.Group("/users")
	users.GET("/:id/pets", h.getPetsByOwner)
}

func (h *Handler) createPet(c *gin.Context) {
	payload, ok := common.BindJSON[petRequest](c)
	if !ok {
		return
	}

	pet, err := h.service.Create(*payload.OwnerID, payload.Name, *payload.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pet"})
		return
	}

	c.JSON(http.StatusCreated, pet)
}

func (h *Handler) getPet(c *gin.Context) {
	id, ok := common.ParseID(c)
	if !ok {
		return
	}

	pet, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pet"})
		return
	}

	c.JSON(http.StatusOK, pet)
}

func (h *Handler) updatePet(c *gin.Context) {
	id, ok := common.ParseID(c)
	if !ok {
		return
	}

	payload, ok := common.BindJSON[petRequest](c)
	if !ok {
		return
	}

	pet, err := h.service.Update(id, *payload.OwnerID, payload.Name, *payload.Age)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pet"})
		return
	}

	c.JSON(http.StatusOK, pet)
}

func (h *Handler) getPetsByOwner(c *gin.Context) {
	id, ok := common.ParseID(c)
	if !ok {
		return
	}

	pets, err := h.service.ListByOwner(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pets"})
		return
	}

	c.JSON(http.StatusOK, pets)
}
