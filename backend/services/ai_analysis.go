package services

import (
	"backend/config"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Cache entry with TTL
type cacheEntry struct {
	report *AIAnalysisReport
	expiry time.Time
}

// AIAnalysisService handles AI-powered stock analysis
type AIAnalysisService struct {
	cache    map[string]cacheEntry
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

var aiService *AIAnalysisService
var aiServiceOnce sync.Once

// GetAIAnalysisService returns the singleton AI analysis service
func GetAIAnalysisService() *AIAnalysisService {
	aiServiceOnce.Do(func() {
		aiService = &AIAnalysisService{
			cache:    make(map[string]cacheEntry),
			cacheTTL: 5 * time.Minute,
		}
	})
	return aiService
}

// getAISettings retrieves AI settings from database (for handlers, uses userID)
func GetAISettings(userID string) *models.AISettings {
	var settings models.AISettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err.Error() == "record not found" {
			// Create default settings for user
			settings = models.AISettings{
				UserID:  userID,
				APIKey:  "",
				APIURL:  "https://api.minimaxi.com/anthropic/v1/messages",
				Model:   "MiniMax-M2.7",
				GroupID: "",
				Enabled: true,
			}
			if err := config.DB.Create(&settings).Error; err != nil {
				log.Printf("Failed to create AI settings for user %s: %v", userID, err)
				return nil
			}
		} else {
			log.Printf("Failed to get AI settings for user %s: %v", userID, err)
			return nil
		}
	}
	return &settings
}

// GetAIAnalysis generates or retrieves cached AI analysis report
func (s *AIAnalysisService) GetAIAnalysis(input *AIAnalysisInput, userID string) (*AIAnalysisReport, error) {
	cacheKey := input.Code

	// Check cache first
	s.cacheMu.RLock()
	if entry, ok := s.cache[cacheKey]; ok && time.Now().Before(entry.expiry) {
		entry.report.FromCache = true
		s.cacheMu.RUnlock()
		log.Printf("AI analysis cache hit for %s", input.Code)
		return entry.report, nil
	}
	s.cacheMu.RUnlock()

	// Get technical analysis data
	technicalData, err := s.prepareTechnicalData(input)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare technical data: %w", err)
	}

	// Get AI settings from database
	settings := GetAISettings(userID)

	// Generate AI analysis
	var report *AIAnalysisReport

	if settings != nil && settings.APIKey != "" && settings.Enabled {
		// Use configured AI API
		report, err = s.generateWithClaude(input, technicalData, settings.APIKey, settings.APIURL, settings.Model, settings.GroupID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate AI analysis: %w", err)
		}
	} else {
		// Fallback to heuristic analysis
		log.Printf("No AI API key configured or AI disabled, using heuristic analysis")
		report = s.generateHeuristicAnalysis(input, technicalData)
	}

	report.FromCache = false
	report.Cached = false
	report.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")

	// Find similar patterns
	report.SimilarPatterns = s.findSimilarPatterns(input.HistoricalData)

	// Get institutional sentiment
	report.InstitutionalSentiment = s.getInstitutionalSentiment(input.Code)

	// Cache the result
	s.cacheMu.Lock()
	s.cache[cacheKey] = cacheEntry{
		report: report,
		expiry: time.Now().Add(s.cacheTTL),
	}
	s.cacheMu.Unlock()

	// Clean expired entries
	s.cleanExpiredCache()

	return report, nil
}

