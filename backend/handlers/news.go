package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

// GetStockNews returns latest news for a stock
func GetStockNews(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Try to get from database first (cached)
	newsItems, err := services.GetRecentNewsFromDB(code, limit)
	if err != nil {
		newsItems = []models.NewsItem{}
	}

	// If no cached news or expired (>1 hour), fetch fresh
	needsRefresh := len(newsItems) == 0
	if len(newsItems) > 0 {
		latest := newsItems[0]
		if time.Since(latest.CreatedAt) > time.Hour {
			needsRefresh = true
		}
	}

	if needsRefresh {
		// Fetch fresh news in background, return cached immediately
		go func() {
			if _, err := services.GetNewsForStock(code); err != nil {
				log.Printf("Background news fetch failed for %s: %v", code, err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"count":     len(newsItems),
		"newsItems": newsItems,
	})
}

// GetNewsSentiment returns aggregated news sentiment for AI analysis
func GetNewsSentiment(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	// Get fresh news sentiment
	result, err := services.GetNewsForStock(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news sentiment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
