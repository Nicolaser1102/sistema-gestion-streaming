package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"streaming-system/internal/models"
	"streaming-system/internal/utils"
)

const (
	CtxUserIDKey = "userId"
	CtxRoleKey   = "role"
)

// Valida JWT y guarda claims en el contexto
func AuthJWT(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta Authorization"})
			c.Abort()
			return
		}

		// formato: "Bearer <token>"
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization debe ser: Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(jwtSecret, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)

		c.Next()
	}
}

// Requiere rol específico (ej. ADMIN)
func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get(CtxRoleKey)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No autorizado"})
			c.Abort()
			return
		}

		currentRole, ok := val.(models.UserRole)
		if !ok || currentRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permisos insuficientes"})
			c.Abort()
			return
		}

		c.Next()
	}
}
