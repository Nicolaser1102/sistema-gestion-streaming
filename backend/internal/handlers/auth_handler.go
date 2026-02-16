package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "pendiente: register"})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "pendiente: login"})
}

// GET /api/auth/me (después, con JWT)
func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "pendiente: me"})
}
