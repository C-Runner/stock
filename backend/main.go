package main

import (
	"backend/config"
	"backend/handlers"
	"backend/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
		config.DB.AutoMigrate(&models.Stock{}, &models.WatchlistItem{})
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

	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/stocks", handlers.GetStocks)
		protected.POST("/stocks", handlers.CreateStock)
		protected.DELETE("/stocks/:code", handlers.DeleteStock)

		protected.GET("/stocks/quote/:code", handlers.GetStockQuote)
		protected.GET("/stocks/analysis/:code", handlers.GetStockAnalysis)
		protected.GET("/stocks/technical/:code", handlers.GetTechnicalAnalysis)

		protected.GET("/watchlist", handlers.GetWatchlist)
		protected.POST("/watchlist", handlers.AddToWatchlist)
		protected.DELETE("/watchlist/:code", handlers.RemoveFromWatchlist)

		protected.GET("/stocks/search", handlers.SearchStocks)
	}

	r.NoRoute(handlers.NotFound)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
