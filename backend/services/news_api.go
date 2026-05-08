package services

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// EastMoneyNewsAPI fetches news/announcements from East Money API
func EastMoneyNewsAPI(stockCode string) ([]models.NewsItem, error) {
	// Convert to East Money format (sh600519 -> 600519)
	code := strings.TrimPrefix(strings.ToLower(stockCode), "sh")
	code = strings.TrimPrefix(code, "sz")

	url := fmt.Sprintf(
		"https://np-anotice-stock.eastmoney.com/api/security/ann?sr=-1&page_size=20&page_index=1&ann_type=SHA,CYB,SZA&client_source=web&stock_list=%s",
		code,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://data.eastmoney.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseEastMoneyNews(body, stockCode)
}

// parseEastMoneyNews parses East Money API response
func parseEastMoneyNews(body []byte, stockCode string) ([]models.NewsItem, error) {
	var resp struct {
		Data struct {
			List []struct {
				ArtCode     string `json:"art_code"`
				Title       string `json:"title"`
				NoticeDate  string `json:"notice_date"`
				DisplayTime string `json:"display_time"`
				Columns     []struct {
					ColumnName string `json:"column_name"`
				} `json:"columns"`
				Codes []struct {
					ShortName string `json:"short_name"`
					StockCode string `json:"stock_code"`
				} `json:"codes"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var newsItems []models.NewsItem
	for _, item := range resp.Data.List {
		noticeDate := item.NoticeDate
		if noticeDate == "" {
			noticeDate = item.DisplayTime
		}

		publishTime, err := time.Parse("2006-01-02 15:04:05", noticeDate)
		if err != nil {
			publishTime, _ = time.Parse("2006-01-02", noticeDate)
			if publishTime.IsZero() {
				publishTime = time.Now()
			}
		}

		// Get column name as source
		source := "东方财富"
		if len(item.Columns) > 0 {
			source = item.Columns[0].ColumnName
		}

		// Build detail URL
		sourceURL := fmt.Sprintf("https://data.eastmoney.com/announcement/detail/AN%s.html", item.ArtCode)

		// Analyze sentiment based on title
		sentiment, score := AnalyzeNewsSentiment(item.Title)
		newsItems = append(newsItems, models.NewsItem{
			StockCode:     stockCode,
			Title:         item.Title,
			Content:       "",
			Source:        source,
			SourceURL:     sourceURL,
			PublishTime:   publishTime,
			Sentiment:     sentiment,
			SentimentScore: score,
		})
	}

	return newsItems, nil
}

// AnalyzeNewsSentiment performs keyword-based sentiment analysis on news text
func AnalyzeNewsSentiment(text string) (string, float64) {
	positiveKeywords := []string{"涨", "突破", "创新高", "增持", "买入", "利好", "增长", "盈利", "超预期", "上升", "上涨", "看好", "推荐", "强劲", "上调", "加仓", "业绩", "回购", "分红", "超配", "买入评级", "大幅增长", "首次覆盖", "目标价"}
	negativeKeywords := []string{"跌", "破", "创新低", "减持", "卖出", "利空", "亏损", "预警", "风险", "下降", "下跌", "下调", "减仓", "业绩下滑", "疲软", "警惕", "下架", "整改", "调查", "警示函", "风险提示", "减持", "大宗交易折价"}

	positiveCount := 0
	negativeCount := 0

	for _, kw := range positiveKeywords {
		positiveCount += strings.Count(text, kw)
	}
	for _, kw := range negativeKeywords {
		negativeCount += strings.Count(text, kw)
	}

	total := positiveCount + negativeCount
	if total == 0 {
		return "neutral", 0.0
	}

	score := float64(positiveCount-negativeCount) / float64(total)

	if score > 0.2 {
		return "positive", score
	} else if score < -0.2 {
		return "negative", score
	}
	return "neutral", score
}

// GetNewsForStock retrieves and stores news for a specific stock
func GetNewsForStock(stockCode string) (*models.NewsSentimentResult, error) {
	newsItems, err := EastMoneyNewsAPI(stockCode)
	if err != nil {
		log.Printf("East Money news fetch failed for %s: %v", stockCode, err)
		return &models.NewsSentimentResult{
			StockCode: stockCode,
			NewsCount: 0,
		}, nil
	}

	if len(newsItems) == 0 {
		return &models.NewsSentimentResult{
			StockCode: stockCode,
			NewsCount: 0,
		}, nil
	}

	// Save to database and update sentiment
	for i := range newsItems {
		saveNewsItem(&newsItems[i])
	}

	// Calculate aggregated sentiment
	return calculateNewsSentiment(stockCode, newsItems), nil
}

// saveNewsItem stores a news item in the database (upsert)
func saveNewsItem(item *models.NewsItem) error {
	if item.SourceURL == "" {
		return nil
	}
	return config.DB.Where("source_url = ? AND stock_code = ?",
		item.SourceURL, item.StockCode).
		Assign(models.NewsItem{
			StockCode:     item.StockCode,
			Title:         item.Title,
			Content:       item.Content,
			Source:        item.Source,
			PublishTime:   item.PublishTime,
			Sentiment:     item.Sentiment,
			SentimentScore: item.SentimentScore,
		}).FirstOrCreate(item).Error
}

// calculateNewsSentiment aggregates sentiment from news items
func calculateNewsSentiment(stockCode string, items []models.NewsItem) *models.NewsSentimentResult {
	result := &models.NewsSentimentResult{
		StockCode: stockCode,
	}

	var totalScore float64
	latestTime := time.Time{}

	for _, item := range items {
		result.NewsCount++
		totalScore += item.SentimentScore

		switch item.Sentiment {
		case "positive":
			result.PositiveCount++
		case "negative":
			result.NegativeCount++
		default:
			result.NeutralCount++
		}

		if item.PublishTime.After(latestTime) {
			latestTime = item.PublishTime
		}
	}

	if result.NewsCount > 0 {
		result.OverallScore = (totalScore / float64(result.NewsCount)) * 100
		result.LatestNewsTime = latestTime.Format("2006-01-02 15:04:05")
	}

	return result
}

// GetRecentNewsFromDB retrieves cached news from database
func GetRecentNewsFromDB(stockCode string, limit int) ([]models.NewsItem, error) {
	var items []models.NewsItem
	err := config.DB.Where("stock_code = ?", stockCode).
		Order("publish_time DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}
