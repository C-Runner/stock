package services

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// KLine represents a single candlestick data point
type KLine struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	Amount   float64
}

// TechnicalAnalysis holds all technical indicators
type TechnicalAnalysis struct {
	Code         string          `json:"code"`
	Name         string          `json:"name"`
	MA           []MAData        `json:"ma"`           // Moving averages
	EMA          []MAData        `json:"ema"`          // Exponential moving averages
	RSI          []RSIData       `json:"rsi"`         // Relative Strength Index
	MACD         MACDData        `json:"macd"`         // MACD indicator
	KDJ          KDJData         `json:"kdj"`         // KDJ indicator
	BOLL         BOLLData        `json:"boll"`         // Bollinger Bands
	RecentPrices []PricePoint    `json:"recentPrices"` // Recent price data for chart
}

// MAData represents moving average data
type MAData struct {
	Period int     `json:"period"`
	Value  float64 `json:"value"`
}

// RSIData represents RSI data
type RSIData struct {
	Period int     `json:"period"`
	Value  float64 `json:"value"`
}

// MACDData represents MACD indicator
type MACDData struct {
	DIF   float64 `json:"dif"`   // MACD line
	DEA   float64 `json:"dea"`   // Signal line
	MACD  float64 `json:"macd"` // Histogram
}

// KDJData represents KDJ indicator
type KDJData struct {
	K float64 `json:"k"` // %K line
	D float64 `json:"d"` // %D line
	J float64 `json:"j"` // %J line
}

// BOLLData represents Bollinger Bands
type BOLLData struct {
	Upper float64 `json:"upper"` // Upper band
	Mid   float64 `json:"mid"`   // Middle band (MA20)
	Lower float64 `json:"lower"` // Lower band
}

// PricePoint represents a price data point for charts
type PricePoint struct {
	Date  string  `json:"date"`
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Volume int64  `json:"volume"`
}

