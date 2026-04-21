package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/config"
	"backend/models"
	"backend/services"
	"github.com/gin-gonic/gin"
)

// ManualBackup triggers a manual backup of all watchlist stocks
func ManualBackup(c *gin.Context) {
	go func() {
		if err := services.BackupAllWatchlist(); err != nil {
			// In production, this would log to a monitoring system
			println("Backup failed:", err.Error())
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":    "Backup started",
		"status":     "running",
		"started_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// BackupSingleStock triggers a manual backup for a single stock
func BackupSingleStock(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	if err := services.BackupStockDaily(code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Backup completed",
		"code":       code,
		"backupDate": time.Now().Format("2006-01-02"),
	})
}

// GetDailySnapshots retrieves historical daily data for a stock
func GetDailySnapshots(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "60")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 60
	}

	snapshots, err := services.GetDailySnapshots(code, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     code,
		"count":    len(snapshots),
		"snapshots": snapshots,
	})
}

// GetDailySnapshot retrieves a specific day's snapshot
func GetDailySnapshot(c *gin.Context) {
	code := c.Param("code")
	date := c.Param("date")

	if code == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code and date are required"})
		return
	}

	snapshot, err := services.GetDailySnapshot(code, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// GetAllSnapshotsByDate retrieves all snapshots for a specific date
func GetAllSnapshotsByDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
		return
	}

	var snapshots []models.StockDailySnapshot
	if err := config.DB.Where("date = ?", date).Find(&snapshots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query snapshots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":      date,
		"count":     len(snapshots),
		"snapshots": snapshots,
	})
}