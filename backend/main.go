package main

import (
	"backend/config"
	"backend/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default config")
	}

	cfg := config.LoadConfig()

	if err := config.ConnectMongoDB(cfg); err != nil {
		log.Printf("Warning: Failed to connect to MongoDB: %v", err)
	} else {
		defer config.DisconnectMongoDB()
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthCheck)

	// Stock routes
	r.GET("/api/stocks", handlers.GetStocks)
	r.POST("/api/stocks", handlers.CreateStock)
	r.DELETE("/api/stocks/:code", handlers.DeleteStock)

	// Stock analysis routes
	r.GET("/api/stocks/quote/:code", handlers.GetStockQuote)
	r.GET("/api/stocks/analysis/:code", handlers.GetStockAnalysis)
	r.GET("/api/stocks/technical/:code", handlers.GetTechnicalAnalysis)

	r.NoRoute(handlers.NotFound)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
