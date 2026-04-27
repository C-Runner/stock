package services

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetTechnicalAnalysis fetches technical analysis for a stock
func GetTechnicalAnalysis(code string) (*TechnicalAnalysis, error) {
	// Get stock basic info
	sinaCode := convertToSinaCode(code)
	quote, err := SinaFinanceAPI(code)
	stockName := code
	if err == nil {
		stockName = quote.Name
	}

	// Fetch historical k-line data
	klines, err := fetchKLineData(sinaCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}

	// Filter to last 180 days (approximately)
	klines = filterKLinesByDays(klines, 180)

	if len(klines) < 50 {
		return nil, fmt.Errorf("insufficient historical data")
	}

	// Calculate indicators
	analysis := &TechnicalAnalysis{
		Code: code,
		Name: stockName,
	}

	// Get recent prices for chart (last 60 days)
	recentCount := len(klines)
	if recentCount > 60 {
		recentCount = 60
	}
	for i := len(klines) - recentCount; i < len(klines); i++ {
		analysis.RecentPrices = append(analysis.RecentPrices, PricePoint{
			Date:   klines[i].Date,
			Open:   klines[i].Open,
			Close:  klines[i].Close,
			High:   klines[i].High,
			Low:    klines[i].Low,
			Volume: klines[i].Volume,
		})
	}

	// Get closes for indicator calculations
	closes := make([]float64, len(klines))
	for i, k := range klines {
		closes[i] = k.Close
	}

	// Calculate Moving Averages (time series)
	analysis.MA = calculateMA(closes, []int{5, 10, 20, 60})
	analysis.EMA = calculateEMA(closes, []int{12, 26})
	analysis.WMA = calculateWMA(closes, []int{5, 10, 20})

	// Calculate RSI (time series)
	analysis.RSI = calculateRSI(closes, []int{6, 12, 24})

	// Calculate MACD (time series)
	analysis.MACD = calculateMACD(closes)

	// Calculate KDJ (time series)
	analysis.KDJ = calculateKDJ(klines)

	// Calculate Bollinger Bands (time series)
	analysis.BOLL = calculateBOLL(closes)

	// Calculate Williams %R (time series)
	analysis.WR = calculateWR(klines, []int{14, 20})

	// Calculate DMI
	analysis.DMI = calculateDMI(klines, 14)

	// Calculate OBV
	analysis.OBV = calculateOBV(klines)

	// Calculate VWAP
	analysis.VWAP = calculateVWAP(klines)

	// Pattern recognition
	analysis.Patterns = analyzePatterns(klines, closes)

	// Support/Resistance levels
	analysis.Levels = detectSupportResistance(klines)

	// Generate recommendation
	analysis.Recommendation = generateRecommendation(analysis)

	return analysis, nil
}

// filterKLinesByDays filters klines to only include data from the last n days
func filterKLinesByDays(klines []KLine, days int) []KLine {
	if len(klines) == 0 || days <= 0 {
		return klines
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	var filtered []KLine

	for _, k := range klines {
		t, err := time.Parse("2006-01-02", k.Date)
		if err != nil {
			filtered = append(filtered, k)
			continue
		}
		if t.After(cutoff) || t.Equal(cutoff) {
			filtered = append(filtered, k)
		}
	}

	if len(filtered) < 30 {
		return klines
	}

	return filtered
}

// fetchKLineData fetches historical k-line data from Sina
func fetchKLineData(sinaCode string) ([]KLine, error) {
	url := fmt.Sprintf("http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=%s&scale=240&ma=5&datalen=1024", sinaCode)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "http://finance.sina.com.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseKLineData(string(body))
}

// parseKLineData parses the k-line data
func parseKLineData(body string) ([]KLine, error) {
	body = strings.Trim(body, "[]")

	entries := strings.Split(body, "},{")
	var klines []KLine

	for _, entry := range entries {
		entry = strings.Trim(entry, "{}")
		entry = strings.ReplaceAll(entry, `"`, ``)

		fields := strings.Split(entry, ",")
		kline := KLine{}

		for _, field := range fields {
			parts := strings.Split(field, ":")
			if len(parts) != 2 {
				continue
			}
			key := strings.Trim(parts[0], " ")
			value := strings.Trim(parts[1], " ")

			switch key {
			case "day":
				kline.Date = value
			case "open":
				kline.Open, _ = strconv.ParseFloat(value, 64)
			case "high":
				kline.High, _ = strconv.ParseFloat(value, 64)
			case "low":
				kline.Low, _ = strconv.ParseFloat(value, 64)
			case "close":
				kline.Close, _ = strconv.ParseFloat(value, 64)
			case "volume":
				kline.Volume, _ = strconv.ParseInt(value, 10, 64)
			}
		}

		if kline.Date != "" {
			klines = append(klines, kline)
		}
	}

	return klines, nil
}