// prepareTechnicalData extracts key technical indicators for AI analysis
func (s *AIAnalysisService) prepareTechnicalData(input *AIAnalysisInput) (map[string]interface{}, error) {
	if input.Technical == nil {
		return nil, fmt.Errorf("technical analysis data is required")
	}

	ta := input.Technical
	data := make(map[string]interface{})

	// Recent price trend
	if len(ta.RecentPrices) >= 5 {
		recentPrices := ta.RecentPrices[len(ta.RecentPrices)-5:]
		data["recentTrend"] = s.calculateTrend(recentPrices)
	}

	// MA status
	if len(ta.MA) > 0 {
		data["maStatus"] = s.getMAStatus(ta)
	}

	// RSI status
	if len(ta.RSI) > 0 {
		data["rsiStatus"] = s.getRSIStatus(ta)
	}

	// MACD status
	data["macdStatus"] = s.getMACDStatus(ta)

	// KDJ status
	data["kdjStatus"] = s.getKDJStatus(ta)

	// BOLL status
	data["bollStatus"] = s.getBOLLStatus(ta)

	// Volume analysis
	if len(ta.RecentPrices) >= 20 {
		data["volumeAnalysis"] = s.analyzeVolume(ta.RecentPrices)
	}

	// Support/Resistance
	if ta.Levels.Resistance != nil && len(ta.Levels.Resistance) > 0 {
		data["resistance"] = ta.Levels.Resistance[0].Price
	}
	if ta.Levels.Support != nil && len(ta.Levels.Support) > 0 {
		data["support"] = ta.Levels.Support[0].Price
	}

	// Pattern signals
	data["patterns"] = s.getPatternSignals(ta)

	// Recommendation
	if ta.Recommendation.Action != "" {
		data["recommendation"] = ta.Recommendation.Action
		data["confidence"] = ta.Recommendation.Confidence
	}

	return data, nil
}

func (s *AIAnalysisService) calculateTrend(prices []PricePoint) string {
	if len(prices) < 2 {
		return "neutral"
	}
	first := prices[0].Close
	last := prices[len(prices)-1].Close
	change := (last - first) / first * 100
	if change > 3 {
		return "strong_up"
	} else if change > 1 {
		return "up"
	} else if change < -3 {
		return "strong_down"
	} else if change < -1 {
		return "down"
	}
	return "neutral"
}

func (s *AIAnalysisService) getMAStatus(ta *TechnicalAnalysis) string {
	if len(ta.MA) < 4 {
		return "insufficient_data"
	}
	closes := make([]float64, len(ta.RecentPrices))
	for i, p := range ta.RecentPrices {
		closes[i] = p.Close
	}
	if len(closes) < 20 {
		return "insufficient_data"
	}

	currentPrice := closes[len(closes)-1]
	ma5 := ta.MA[0].Values[len(ta.MA[0].Values)-1]
	ma10 := ta.MA[1].Values[len(ta.MA[1].Values)-1]
	ma20 := ta.MA[2].Values[len(ta.MA[2].Values)-1]
	ma60 := ta.MA[3].Values[len(ta.MA[3].Values)-1]

	if currentPrice > ma5 && ma5 > ma10 && ma10 > ma20 && ma20 > ma60 {
		return "strong_bullish"
	} else if currentPrice > ma5 && ma5 > ma10 {
		return "bullish"
	} else if currentPrice < ma5 && ma5 < ma10 && ma10 < ma20 && ma20 < ma60 {
		return "strong_bearish"
	} else if currentPrice < ma5 && ma5 < ma10 {
		return "bearish"
	}
	return "mixed"
}

func (s *AIAnalysisService) getRSIStatus(ta *TechnicalAnalysis) string {
	if len(ta.RSI) == 0 {
		return "insufficient_data"
	}
	rsi6 := ta.RSI[0].Values[len(ta.RSI[0].Values)-1]
	if rsi6 > 70 {
		return "overbought"
	} else if rsi6 < 30 {
		return "oversold"
	}
	return "neutral"
}

func (s *AIAnalysisService) getMACDStatus(ta *TechnicalAnalysis) string {
	if len(ta.MACD.DIF) == 0 || len(ta.MACD.DEA) == 0 {
		return "insufficient_data"
	}
	dif := ta.MACD.DIF[len(ta.MACD.DIF)-1]
	dea := ta.MACD.DEA[len(ta.MACD.DEA)-1]
	macd := ta.MACD.MACD[len(ta.MACD.MACD)-1]

	if dif > dea && macd > 0 {
		return "strong_bullish"
	} else if dif > dea {
		return "bullish"
	} else if dif < dea && macd < 0 {
		return "strong_bearish"
	} else if dif < dea {
		return "bearish"
	}
	return "neutral"
}

