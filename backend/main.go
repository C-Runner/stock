package main

import (
	"backend/config"
	"backend/handlers"
	"backend/models"
	"backend/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default config")
	}

	cfg := config.LoadConfig()

	if err := config.ConnectPostgres(cfg); err != nil {
		log.Printf("Warning: Failed to connect to PostgreSQL: %v", err)
	} else {
		defer config.DisconnectPostgres()
		config.DB.AutoMigrate(&models.Stock{}, &models.WatchlistItem{}, &models.StockDailySnapshot{})
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authMiddleware := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		cfg := config.LoadConfig()

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}

	r.GET("/health", handlers.HealthCheck)
	r.POST("/api/login", handlers.Login)
	r.POST("/api/stocks/backup", handlers.ManualBackup) // Manual backup trigger
	r.POST("/api/stocks/backup/:code", handlers.BackupSingleStock) // Backup single stock

	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/stocks", handlers.GetStocks)
		protected.POST("/stocks", handlers.CreateStock)
		protected.DELETE("/stocks/:code", handlers.DeleteStock)

		protected.GET("/stocks/quote/:code", handlers.GetStockQuote)
		protected.GET("/stocks/analysis/:code", handlers.GetStockAnalysis)
		protected.GET("/stocks/technical/:code", handlers.GetTechnicalAnalysis)

		protected.GET("/stocks/daily/:code", handlers.GetDailySnapshots)
		protected.GET("/stocks/daily/:code/:date", handlers.GetDailySnapshot)
		protected.GET("/stocks/daily/date/:date", handlers.GetAllSnapshotsByDate)

		protected.GET("/watchlist", handlers.GetWatchlist)
		protected.POST("/watchlist", handlers.AddToWatchlist)
		protected.DELETE("/watchlist/:code", handlers.RemoveFromWatchlist)

		protected.GET("/stocks/search", handlers.SearchStocks)
	}

	r.NoRoute(handlers.NotFound)

	// Start cron scheduler for daily backup at 15:30 China time
	cronScheduler := cron.New()
	// Run at 15:30 Beijing time (07:30 UTC) on weekdays (Monday-Friday)
	cronScheduler.AddFunc("0 7 * * 1-5", func() {
		log.Println("Starting scheduled daily backup")
		if err := services.BackupAllWatchlist(); err != nil {
			log.Printf("Scheduled backup failed: %v", err)
		}
	})
	cronScheduler.Start()
	defer cronScheduler.Stop()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
