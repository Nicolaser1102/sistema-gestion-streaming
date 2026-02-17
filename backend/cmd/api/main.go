package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"streaming-system/internal/handlers"
	"streaming-system/internal/middleware"
	"streaming-system/internal/repositories"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")

	// ✅ CORS (para permitir llamadas desde Live Server: 5500)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	userRepo := repositories.NewUserRepoJSON("data/users.json")
	authHandler := handlers.NewAuthHandler(userRepo, secret)

	contentRepo := repositories.NewContentRepoJSON("data/contents.json")
	contentHandler := handlers.NewContentHandler(contentRepo)

	myListRepo := repositories.NewMyListRepoJSON("data/mylist.json")
	myListHandler := handlers.NewMyListHandler(myListRepo, contentRepo)

	subRepo := repositories.NewSubscriptionRepoJSON("data/subscriptions.json")
	progressRepo := repositories.NewProgressRepoJSON("data/progress.json")
	playbackHandler := handlers.NewPlaybackHandler(progressRepo, contentRepo)

	api := r.Group("/api")
	{
		// ✅ Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthJWT(secret), authHandler.Me)
		}

		// ✅ Catálogo (usuario)
		api.GET("/contents", contentHandler.ListContents)
		api.GET("/contents/:id", contentHandler.GetContentByID)

		// ✅ Admin (prueba de roles)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthJWT(secret))
		admin.Use(middleware.RequireRole("ADMIN"))
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin ok"})
			})
			admin.GET("/contents", contentHandler.AdminListContents)
			admin.POST("/contents", contentHandler.CreateContent)
			admin.PUT("/contents/:id", contentHandler.UpdateContent)
			admin.DELETE("/contents/:id", contentHandler.DeleteContent)

		}

		myList := api.Group("/my-list")
		myList.Use(middleware.AuthJWT(secret))
		{
			myList.GET("", myListHandler.GetMyList)
			myList.POST("/:contentId", myListHandler.AddToMyList)
			myList.DELETE("/:contentId", myListHandler.RemoveFromMyList)
		}

		play := api.Group("/playback")
		play.Use(middleware.AuthJWT(secret))
		play.Use(middleware.RequireSubscriptionActive(subRepo)) // 🔒 requiere suscripción
		{
			play.GET("/:contentId", playbackHandler.GetProgress)
			play.PUT("/:contentId/progress", playbackHandler.UpdateProgress)
		}
	}

	r.Run(":8080")
}