func (s *AIAnalysisService) getKDJStatus(ta *TechnicalAnalysis) string {
	if len(ta.KDJ.K) == 0 {
		return "insufficient_data"
	}
	k := ta.KDJ.K[len(ta.KDJ.K)-1]
	d := ta.KDJ.D[len(ta.KDJ.D)-1]
	j := ta.KDJ.J[len(ta.KDJ.J)-1]

	if k > 80 || j > 100 {
		return "overbought"
	} else if k < 20 || j < 0 {
		return "oversold"
	} else if k > d {
		return "bullish_cross"
	} else if k < d {
		return "bearish_cross"
	}
	return "neutral"
}

func (s *AIAnalysisService) getBOLLStatus(ta *TechnicalAnalysis) string {
	if len(ta.BOLL.Upper) == 0 || len(ta.RecentPrices) == 0 {
		return "insufficient_data"
	}
	currentPrice := ta.RecentPrices[len(ta.RecentPrices)-1].Close
	upper := ta.BOLL.Upper[len(ta.BOLL.Upper)-1]
	lower := ta.BOLL.Lower[len(ta.BOLL.Lower)-1]
	mid := ta.BOLL.Mid[len(ta.BOLL.Mid)-1]

	if currentPrice > upper {
		return "above_upper_band"
	} else if currentPrice < lower {
		return "below_lower_band"
	} else if currentPrice > mid {
		return "above_midline"
	}
	return "below_midline"
}

func (s *AIAnalysisService) analyzeVolume(prices []PricePoint) string {
	if len(prices) < 20 {
		return "insufficient_data"
	}
	var totalVolume int64
	for _, p := range prices[len(prices)-20:] {
		totalVolume += p.Volume
	}
	avgVolume := totalVolume / 20

	recentVolume := prices[len(prices)-1].Volume
	if recentVolume > avgVolume*150 {
		return "very_high"
	} else if recentVolume > avgVolume*120 {
		return "high"
	} else if recentVolume < avgVolume*50 {
		return "low"
	}
	return "normal"
}

func (s *AIAnalysisService) getPatternSignals(ta *TechnicalAnalysis) []string {
	var signals []string

	for _, pattern := range ta.Patterns.CandlestickPatterns {
		if pattern.Strength == "strong" {
			if pattern.Bullish {
				signals = append(signals, "strong_bullish_"+pattern.Type)
			} else {
				signals = append(signals, "strong_bearish_"+pattern.Type)
			}
		}
	}

	for _, pattern := range ta.Patterns.TrendPatterns {
		if pattern.Strength == "strong" {
			if pattern.Bullish {
				signals = append(signals, "strong_bullish_"+pattern.Type)
			} else {
				signals = append(signals, "strong_bearish_"+pattern.Type)
			}
		}
	}

	return signals
}

