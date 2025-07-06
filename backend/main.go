// PlaySlate Backend API - Passwordless auth + Stripe billing for kids' drawing app
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"playslate-backend/config"
	"playslate-backend/database"
	"playslate-backend/handlers"
	"playslate-backend/middleware"
)

func main() {
	// Load .env in dev only - production uses real env vars
	if err := godotenv.Load(); err != nil {
		log.Println("Using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database with auto-migrations
	database.InitDatabase(cfg)
	db := database.GetDB()

	// Setup router with CORS
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORSOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Health check for load balancers
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "playslate-backend",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/magic-link", handlers.HandleMagicLinkRequest(db, cfg))
			auth.POST("/verify", handlers.HandleMagicLinkVerify(db, cfg))
			auth.POST("/google", handlers.HandleGoogleAuth(db, cfg))
			auth.POST("/logout", handlers.HandleLogout)
			auth.GET("/me", middleware.AuthMiddleware(cfg), handlers.HandleGetUser)
		}

		// Subscription endpoints (protected)
		subscription := v1.Group("/subscription")
		subscription.Use(middleware.AuthMiddleware(cfg))
		{
			subscription.POST("/create", handlers.HandleCreateSubscription(db, cfg))
			subscription.GET("/status", handlers.HandleGetSubscriptionStatus(db))
			subscription.POST("/cancel", handlers.HandleCancelSubscription(db, cfg))
			subscription.POST("/webhook", handlers.HandleStripeWebhook(db, cfg))
		}

		// User endpoints (protected)
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(cfg))
		{
			user.GET("/profile", handlers.HandleGetProfile(db))
			user.PUT("/profile", handlers.HandleUpdateProfile(db))
		}
	}

	// Start server
	log.Printf("Starting PlaySlate API on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}