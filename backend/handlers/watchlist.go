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
	userID := getUserID(c)
	var items []models.UserWatchlist
	if err := config.DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []models.UserWatchlist{}
	}

	// Fetch stock names for each code
	type WatchlistResponse struct {
		Code    string    `json:"code"`
		Name    string    `json:"name"`
		AddedAt time.Time `json:"addedAt"`
	}
	result := make([]WatchlistResponse, 0, len(items))
	for _, item := range items {
		var stock models.Stock
		if err := config.DB.Select("name").Where("code = ? AND user_id = ?", item.Code, userID).First(&stock).Error; err == nil {
			result = append(result, WatchlistResponse{
				Code:    item.Code,
				Name:    stock.Name,
				AddedAt: item.AddedAt,
			})
		} else {
			// Fallback to watchlist backup table
			var wl models.WatchlistItem
			if err := config.DB.Select("name").Where("code = ?", item.Code).First(&wl).Error; err == nil {
				result = append(result, WatchlistResponse{
					Code:    item.Code,
					Name:    wl.Name,
					AddedAt: item.AddedAt,
				})
			} else {
				result = append(result, WatchlistResponse{
					Code:    item.Code,
					Name:    item.Code,
					AddedAt: item.AddedAt,
				})
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// AddToWatchlist adds a stock to user's watchlist
func AddToWatchlist(c *gin.Context) {
	userID := getUserID(c)
	var req struct {
		Code string `json:"code" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.UserWatchlist
	if err := config.DB.Where("code = ? AND user_id = ?", req.Code, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Stock already in watchlist"})
		return
	}

	item := models.UserWatchlist{
		UserID:  userID,
		Code:    req.Code,
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
	userID := getUserID(c)
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	result := config.DB.Where("code = ? AND user_id = ?", code, userID).Delete(&models.UserWatchlist{})
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
