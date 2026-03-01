package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type TechnicalAnalysis = services.TechnicalAnalysis

// GetStockQuote fetches real-time stock quote from external API
func GetStockQuote(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	quote, err := services.SinaFinanceAPI(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// GetStockAnalysis fetches stock analysis with position data
func GetStockAnalysis(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	// Fetch stock from database
	collection := config.GetCollection("stocks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var stock models.Stock
	err := collection.FindOne(ctx, bson.M{"code": code}).Decode(&stock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found in portfolio"})
		return
	}

	// Fetch real-time quote from external API
	quote, err := services.SinaFinanceAPI(code)
	if err != nil {
		// If external API fails, use stored price
		quote = &models.StockQuote{
			Code:       stock.Code,
			Name:       stock.Name,
			Current:    stock.CurrentPrice,
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	// Calculate analysis
	currentPrice := quote.Current
	quantity := stock.Quantity
	buyPrice := stock.BuyPrice

	marketValue := currentPrice * float64(quantity)
	cost := buyPrice * float64(quantity)
	profitLoss := marketValue - cost

	var profitRate float64
	if cost > 0 {
		profitRate = (profitLoss / cost) * 100
	}

	// Calculate change from previous close (using stored price as previous close)
	change := 0.0
	changeAmount := 0.0
	if stock.CurrentPrice > 0 {
		change = ((currentPrice - stock.CurrentPrice) / stock.CurrentPrice) * 100
		changeAmount = currentPrice - stock.CurrentPrice
	}

	analysis := models.StockAnalysis{
		Code:         stock.Code,
		Name:         stock.Name,
		CurrentPrice: currentPrice,
		Quantity:     quantity,
		BuyPrice:     buyPrice,
		MarketValue:  marketValue,
		Cost:         cost,
		ProfitLoss:   profitLoss,
		ProfitRate:   profitRate,
		Change:       change,
		ChangeAmount: changeAmount,
	}

	c.JSON(http.StatusOK, analysis)
}

// GetTechnicalAnalysis fetches technical analysis data
func GetTechnicalAnalysis(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	analysis, err := services.GetTechnicalAnalysis(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
