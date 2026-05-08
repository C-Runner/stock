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

	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	// Get news sentiment data
	newsData := prepareNewsData(code)

	// Build AI input - use *models.StockQuote which is aliased to services.StockQuote
	input := &services.AIAnalysisInput{
		Code:           code,
		Name:           quote.Name,
		Quote:          (*services.StockQuote)(quote),
		Technical:      techAnalysis,
		HistoricalData: klines,
		NewsData:       newsData,
	}

	// Get AI analysis (from cache or generate)
	aiService := services.GetAIAnalysisService()
	report, err := aiService.GetAIAnalysis(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI analysis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// prepareNewsData fetches and prepares news sentiment for AI
func prepareNewsData(code string) *services.NewsSentimentData {
	result, err := services.GetNewsForStock(code)
	if err != nil {
		return nil
	}

	newsItems, err := services.GetRecentNewsFromDB(code, 5)
	if err != nil {
		return nil
	}

	var recentNews []services.NewsSummary
	for _, item := range newsItems {
		if item.PublishTime.IsZero() {
			continue
		}
		recentNews = append(recentNews, services.NewsSummary{
			Title:       item.Title,
			PublishTime: item.PublishTime.Format("2006-01-02 15:04"),
			Sentiment:   item.Sentiment,
		})
	}

	return &services.NewsSentimentData{
		OverallScore:   result.OverallScore,
		PositiveCount:  result.PositiveCount,
		NeutralCount:   result.NeutralCount,
		NegativeCount:  result.NegativeCount,
		LatestNewsTime: result.LatestNewsTime,
		RecentNews:     recentNews,
	}
}
