package handlers

import (
	"net/http"

	"backend/services"
	"github.com/gin-gonic/gin"
)

// GetAIAnalysis generates AI-powered stock analysis report
func GetAIAnalysis(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock code is required"})
		return
	}

	// Get quote
	quote, err := services.SinaFinanceAPI(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock quote: " + err.Error()})
		return
	}

	// Get technical analysis
	techAnalysis, err := services.GetTechnicalAnalysis(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch technical analysis: " + err.Error()})
		return
	}

	// Fetch historical data for pattern matching (need full history for patterns)
	sinaCode := services.ConvertToSinaCode(code)
	klines, err := services.FetchKLineData(sinaCode)
	if err != nil {
		klines = nil // Continue without historical data
	}

	// Build AI input - use *models.StockQuote which is aliased to services.StockQuote
	input := &services.AIAnalysisInput{
		Code:           code,
		Name:           quote.Name,
		Quote:          (*services.StockQuote)(quote),
		Technical:      techAnalysis,
		HistoricalData: klines,
	}

	// Get AI analysis (from cache or generate)
	aiService := services.GetAIAnalysisService()
	report, err := aiService.GetAIAnalysis(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI analysis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
