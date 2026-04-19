package handlers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/services"
	"github.com/gin-gonic/gin"
)

// GetWatchlist returns all stocks in user's watchlist
func GetWatchlist(c *gin.Context) {
	var items []models.WatchlistItem
	if err := config.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []models.WatchlistItem{}
	}

	c.JSON(http.StatusOK, items)
}

// AddToWatchlist adds a stock to user's watchlist
func AddToWatchlist(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.WatchlistItem
	if err := config.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Stock already in watchlist"})
		return
	}

	item := models.WatchlistItem{
		Code:    req.Code,
		Name:    req.Name,
		AddedAt: time.Now(),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// RemoveFromWatchlist removes a stock from watchlist
func RemoveFromWatchlist(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	result := config.DB.Where("code = ?", code).Delete(&models.WatchlistItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found in watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed from watchlist"})
}

// SearchStocks searches for stocks by code or name
func SearchStocks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	quote, err := services.SinaFinanceAPI(query)
	if err == nil && quote.Name != "" {
		c.JSON(http.StatusOK, []models.StockQuote{*quote})
		return
	}

	var stocks []models.Stock
	if err := config.DB.Where("code ILIKE ?", query+"%").Or("name ILIKE ?", "%"+query+"%").Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stocks == nil {
		stocks = []models.Stock{}
	}

	c.JSON(http.StatusOK, stocks)
}
