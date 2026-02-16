package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"streaming-system/internal/handlers"
	"streaming-system/internal/repositories"
)

func main() {

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Inicializar repositorio
	userRepo := repositories.NewUserRepoJSON("data/users.json")

	// Inicializar handler
	authHandler := handlers.NewAuthHandler(userRepo)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
		}
	}

	r.Run(":8080")
}