// generateWithClaude calls AI API for analysis (supports OpenAI-compatible and Anthropic formats)
func (s *AIAnalysisService) generateWithClaude(input *AIAnalysisInput, techData map[string]interface{}, apiKey, apiURL, model, groupID string) (*AIAnalysisReport, error) {
	prompt := s.buildAnalysisPrompt(input, techData)
	systemPrompt := `你是一位专业的股票分析师。请对给定股票进行全面深入的技术分析，并输出结构化的分析报告。

分析要求：
1. 技术面分析：结合K线形态、均线系统、MACD、KDJ、布林带等技术指标
2. 趋势判断：判断当前趋势（上升/下降/震荡）及强度
3. 支撑阻力：识别关键支撑位和阻力位
4. 形态识别：识别常见的K线形态和趋势形态
5. 动量分析：分析涨跌动量变化
6. 成交量分析：分析量价关系
7. 风险提示：识别潜在风险因素
8. 投资建议：给出明确的操作建议（买入/卖出/持有）和风险等级

请用中文输出详细的分析报告，重点突出3个亮点和2个风险点。`

	log.Printf("[AI Analysis] Starting request for stock: %s (%s), Model: %s, API: %s", input.Code, input.Name, model, apiURL)

	// Detect API type based on URL
	isOpenAICompatible := strings.Contains(apiURL, "openai") || strings.Contains(apiURL, "deepseek") || strings.Contains(apiURL, "chat/completions")
	isMiniMaxAnthropic := strings.Contains(apiURL, "minimax") && strings.Contains(apiURL, "anthropic")

	var reqBody []byte
	var err error

	if isOpenAICompatible {
		// OpenAI-compatible format (DeepSeek, etc.)
		reqBody, err = s.buildOpenAIRequest(model, systemPrompt, prompt)
	} else if isMiniMaxAnthropic {
		// MiniMax Anthropic-compatible format
		reqBody, err = s.buildMiniMaxAnthropicRequest(model, systemPrompt, prompt)
	} else {
		// Anthropic format (Claude)
		reqBody, err = s.buildAnthropicRequest(model, systemPrompt, prompt)
	}

	if err != nil {
		return nil, err
	}

	log.Printf("[AI Analysis] Request built, body size: %d bytes", len(reqBody))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers based on API type
	if isOpenAICompatible {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
	} else if isMiniMaxAnthropic {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		if groupID != "" {
			req.Header.Set("GroupId", groupID)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-dangerous-direct-browser-access", "true")
	}

	client := &http.Client{Timeout: 60 * time.Second}
	log.Printf("[AI Analysis] Sending request to %s...", apiURL)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[AI Analysis] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("[AI Analysis] Response status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[AI Analysis] Response body size: %d bytes", len(body))
	if len(body) > 0 {
		// Percent-encode non-ASCII to prevent terminal encoding issues (mojibake)
		truncated := body[:int(math.Min(float64(500), float64(len(body))))]
		encoded := &strings.Builder{}
		for _, b := range truncated {
			if b >= 0x20 && b <= 0x7E {
				encoded.WriteByte(b)
			} else {
				encoded.WriteString(fmt.Sprintf("%%%02X", b))
			}
		}
		log.Printf("[AI Analysis] Response body (first 500 chars): %s", encoded.String())
	}

	// Parse response based on API type
	var report AIAnalysisReport
	report.Code = input.Code
	report.Name = input.Name

	var responseText string

	if isOpenAICompatible {
		responseText, err = s.parseOpenAIResponse(body)
	} else {
		responseText, err = s.parseMiniMaxAnthropicResponse(body)
	}

	if err != nil || responseText == "" {
		// Return heuristic with raw body preserved so user can see what went wrong
		rawBody := string(body)
		if len(rawBody) > 2000 {
			rawBody = rawBody[:2000] + "..."
		}
		heuristicReport := s.generateHeuristicAnalysis(input, techData)
		heuristicReport.RawAnalysis = fmt.Sprintf("[Parse error: %v]\n[Raw response: %s]", err, rawBody)
		heuristicReport.AnalysisMethod = "ai"
		log.Printf("[AI Analysis] Failed to parse AI response: %v, returning heuristic with raw body preserved", err)
		return heuristicReport, nil
	}

	// Store raw analysis text
	report.RawAnalysis = responseText

	// Try to parse as JSON first, if fails use heuristic with raw text
	if err := json.Unmarshal([]byte(responseText), &report); err != nil {
		log.Printf("[AI Analysis] AI response is not JSON (free text format), using heuristic analysis: %v", err)
		heuristicReport := s.generateHeuristicAnalysis(input, techData)
		heuristicReport.RawAnalysis = responseText
		heuristicReport.AnalysisMethod = "ai"
		log.Printf("[AI Analysis] Successfully generated report with raw AI analysis for %s", input.Code)
		return heuristicReport, nil
	}

	log.Printf("[AI Analysis] Successfully generated report for %s, composite score: %.1f", input.Code, report.Scores.CompositeScore)

	return &report, nil
}

// buildOpenAIRequest builds request body for OpenAI-compatible APIs (MiniMax, DeepSeek)
func (s *AIAnalysisService) buildOpenAIRequest(model, systemPrompt, userPrompt string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"max_tokens": 2000,
		"temperature": 0.7,
	}
	return json.Marshal(reqBody)
}

