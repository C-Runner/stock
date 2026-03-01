package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "stock-api",
	}
	c.JSON(http.StatusOK, response)
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Not Found",
	})
}

func GetStocks(c *gin.Context) {
	collection := config.GetCollection("stocks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var stocks []models.Stock
	if err := cursor.All(ctx, &stocks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stocks == nil {
		stocks = []models.Stock{}
	}

	c.JSON(http.StatusOK, stocks)
}

func CreateStock(c *gin.Context) {
	var req models.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := config.GetCollection("stocks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if stock already exists
	var existing models.Stock
	err := collection.FindOne(ctx, bson.M{"code": req.Code}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Stock with this code already exists"})
		return
	}

	stock := models.Stock{
		Code:         req.Code,
		Name:         req.Name,
		CurrentPrice: req.CurrentPrice,
		Quantity:     req.Quantity,
		BuyPrice:     req.BuyPrice,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = collection.InsertOne(ctx, stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

func DeleteStock(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	collection := config.GetCollection("stocks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock deleted successfully"})
}