// GetTechnicalAnalysis fetches technical analysis for a stock
func GetTechnicalAnalysis(code string) (*TechnicalAnalysis, error) {
	// Get stock basic info
	sinaCode := convertToSinaCode(code)
	quote, err := SinaFinanceAPI(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Fetch historical k-line data
	klines, err := fetchKLineData(sinaCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}

	if len(klines) < 50 {
		return nil, fmt.Errorf("insufficient historical data")
	}

	// Calculate indicators
	analysis := &TechnicalAnalysis{
		Code: code,
		Name: quote.Name,
	}

	// Get recent prices for chart (last 60 days)
	recentCount := len(klines)
	if recentCount > 60 {
		recentCount = 60
	}
	for i := len(klines) - recentCount; i < len(klines); i++ {
		analysis.RecentPrices = append(analysis.RecentPrices, PricePoint{
			Date:  klines[i].Date,
			Open:  klines[i].Open,
			Close: klines[i].Close,
			High:  klines[i].High,
			Low:   klines[i].Low,
			Volume: klines[i].Volume,
		})
	}

	// Calculate Moving Averages
	closes := make([]float64, len(klines))
	for i, k := range klines {
		closes[i] = k.Close
	}

	analysis.MA = calculateMA(closes, []int{5, 10, 20, 60})
	analysis.EMA = calculateEMA(closes, []int{12, 26})

	// Calculate RSI
	analysis.RSI = calculateRSI(closes, []int{6, 12, 24})

	// Calculate MACD
	analysis.MACD = calculateMACD(closes)

	// Calculate KDJ
	analysis.KDJ = calculateKDJ(klines)

	// Calculate Bollinger Bands
	analysis.BOLL = calculateBOLL(closes)

	return analysis, nil
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
	// Response format: [{"day":"2024-01-01","open":100.00,"high":101.00,"low":99.00,"close":100.50,"volume":1234567},...]
	body = strings.Trim(body, "[]")

	// Simple parsing - split by }{ and parse each object
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

// calculateMA calculates moving averages
func calculateMA(closes []float64, periods []int) []MAData {
	var result []MAData
	for _, period := range periods {
		if len(closes) < period {
			continue
		}
		sum := 0.0
		for i := len(closes) - period; i < len(closes); i++ {
			sum += closes[i]
		}
		result = append(result, MAData{
			Period: period,
			Value:  sum / float64(period),
		})
	}
	return result
}

// calculateEMA calculates exponential moving averages
func calculateEMA(closes []float64, periods []int) []MAData {
	var result []MAData
	for _, period := range periods {
		if len(closes) < period {
			continue
		}
		// Simple EMA calculation
		multiplier := 2.0 / float64(period+1)
		ema := closes[0]
		for i := 1; i < len(closes); i++ {
			ema = (closes[i]-ema)*multiplier + ema
		}
		result = append(result, MAData{
			Period: period,
			Value:  ema,
		})
	}
	return result
}

// calculateRSI calculates Relative Strength Index
func calculateRSI(closes []float64, periods []int) []RSIData {
	var result []RSIData
	for _, period := range periods {
		if len(closes) < period+1 {
			continue
		}

		var gains, losses float64
		for i := len(closes) - period; i < len(closes); i++ {
			change := closes[i] - closes[i-1]
			if change > 0 {
				gains += change
			} else {
				losses -= change
			}
		}

		avgGain := gains / float64(period)
		avgLoss := losses / float64(period)

		if avgLoss == 0 {
			result = append(result, RSIData{Period: period, Value: 100})
			continue
		}

		rs := avgGain / avgLoss
		rsi := 100 - (100 / (1 + rs))
		result = append(result, RSIData{Period: period, Value: rsi})
	}
	return result
}

// calculateMACD calculates MACD indicator
func calculateMACD(closes []float64) MACDData {
	if len(closes) < 34 {
		return MACDData{}
	}

	// Calculate 12-day and 26-day EMA
	ema12 := calculateEMAValue(closes, 12)
	ema26 := calculateEMAValue(closes, 26)

	dif := ema12 - ema26

	// Calculate 9-day DEA (signal line) using last 9 MACD values
	var macdValues []float64
	for i := len(closes) - 9; i < len(closes); i++ {
		if i >= 26 {
			e12 := calculateEMAValue(closes[:i+1], 12)
			e26 := calculateEMAValue(closes[:i+1], 26)
			macdValues = append(macdValues, e12-e26)
		}
	}

	dea := 0.0
	if len(macdValues) > 0 {
		for _, v := range macdValues {
			dea += v
		}
		dea /= float64(len(macdValues))
	}

	macd := dif - dea

	return MACDData{
		DIF:  dif,
		DEA:  dea,
		MACD: macd * 2, // Some sources multiply by 2
	}
}

// calculateEMAValue calculates a single EMA value
func calculateEMAValue(closes []float64, period int) float64 {
	if len(closes) < period {
		return 0
	}
	multiplier := 2.0 / float64(period+1)
	ema := closes[0]
	for i := 1; i < len(closes); i++ {
		ema = (closes[i]-ema)*multiplier + ema
	}
	return ema
}

// calculateKDJ calculates KDJ indicator
func calculateKDJ(klines []KLine) KDJData {
	if len(klines) < 9 {
		return KDJData{}
	}

	// Use last 9 periods for KDJ calculation
	n := 9
	if len(klines) < n {
		n = len(klines)
	}

	var rsvs []float64
	for i := len(klines) - n; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		for j := len(klines) - n; j <= i; j++ {
			if klines[j].High > high {
				high = klines[j].High
			}
			if klines[j].Low < low {
				low = klines[j].Low
			}
		}
		if high == low {
			rsvs = append(rsvs, 50)
		} else {
			rsv := (klines[i].Close - low) / (high - low) * 100
			rsvs = append(rsvs, rsv)
		}
	}

	// Calculate K, D, J
	k := 50.0
	d := 50.0
	for _, rsv := range rsvs {
		k = (2.0*k + rsv) / 3.0
		d = (2.0*d + k) / 3.0
	}
	j := 3*k - 2*d

	return KDJData{
		K: k,
		D: d,
		J: j,
	}
}

// calculateBOLL calculates Bollinger Bands
func calculateBOLL(closes []float64) BOLLData {
	if len(closes) < 20 {
		return BOLLData{}
	}

	// Calculate 20-day MA
	sum := 0.0
	for i := len(closes) - 20; i < len(closes); i++ {
		sum += closes[i]
	}
	ma20 := sum / 20.0

	// Calculate standard deviation
	var variance float64
	for i := len(closes) - 20; i < len(closes); i++ {
		diff := closes[i] - ma20
		variance += diff * diff
	}
	stdDev := variance / 20.0
	stdDev = sqrt(stdDev)

	upper := ma20 + 2*stdDev
	lower := ma20 - 2*stdDev

	return BOLLData{
		Upper: upper,
		Mid:   ma20,
		Lower: lower,
	}
}

// sqrt returns square root
func sqrt(n float64) float64 {
	if n <= 0 {
		return 0
	}
	x := n
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + n/x) / 2
	}
	return x
}