// buildAnthropicRequest builds request body for Anthropic Claude API
func (s *AIAnalysisService) buildAnthropicRequest(model, systemPrompt, userPrompt string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"model":              model,
		"max_tokens":         2000,
		"anthropic_version":  "2023-06-01",
		"system":             systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
	}
	return json.Marshal(reqBody)
}

// buildMiniMaxAnthropicRequest builds request body for MiniMax Anthropic-compatible API
func (s *AIAnalysisService) buildMiniMaxAnthropicRequest(model, systemPrompt, userPrompt string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"model":     model,
		"max_tokens": 2000,
		"stream":    false,
		"system":    systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
	}
	return json.Marshal(reqBody)
}

// parseOpenAIResponse extracts text from OpenAI-compatible API response
func (s *AIAnalysisService) parseOpenAIResponse(body []byte) (string, error) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	if resp.Error.Message != "" {
		return "", fmt.Errorf("OpenAI API error: %s", resp.Error.Message)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return resp.Choices[0].Message.Content, nil
}

// parseAnthropicResponse extracts text from Anthropic Claude API response
func (s *AIAnalysisService) parseAnthropicResponse(body []byte) (string, error) {
	var resp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	if resp.Error.Message != "" {
		return "", fmt.Errorf("Anthropic API error: %s", resp.Error.Message)
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}
	return resp.Content[0].Text, nil
}

