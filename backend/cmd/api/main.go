package main

import (
	"net/http"

	"streaming-system/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		authHandler := handlers.NewAuthHandler()

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.Me) // después lo protegemos con middleware JWT
		}
	}

	r.Run(":8080")
}
