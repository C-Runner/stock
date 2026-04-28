package handlers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
)

// getUserID extracts userID from JWT context
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return userID.(string)
	}
	return ""
}

func HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "stock-api",
	}

	dbStatus := "connected"
	if err := config.PingPostgres(); err != nil {
		dbStatus = "disconnected"
		response.Status = "degraded"
	}

	response.Database = dbStatus
	c.JSON(http.StatusOK, response)
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Not Found",
	})
}

func GetStocks(c *gin.Context) {
	userID := getUserID(c)
	var stocks []models.Stock
	if err := config.DB.Where("user_id = ?", userID).Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stocks == nil {
		stocks = []models.Stock{}
	}

	c.JSON(http.StatusOK, stocks)
}

func CreateStock(c *gin.Context) {
	userID := getUserID(c)
	var req models.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Stock
	if err := config.DB.Where("code = ? AND user_id = ?", req.Code, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Stock with this code already exists"})
		return
	}

	stock := models.Stock{
		Code:         req.Code,
		UserID:       userID,
		Name:         req.Name,
		CurrentPrice: req.CurrentPrice,
		Quantity:     req.Quantity,
		BuyPrice:     req.BuyPrice,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-add to watchlist if not already present
	var existingWL models.UserWatchlist
	if err := config.DB.Where("code = ? AND user_id = ?", req.Code, userID).First(&existingWL).Error; err != nil {
		// Not in watchlist, add it
		wlItem := models.UserWatchlist{
			UserID:  userID,
			Code:    req.Code,
			AddedAt: time.Now(),
		}
		config.DB.Create(&wlItem)
	}

	c.JSON(http.StatusCreated, stock)
}

func UpdateStock(c *gin.Context) {
	userID := getUserID(c)
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	var req models.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stock models.Stock
	if err := config.DB.Where("code = ? AND user_id = ?", code, userID).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	stock.Name = req.Name
	stock.CurrentPrice = req.CurrentPrice
	stock.Quantity = req.Quantity
	stock.BuyPrice = req.BuyPrice
	stock.UpdatedAt = time.Now()

	if err := config.DB.Save(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func DeleteStock(c *gin.Context) {
	userID := getUserID(c)
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	result := config.DB.Where("code = ? AND user_id = ?", code, userID).Delete(&models.Stock{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock deleted successfully"})
}