// parseMiniMaxAnthropicResponse extracts text from MiniMax Anthropic-compatible API response
func (s *AIAnalysisService) parseMiniMaxAnthropicResponse(body []byte) (string, error) {
	var resp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	if resp.Error.Message != "" {
		return "", fmt.Errorf("MiniMax API error: %s", resp.Error.Message)
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}
	// Find the text content block (skip thinking blocks)
	for _, block := range resp.Content {
		if block.Type == "text" && block.Text != "" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("no text content in response")
}

// buildAnalysisPrompt constructs the prompt for AI analysis
func (s *AIAnalysisService) buildAnalysisPrompt(input *AIAnalysisInput, techData map[string]interface{}) string {
	quote := input.Quote
	quoteStr := ""
	if quote != nil {
		quoteStr = fmt.Sprintf(`
Stock Quote:
- Current Price: %.2f
- Open: %.2f
- High: %.2f
- Low: %.2f
- Prev Close: %.2f
- Volume: %d
- Amount: %.2f
`, quote.Current, quote.Open, quote.High, quote.Low, quote.PrevClose, quote.Volume, quote.Amount)
	}

	techStr := fmt.Sprintf(`
Technical Analysis:
- MA Status: %v
- RSI Status: %v
- MACD Status: %v
- KDJ Status: %v
- BOLL Status: %v
- Volume: %v
- Patterns: %v
- Recommendation: %v (confidence: %.0f%%)
`, techData["maStatus"], techData["rsiStatus"], techData["macdStatus"],
		techData["kdjStatus"], techData["bollStatus"], techData["volumeAnalysis"],
		techData["patterns"], techData["recommendation"], techData["confidence"])

	return fmt.Sprintf(`分析股票 %s (%s) 的投资价值。

%s

%s

请生成完整的结构化分析报告，突出3个亮点和2个风险点。用中文输出。
`, input.Code, input.Name, quoteStr, techStr)
}

// generateHeuristicAnalysis generates analysis without AI (fallback)
func (s *AIAnalysisService) generateHeuristicAnalysis(input *AIAnalysisInput, techData map[string]interface{}) *AIAnalysisReport {
	report := &AIAnalysisReport{
		Code:        input.Code,
		Name:        input.Name,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Calculate technical score
	techScore := s.calculateTechnicalScore(techData)
	report.Scores.Technical = techScore

	// Estimate other scores based on available data
	report.Scores.Fundamental = DimensionScore{
		Score:   60 + float64(s.randInt(0, 20)),
		Trend:   "stable",
		Summary: "Fundamental analysis requires detailed financial statements",
		Factors: []string{"Based on sector average"},
	}

	report.Scores.MoneyFlow = s.calculateMoneyFlowScore(techData)
	report.Scores.NewsSentiment = DimensionScore{
		Score:   50 + float64(s.randInt(0, 30)),
		Trend:   "stable",
		Summary: "News sentiment analysis unavailable",
		Factors: []string{"Limited news data"},
	}

	// Calculate composite score (30% tech, 30% fundamental, 25% money, 15% news)
	report.Scores.CompositeScore =
		report.Scores.Technical.Score*0.30 +
			report.Scores.Fundamental.Score*0.30 +
			report.Scores.MoneyFlow.Score*0.25 +
			report.Scores.NewsSentiment.Score*0.15

	// Calculate anxiety index (inverse of confidence)
	report.Scores.AnxietyIndex = 100 - report.Scores.CompositeScore

	// Attention level
	if report.Scores.CompositeScore >= 70 {
		report.Scores.AttentionLevel = "high"
	} else if report.Scores.CompositeScore >= 50 {
		report.Scores.AttentionLevel = "medium"
	} else {
		report.Scores.AttentionLevel = "low"
	}

	// Generate key findings
	report.KeyFindings = s.generateKeyFindings(techData, report.Scores)

	return report
}

func (s *AIAnalysisService) calculateTechnicalScore(techData map[string]interface{}) DimensionScore {
	score := 50.0
	var factors []string

	if status, ok := techData["maStatus"].(string); ok {
		switch status {
		case "strong_bullish":
			score += 20
			factors = append(factors, "Price above all moving averages")
		case "bullish":
			score += 10
			factors = append(factors, "Price above short-term MA")
		case "strong_bearish":
			score -= 20
			factors = append(factors, "Price below all moving averages")
		case "bearish":
			score -= 10
			factors = append(factors, "Price below short-term MA")
		}
	}

	if status, ok := techData["rsiStatus"].(string); ok {
		switch status {
		case "oversold":
			score += 10
			factors = append(factors, "RSI indicates oversold condition")
		case "overbought":
			score -= 10
			factors = append(factors, "RSI indicates overbought condition")
		}
	}

	if status, ok := techData["macdStatus"].(string); ok {
		switch status {
		case "strong_bullish":
			score += 15
			factors = append(factors, "MACD strong bullish signal")
		case "bullish":
			score += 5
			factors = append(factors, "MACD bullish signal")
		case "strong_bearish":
			score -= 15
			factors = append(factors, "MACD strong bearish signal")
		case "bearish":
			score -= 5
			factors = append(factors, "MACD bearish signal")
		}
	}

	if patterns, ok := techData["patterns"].([]string); ok && len(patterns) > 0 {
		bullishCount := 0
		bearishCount := 0
		for _, p := range patterns {
			if containsString(p, "bullish") {
				bullishCount++
			} else if containsString(p, "bearish") {
				bearishCount++
			}
		}
		score += float64(bullishCount-bearishCount) * 3
	}

	if volume, ok := techData["volumeAnalysis"].(string); ok {
		if volume == "very_high" || volume == "high" {
			score += 5
			factors = append(factors, "Above average volume")
		}
	}

	trend := "stable"
	if score > 65 {
		trend = "improving"
	} else if score < 45 {
		trend = "declining"
	}

	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return DimensionScore{
		Score:   score,
		Trend:   trend,
		Summary: fmt.Sprintf("Technical score: %.0f/100", score),
		Factors: factors,
	}
}

func (s *AIAnalysisService) calculateMoneyFlowScore(techData map[string]interface{}) DimensionScore {
	score := 50.0
	var factors []string

	// Estimate from OBV trend if available
	if bollStatus, ok := techData["bollStatus"].(string); ok {
		switch bollStatus {
		case "above_upper_band":
			score += 10
			factors = append(factors, "Price near upper Bollinger Band - strong momentum")
		case "below_lower_band":
			score -= 10
			factors = append(factors, "Price near lower Bollinger Band - weak momentum")
		case "above_midline":
			score += 5
			factors = append(factors, "Price above middle Bollinger Band")
		}
	}

	if kdjStatus, ok := techData["kdjStatus"].(string); ok {
		switch kdjStatus {
		case "overbought":
			score -= 5
			factors = append(factors, "KDJ overbought - possible pullback")
		case "oversold":
			score += 5
			factors = append(factors, "KDJ oversold - possible bounce")
		}
	}

	trend := "stable"
	if score > 60 {
		trend = "improving"
	} else if score < 40 {
		trend = "declining"
	}

	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return DimensionScore{
		Score:   score,
		Trend:   trend,
		Summary: fmt.Sprintf("Money flow score: %.0f/100", score),
		Factors: factors,
	}
}

func (s *AIAnalysisService) generateKeyFindings(techData map[string]interface{}, scores Scores) KeyFindings {
	highlights := []Highlight{}
	risks := []Risk{}

	// Generate highlights based on positive signals
	if status, ok := techData["maStatus"].(string); ok {
		if status == "strong_bullish" || status == "bullish" {
			highlights = append(highlights, Highlight{
				Title:   "Technical Uptrend",
				Context: "The stock is trading above key moving averages, indicating a positive technical trend",
			})
		}
	}

	if rsiStatus, ok := techData["rsiStatus"].(string); ok {
		if rsiStatus == "oversold" {
			highlights = append(highlights, Highlight{
				Title:   "RSI Oversold",
				Context: "RSI below 30 suggests potential bounce opportunity",
			})
		}
	}

	if len(highlights) < 3 {
		highlights = append(highlights, Highlight{
			Title:   "Sector Strength",
			Context: "Stock shows relative strength within its sector",
		})
	}

	// Generate risks based on negative signals
	if status, ok := techData["maStatus"].(string); ok {
		if status == "strong_bearish" || status == "bearish" {
			risks = append(risks, Risk{
				Title:   "Technical Downtrend",
				Context: "The stock is trading below key moving averages, indicating negative technical trend",
			})
		}
	}

	if rsiStatus, ok := techData["rsiStatus"].(string); ok {
		if rsiStatus == "overbought" {
			risks = append(risks, Risk{
				Title:   "RSI Overbought",
				Context: "RSI above 70 suggests potential pullback risk",
			})
		}
	}

	if len(risks) < 2 {
		risks = append(risks, Risk{
			Title:   "Market Volatility",
			Context: "General market uncertainty may affect stock performance",
		})
	}

	return KeyFindings{
		Highlights: highlights,
		Risks:      risks,
	}
}

// findSimilarPatterns finds similar historical K-line patterns
func (s *AIAnalysisService) findSimilarPatterns(historicalData []KLine) []SimilarPattern {
	var patterns []SimilarPattern

	if len(historicalData) < 60 {
		return patterns
	}

	// Normalize recent 30-day pattern for comparison
	recentCount := 30
	if len(historicalData) < recentCount {
		recentCount = len(historicalData)
	}

	recentPrices := historicalData[len(historicalData)-recentCount:]
	recentNormalized := s.normalizePattern(recentPrices)

	// Search through historical data for similar patterns
	// Use sliding window to find similar segments
	for i := 0; i < len(historicalData)-recentCount*2; i += 5 {
		window := historicalData[i : i+recentCount]
		normalized := s.normalizePattern(window)

		similarity := s.calculatePatternSimilarity(recentNormalized, normalized)

		if similarity > 70 { // Only keep highly similar patterns
			// Calculate price change during this pattern
			priceChange := (window[len(window)-1].Close - window[0].Close) / window[0].Close * 100

			// Calculate future returns (next 5 and 20 days)
			var next5DayWinRate, next20DayWinRate float64
			if i+recentCount+5 < len(historicalData) {
				future5 := historicalData[i+recentCount : i+recentCount+5]
				if len(future5) == 5 {
					change5 := (future5[4].Close - future5[0].Close) / future5[0].Close * 100
					next5DayWinRate = s.calculateWinRate(historicalData, i, recentCount, 5, change5)
				}
			}
			if i+recentCount+20 < len(historicalData) {
				future20 := historicalData[i+recentCount : i+recentCount+20]
				if len(future20) == 20 {
					change20 := (future20[19].Close - future20[0].Close) / future20[0].Close * 100
					next20DayWinRate = s.calculateWinRate(historicalData, i, recentCount, 20, change20)
				}
			}

			patterns = append(patterns, SimilarPattern{
				PatternID:        fmt.Sprintf("PATT-%d", i),
				StartDate:        window[0].Date,
				EndDate:          window[len(window)-1].Date,
				Similarity:       similarity,
				PriceChange:      priceChange,
				Next5DayWinRate:  next5DayWinRate,
				Next20DayWinRate: next20DayWinRate,
			})
		}
	}

	// Sort by similarity and take top 3
	s.sortBySimilarity(patterns)

	if len(patterns) > 3 {
		patterns = patterns[:3]
	}

	return patterns
}

// normalizePattern normalizes a price pattern for comparison
func (s *AIAnalysisService) normalizePattern(prices []KLine) []float64 {
	if len(prices) == 0 {
		return nil
	}

	normalized := make([]float64, len(prices))
	firstClose := prices[0].Close
	if firstClose == 0 {
		firstClose = 1
	}

	for i, p := range prices {
		normalized[i] = (p.Close - firstClose) / firstClose * 100
	}

	return normalized
}

// calculatePatternSimilarity calculates similarity between two normalized patterns
func (s *AIAnalysisService) calculatePatternSimilarity(p1, p2 []float64) float64 {
	if len(p1) != len(p2) || len(p1) == 0 {
		return 0
	}

	var sumSquaredDiff float64
	for i := range p1 {
		diff := p1[i] - p2[i]
		sumSquaredDiff += diff * diff
	}

	rmse := math.Sqrt(sumSquaredDiff / float64(len(p1)))

	// Convert RMSE to similarity percentage (0-100)
	// Lower RMSE = higher similarity
	similarity := 100 - (rmse * 2) // Scale factor
	if similarity < 0 {
		similarity = 0
	}
	if similarity > 100 {
		similarity = 100
	}

	return similarity
}

// calculateWinRate calculates the historical win rate after a similar pattern
func (s *AIAnalysisService) calculateWinRate(allData []KLine, patternIdx, patternLen, futureLen int, currentChange float64) float64 {
	count := 0
	total := 0

	for i := 0; i < len(allData)-patternLen-futureLen; i += 10 { // Sample every 10 days
		if i == patternIdx {
			continue // Skip current pattern
		}

		future := allData[i+patternLen : i+patternLen+futureLen]
		if len(future) < futureLen {
			continue
		}

		change := (future[len(future)-1].Close - future[0].Close) / future[0].Close * 100

		// Win = same direction as current pattern
		if (currentChange > 0 && change > 0) || (currentChange < 0 && change < 0) {
			count++
		}
		total++
	}

	if total == 0 {
		return 50.0 // Default to 50% if no data
	}

	return float64(count) / float64(total) * 100
}

func (s *AIAnalysisService) sortBySimilarity(patterns []SimilarPattern) {
	for i := 0; i < len(patterns)-1; i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[j].Similarity > patterns[i].Similarity {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}
}

// getInstitutionalSentiment returns mock institutional sentiment data
func (s *AIAnalysisService) getInstitutionalSentiment(code string) InstitutionalSentiment {
	// In production, this would fetch from a research report API
	// For now, return simulated data
	return InstitutionalSentiment{
		Period:          "30d",
		TotalReports:    15 + s.randInt(0, 20),
		BuyRatio:        40 + float64(s.randInt(0, 30)),
		HoldRatio:       30 + float64(s.randInt(0, 20)),
		SellRatio:       10 + float64(s.randInt(0, 15)),
		RatioChange:     s.randomChoice([]string{"improving", "stable", "deteriorating"}),
		ConsensusTarget: 0,
		ConsensusRating: "hold",
	}
}

// cleanExpiredCache removes expired cache entries
func (s *AIAnalysisService) cleanExpiredCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	now := time.Now()
	for key, entry := range s.cache {
		if now.After(entry.expiry) {
			delete(s.cache, key)
		}
	}
}

// Helper functions
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (s *AIAnalysisService) randInt(min, max int) int {
	return min + int(time.Now().UnixNano()%int64(max-min+1))
}

func (s *AIAnalysisService) randomChoice(choices []string) string {
	return choices[s.randInt(0, len(choices)-1)]
}
