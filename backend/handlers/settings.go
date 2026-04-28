package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
)

// GetAISettings returns the current AI settings
func GetAISettings(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	var settings models.AISettings
	result := config.DB.Where("user_id = ?", userIDStr).First(&settings)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			settings = models.AISettings{
				UserID:  userIDStr,
				APIKey:  "",
				APIURL:  "https://api.minimaxi.com/v1/text/chatcompletion",
				Model:   "MiniMax-Text-01",
				Enabled: true,
			}
			if err := config.DB.Create(&settings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI settings"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI settings"})
			return
		}
	}

	// Don't return the actual API key for security
	hasAPIKey := settings.APIKey != ""
	settings.APIKey = ""
	if hasAPIKey {
		settings.APIKey = "********"
	}

	c.JSON(http.StatusOK, gin.H{
		"apiKey":  settings.APIKey,
		"apiUrl":  settings.APIURL,
		"model":   settings.Model,
		"groupId": settings.GroupID,
		"enabled": settings.Enabled,
		"hasApiKey": hasAPIKey,
	})
}

// UpdateAISettings updates the AI settings
func UpdateAISettings(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		APIKey  string `json:"apiKey"`
		APIURL  string `json:"apiUrl"`
		Model   string `json:"model"`
		GroupID string `json:"groupId"`
		Enabled bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("UpdateAISettings input: userID=%s, apiUrl=%s, model=%s, enabled=%v, hasApiKey=%v",
		userID, input.APIURL, input.Model, input.Enabled, input.APIKey != "")

	userIDStr := userID.(string)
	var settings models.AISettings
	result := config.DB.Where("user_id = ?", userIDStr).First(&settings)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			settings = models.AISettings{
				UserID:  userIDStr,
				APIKey:  "",
				APIURL:  "https://api.minimaxi.com/v1/text/chatcompletion",
				Model:   "MiniMax-Text-01",
				Enabled: true,
			}
			if err := config.DB.Create(&settings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI settings"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI settings"})
			return
		}
	}

	log.Printf("UpdateAISettings after First/Create: ID=%d, apiUrl=%s, model=%s",
		settings.ID, settings.APIURL, settings.Model)

	// Only update API key if it's not the placeholder
	if input.APIKey != "" && input.APIKey != "********" {
		settings.APIKey = input.APIKey
	}

	if input.APIURL != "" {
		settings.APIURL = input.APIURL
	}

	if input.Model != "" {
		settings.Model = input.Model
	}

	if input.GroupID != "" {
		settings.GroupID = input.GroupID
	}

	settings.Enabled = input.Enabled

	log.Printf("UpdateAISettings before Save: ID=%d, apiUrl=%s, model=%s",
		settings.ID, settings.APIURL, settings.Model)

	if err := config.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AI settings"})
		return
	}

	log.Printf("UpdateAISettings success: ID=%d, apiUrl=%s", settings.ID, settings.APIURL)

	// Return sanitized response
	hasAPIKey := settings.APIKey != ""
	settings.APIKey = ""
	if hasAPIKey {
		settings.APIKey = "********"
	}

	c.JSON(http.StatusOK, gin.H{
		"apiKey":  settings.APIKey,
		"apiUrl":  settings.APIURL,
		"model":   settings.Model,
		"groupId": settings.GroupID,
		"enabled": settings.Enabled,
		"hasApiKey": hasAPIKey,
	})
}

// TestAISettings tests the AI configuration by sending a simple request
func TestAISettings(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr := userID.(string)

	var settings models.AISettings
	result := config.DB.Where("user_id = ?", userIDStr).First(&settings)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI settings not found. Please save settings first."})
		return
	}

	if settings.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key is not configured"})
		return
	}

	// Simple connectivity test - send a minimal request to the API
	startTime := time.Now()

	// Build a minimal test prompt
	testPrompt := "Hello, please respond with exactly the word \"OK\" if you receive this message."

	// Detect API type based on URL
	isOpenAICompatible := strings.Contains(settings.APIURL, "openai") ||
		strings.Contains(settings.APIURL, "minimax") ||
		strings.Contains(settings.APIURL, "deepseek") ||
		strings.Contains(settings.APIURL, "chat/completions")
	isAnthropic := strings.Contains(settings.APIURL, "anthropic") || strings.Contains(settings.APIURL, "messages")

	var reqBody []byte
	var err error

	if isOpenAICompatible && !isAnthropic {
		reqBody, err = json.Marshal(map[string]interface{}{
			"model": settings.Model,
			"messages": []map[string]string{
				{"role": "user", "content": testPrompt},
			},
			"max_tokens": 10,
		})
	} else {
		reqBody, err = json.Marshal(map[string]interface{}{
			"model": settings.Model,
			"messages": []map[string]string{
				{"role": "user", "content": testPrompt},
			},
			"max_tokens": 10,
		})
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Failed to build request: " + err.Error(),
		})
		return
	}

	req, err := http.NewRequest("POST", settings.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Failed to create request: " + err.Error(),
		})
		return
	}

	// Set headers based on API type
	if isOpenAICompatible && !isAnthropic {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+settings.APIKey)
		// Add GroupID header for MiniMax API
		if settings.GroupID != "" {
			req.Header.Set("GroupId", settings.GroupID)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", settings.APIKey)
		req.Header.Set("anthropic-dangerous-direct-browser-access", "true")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Connection failed: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	elapsed := time.Since(startTime)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"message":     "AI connection successful",
			"model":       settings.Model,
			"statusCode":  resp.StatusCode,
			"responseTime": elapsed.String(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success":    false,
			"error":      fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)),
			"statusCode": resp.StatusCode,
		})
	}
}
