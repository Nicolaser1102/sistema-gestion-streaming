package middleware

import (
	"net/http"

	"streaming-system/internal/repositories"

	"github.com/gin-gonic/gin"
)

func RequireSubscriptionActive(subRepo *repositories.SubscriptionRepoJSON) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
			c.Abort()
			return
		}
		userID, _ := userIDVal.(string)

		active, err := subRepo.IsActive(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
			c.Abort()
			return
		}
		if !active {
			c.JSON(http.StatusForbidden, gin.H{"error": "Suscripción no activa"})
			c.Abort()
			return
		}

		c.Next()
	}
}
