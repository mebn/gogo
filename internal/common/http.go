package common

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context) (*T, bool) {
	var payload T
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return nil, false
	}
	return &payload, true
}

func ParseID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return 0, false
	}

	return uint(id), true
}
