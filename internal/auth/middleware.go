package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gogo/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ContextUserKey = "currentUser"

var ErrInvalidAccessToken = errors.New("invalid access token")

func (s *Service) RequireAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		scheme, token, ok := strings.Cut(header, " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := s.parseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		var dbUser user.User
		if err := s.db.First(&dbUser, uint(userID)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate user"})
			c.Abort()
			return
		}

		c.Set(ContextUserKey, &dbUser)
		c.Next()
	}
}
