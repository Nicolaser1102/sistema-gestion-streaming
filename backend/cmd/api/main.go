package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"streaming-system/internal/handlers"
	"streaming-system/internal/middleware"
	"streaming-system/internal/repositories"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	userRepo := repositories.NewUserRepoJSON("data/users.json")
	authHandler := handlers.NewAuthHandler(userRepo, secret)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// protegido con JWT
			auth.GET("/me", middleware.AuthJWT(secret), authHandler.Me)
		}

		// ejemplo de ruta admin (para probar roles)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthJWT(secret))
		admin.Use(middleware.RequireRole("ADMIN"))
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin ok"})
			})
		}
	}

	r.Run(":8080")
}
