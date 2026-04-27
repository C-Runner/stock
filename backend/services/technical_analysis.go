package services

import (
	"fmt"
	"io"
	"math"
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
	Code            string              `json:"code"`
	Name            string              `json:"name"`
	MA              []MAData            `json:"ma"`              // Moving averages
	EMA             []MAData            `json:"ema"`             // Exponential moving averages
	WMA             []MAData            `json:"wma"`             // Weighted moving averages
	RSI             []RSIData           `json:"rsi"`             // Relative Strength Index
	MACD            MACDData            `json:"macd"`            // MACD indicator
	KDJ             KDJData             `json:"kdj"`             // KDJ indicator
	BOLL            BOLLData            `json:"boll"`            // Bollinger Bands
	WR              []WRData            `json:"wr"`              // Williams %R
	DMI             DMIData             `json:"dmi"`             // Directional Movement Index
	OBV             OBVData             `json:"obv"`             // On-Balance Volume
	VWAP            VWAPData            `json:"vwap"`            // Volume Weighted Average Price
	Patterns        PatternAnalysis     `json:"patterns"`        // Pattern recognition results
	Levels          PriceLevels        `json:"levels"`          // Support/Resistance levels
	Recommendation  Recommendation     `json:"recommendation"` // Buy/Hold/Watch recommendation
	RecentPrices    []PricePoint       `json:"recentPrices"`    // Recent price data for chart
}

// MAData represents moving average data for each date
type MAData struct {
	Period int       `json:"period"`
	Values []float64 `json:"values"` // Array of MA values, one per date
}

// RSIData represents RSI data for each date
type RSIData struct {
	Period int       `json:"period"`
	Values []float64 `json:"values"`
}

// MACDData represents MACD indicator time series
type MACDData struct {
	DIF  []float64 `json:"dif"`  // MACD line time series
	DEA  []float64 `json:"dea"`  // Signal line time series
	MACD []float64 `json:"macd"` // Histogram time series
}

// KDJData represents KDJ indicator time series
type KDJData struct {
	K []float64 `json:"k"` // %K line time series
	D []float64 `json:"d"` // %D line time series
	J []float64 `json:"j"` // %J line time series
}

// BOLLData represents Bollinger Bands
type BOLLData struct {
	Upper []float64 `json:"upper"` // Upper band time series
	Mid   []float64 `json:"mid"`   // Middle band time series (MA20)
	Lower []float64 `json:"lower"` // Lower band time series
}

// WRData represents Williams %R
type WRData struct {
	Period int       `json:"period"`
	Values []float64 `json:"values"`
}

// DMIData represents Directional Movement Index
type DMIData struct {
	PlusDI  float64 `json:"plusDI"`  // +DI
	MinusDI float64 `json:"minusDI"` // -DI
	ADX     float64 `json:"adx"`     // ADX
}

// OBVData represents On-Balance Volume
type OBVData struct {
	Values []float64 `json:"values"` // OBV values time series
	Trend  string    `json:"trend"`  // "rising", "falling", "neutral"
}

// VWAPData represents Volume Weighted Average Price
type VWAPData struct {
	Values []float64 `json:"values"` // VWAP values time series
}

// PatternAnalysis holds all pattern recognition results
type PatternAnalysis struct {
	CandlestickPatterns []CandlestickPattern `json:"candlestickPatterns"` // K-line patterns
	TrendPatterns       []TrendPattern      `json:"trendPatterns"`      // Trend/chart patterns
}

// CandlestickPattern represents a recognized candlestick pattern
type CandlestickPattern struct {
	Type     string `json:"type"`     // Pattern name
	Date     string `json:"date"`     // Date when pattern was detected
	Strength string `json:"strength"` // "strong", "medium", "weak"
	Bullish  bool   `json:"bullish"`  // True if bullish pattern
	Meaning  string `json:"meaning"`  // Pattern meaning in English
}

// TrendPattern represents a recognized chart pattern
type TrendPattern struct {
	Type      string  `json:"type"`      // "double_top", "head_shoulders", etc.
	StartDate string  `json:"startDate"` // Pattern start date
	EndDate   string  `json:"endDate"`   // Pattern end date
	Strength  string  `json:"strength"`  // "strong", "medium", "weak"
	Bullish   bool    `json:"bullish"`   // True if bullish pattern
	Breakout  float64 `json:"breakout"` // Breakout price level
	Meaning   string  `json:"meaning"`   // Pattern meaning in English
}

// PriceLevels represents support and resistance levels
type PriceLevels struct {
	Resistance []PriceLevel `json:"resistance"` // Resistance levels
	Support    []PriceLevel `json:"support"`    // Support levels
	TrendLine  *TrendLine   `json:"trendLine"`  // Current trend line
	Channel    *Channel     `json:"channel"`    // Price channel if detected
}

// PriceLevel represents a single price level
type PriceLevel struct {
	Price    float64 `json:"price"`    // Price level
	Strength string  `json:"strength"` // "strong", "medium", "weak"
	Touches  int     `json:"touches"`  // Number of times price touched
	Type     string  `json:"type"`     // "resistance" or "support"
}

// TrendLine represents a detected trend line
type TrendLine struct {
	Slope     float64 `json:"slope"`     // Slope of the trend line
	Intercept float64 `json:"intercept"` // Y-intercept
	StartDate string  `json:"startDate"` // Trend line start date
	EndDate   string  `json:"endDate"`   // Trend line end date
	Strength  string  `json:"strength"`  // "strong", "medium", "weak"
	Direction string  `json:"direction"` // "up" or "down"
}

// Channel represents a price channel
type Channel struct {
	UpperLine TrendLine `json:"upperLine"` // Upper boundary
	LowerLine TrendLine `json:"lowerLine"` // Lower boundary
	StartDate string    `json:"startDate"` // Channel start date
	EndDate   string    `json:"endDate"`   // Channel end date
	Strength  string    `json:"strength"`  // "strong", "medium", "weak"
}

// Recommendation represents trading recommendation
type Recommendation struct {
	Action     string   `json:"action"`     // "buy", "hold", "watch"
	Confidence float64  `json:"confidence"` // 0-100 confidence score
	Reasons    []string `json:"reasons"`   // Reasons for the recommendation
	Summary    string   `json:"summary"`   // English summary of the analysis
	RiskLevel string   `json:"riskLevel"`  // "low", "medium", "high"
}

// PricePoint represents a price data point for charts
type PricePoint struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume int64   `json:"volume"`
}

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

	// Parse dates and filter
	cutoff := time.Now().AddDate(0, 0, -days)
	var filtered []KLine

	for _, k := range klines {
		// Try to parse date (format: "2006-01-02")
		t, err := time.Parse("2006-01-02", k.Date)
		if err != nil {
			// If can't parse date, include it
			filtered = append(filtered, k)
			continue
		}
		if t.After(cutoff) || t.Equal(cutoff) {
			filtered = append(filtered, k)
		}
	}

	// If filtered is too small, return original
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

// calculateMA calculates moving averages as time series
func calculateMA(closes []float64, periods []int) []MAData {
	var result []MAData
	recentCount := len(closes)
	if recentCount > 60 {
		recentCount = 60
	}

	for _, period := range periods {
		if len(closes) < period {
			continue
		}
		values := make([]float64, 0, recentCount)
		for i := period - 1; i < len(closes); i++ {
			sum := 0.0
			for j := i - period + 1; j <= i; j++ {
				sum += closes[j]
			}
			values = append(values, sum/float64(period))
		}
		result = append(result, MAData{
			Period: period,
			Values: values,
		})
	}
	return result
}

// calculateEMA calculates exponential moving averages as time series
func calculateEMA(closes []float64, periods []int) []MAData {
	var result []MAData
	recentCount := len(closes)
	if recentCount > 60 {
		recentCount = 60
	}

	for _, period := range periods {
		if len(closes) < period {
			continue
		}
		multiplier := 2.0 / float64(period+1)
		values := make([]float64, 0, recentCount)

		// Start EMA from period-1 index
		for i := period - 1; i < len(closes); i++ {
			if i == period-1 {
				// First EMA is SMA
				sum := 0.0
				for j := 0; j <= i; j++ {
					sum += closes[j]
				}
				values = append(values, sum/float64(period))
			} else {
				prevEMA := values[len(values)-1]
				ema := (closes[i]-prevEMA)*multiplier + prevEMA
				values = append(values, ema)
			}
		}
		result = append(result, MAData{
			Period: period,
			Values: values,
		})
	}
	return result
}

// calculateRSI calculates RSI as time series
func calculateRSI(closes []float64, periods []int) []RSIData {
	var result []RSIData
	recentCount := len(closes)
	if recentCount > 60 {
		recentCount = 60
	}

	for _, period := range periods {
		if len(closes) < period+1 {
			continue
		}
		values := make([]float64, 0, recentCount)

		for i := period; i < len(closes); i++ {
			var gains, losses float64
			for j := i - period + 1; j <= i; j++ {
				change := closes[j] - closes[j-1]
				if change > 0 {
					gains += change
				} else {
					losses -= change
				}
			}
			avgGain := gains / float64(period)
			avgLoss := losses / float64(period)

			var rsi float64
			if avgLoss == 0 {
				rsi = 100
			} else {
				rs := avgGain / avgLoss
				rsi = 100 - (100 / (1 + rs))
			}
			values = append(values, rsi)
		}
		result = append(result, RSIData{
			Period: period,
			Values: values,
		})
	}
	return result
}

// calculateMACD calculates MACD as time series
func calculateMACD(closes []float64) MACDData {
	if len(closes) < 34 {
		return MACDData{}
	}

	recentCount := len(closes)
	if recentCount > 60 {
		recentCount = 60
	}

	dif := make([]float64, 0, recentCount)
	dea := make([]float64, 0, recentCount)
	macd := make([]float64, 0, recentCount)

	for i := 25; i < len(closes); i++ {
		ema12 := calculateEMAValue(closes[:i+1], 12)
		ema26 := calculateEMAValue(closes[:i+1], 26)
		dif = append(dif, ema12-ema26)
	}

	// Calculate DEA (9-day EMA of DIF)
	for i := 8; i < len(dif); i++ {
		window := dif[i-8 : i+1]
		sum := 0.0
		for _, v := range window {
			sum += v
		}
		dea = append(dea, sum/9.0)
	}

	// Calculate MACD histogram
	for i := 0; i < len(dif) && i < len(dea); i++ {
		macd = append(macd, (dif[i]-dea[i])*2)
	}

	return MACDData{
		DIF:  dif,
		DEA:  dea,
		MACD: macd,
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

// calculateKDJ calculates KDJ as time series
func calculateKDJ(klines []KLine) KDJData {
	if len(klines) < 9 {
		return KDJData{}
	}

	n := 9
	kValues := make([]float64, 0, len(klines))
	dValues := make([]float64, 0, len(klines))
	jValues := make([]float64, 0, len(klines))

	for i := n - 1; i < len(klines); i++ {
		// Find highest high and lowest low in the window
		high := klines[i].High
		low := klines[i].Low
		for j := i - n + 1; j <= i; j++ {
			if klines[j].High > high {
				high = klines[j].High
			}
			if klines[j].Low < low {
				low = klines[j].Low
			}
		}

		var rsv float64
		if high == low {
			rsv = 50
		} else {
			rsv = (klines[i].Close - low) / (high - low) * 100
		}

		// K = (2/3)*K_prev + (1/3)*RSV
		var k, d float64
		if len(kValues) > 0 {
			k = (2.0*kValues[len(kValues)-1] + rsv) / 3.0
			d = (2.0*dValues[len(dValues)-1] + k) / 3.0
		} else {
			k = 50.0
			d = 50.0
		}
		j := 3*k - 2*d

		kValues = append(kValues, k)
		dValues = append(dValues, d)
		jValues = append(jValues, j)
	}

	return KDJData{
		K: kValues,
		D: dValues,
		J: jValues,
	}
}

// calculateBOLL calculates Bollinger Bands as time series
func calculateBOLL(closes []float64) BOLLData {
	if len(closes) < 20 {
		return BOLLData{}
	}

	upper := make([]float64, 0, len(closes))
	mid := make([]float64, 0, len(closes))
	lower := make([]float64, 0, len(closes))

	for i := 19; i < len(closes); i++ {
		// Calculate 20-day MA
		sum := 0.0
		for j := i - 19; j <= i; j++ {
			sum += closes[j]
		}
		ma20 := sum / 20.0

		// Calculate standard deviation
		var variance float64
		for j := i - 19; j <= i; j++ {
			diff := closes[j] - ma20
			variance += diff * diff
		}
		stdDev := sqrt(variance / 20.0)

		upper = append(upper, ma20+2*stdDev)
		mid = append(mid, ma20)
		lower = append(lower, ma20-2*stdDev)
	}

	return BOLLData{
		Upper: upper,
		Mid:   mid,
		Lower: lower,
	}
}

// calculateWR calculates Williams %R as time series
func calculateWR(klines []KLine, periods []int) []WRData {
	var result []WRData

	for _, period := range periods {
		if len(klines) < period {
			continue
		}
		values := make([]float64, 0, len(klines))

		for i := period - 1; i < len(klines); i++ {
			high := klines[i].High
			low := klines[i].Low
			for j := i - period + 1; j <= i; j++ {
				if klines[j].High > high {
					high = klines[j].High
				}
				if klines[j].Low < low {
					low = klines[j].Low
				}
			}

			var wr float64
			if high == low {
				wr = -50
			} else {
				wr = (high - klines[i].Close) / (high - low) * -100
			}
			values = append(values, wr)
		}
		result = append(result, WRData{
			Period: period,
			Values: values,
		})
	}
	return result
}

// calculateDMI calculates DMI indicator
func calculateDMI(klines []KLine, period int) DMIData {
	if len(klines) < period+1 {
		return DMIData{}
	}

	var trList, plusDMList, minusDMList []float64

	for i := 1; i < len(klines); i++ {
		high := klines[i].High
		low := klines[i].Low
		prevHigh := klines[i-1].High
		prevLow := klines[i-1].Low
		prevClose := klines[i-1].Close

		// True Range
		tr := high - low
		if abs(high-prevClose) > tr {
			tr = abs(high - prevClose)
		}
		if abs(low-prevClose) > tr {
			tr = abs(low - prevClose)
		}
		trList = append(trList, tr)

		// +DM and -DM
		plusDM := 0.0
		minusDM := 0.0

		diffHigh := high - prevHigh
		diffLow := prevLow - low

		if diffHigh > diffLow && diffHigh > 0 {
			plusDM = diffHigh
		}
		if diffLow > diffHigh && diffLow > 0 {
			minusDM = diffLow
		}

		plusDMList = append(plusDMList, plusDM)
		minusDMList = append(minusDMList, minusDM)
	}

	// Wilder's smoothing
	smoothedTR := smoothValues(trList, period)
	smoothedPlusDM := smoothValues(plusDMList, period)
	smoothedMinusDM := smoothValues(minusDMList, period)

	if len(smoothedTR) == 0 || smoothedTR[len(smoothedTR)-1] == 0 {
		return DMIData{}
	}

	latestTR := smoothedTR[len(smoothedTR)-1]
	latestPlusDM := smoothedPlusDM[len(smoothedPlusDM)-1]
	latestMinusDM := smoothedMinusDM[len(smoothedMinusDM)-1]

	plusDI := (latestPlusDM / latestTR) * 100
	minusDI := (latestMinusDM / latestTR) * 100

	return DMIData{
		PlusDI:  plusDI,
		MinusDI: minusDI,
		ADX:     calculateADX(smoothedTR, smoothedPlusDM, smoothedMinusDM, period),
	}
}

// smoothValues applies Wilder's smoothing
func smoothValues(values []float64, period int) []float64 {
	if len(values) < period {
		return values
	}

	result := make([]float64, 0, len(values))
	var sum float64
	for i := 0; i < period; i++ {
		sum += values[i]
	}
	result = append(result, sum)

	for i := period; i < len(values); i++ {
		sum = sum - (sum/float64(period)) + values[i]
		result = append(result, sum)
	}

	return result
}

// calculateADX calculates ADX from smoothed values
func calculateADX(trList, plusDMList, minusDMList []float64, period int) float64 {
	if len(trList) < period*2 {
		return 0
	}

	var dxValues []float64
	for i := period; i < len(trList); i++ {
		if trList[i] == 0 {
			continue
		}
		plusDI := (plusDMList[i] / trList[i]) * 100
		minusDI := (minusDMList[i] / trList[i]) * 100
		dx := abs(plusDI-minusDI) / abs(plusDI+minusDI) * 100
		dxValues = append(dxValues, dx)
	}

	if len(dxValues) < period {
		return 0
	}

	// Smooth DX to get ADX
	var adx float64
	for i := len(dxValues) - period; i < len(dxValues); i++ {
		adx += dxValues[i]
	}
	adx /= float64(period)

	return adx
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

// calculateWMA calculates Weighted Moving Average
func calculateWMA(closes []float64, periods []int) []MAData {
	var result []MAData
	recentCount := len(closes)
	if recentCount > 60 {
		recentCount = 60
	}

	for _, period := range periods {
		if len(closes) < period {
			continue
		}
		values := make([]float64, 0, recentCount)

		for i := period - 1; i < len(closes); i++ {
			sum := 0.0
			weightSum := 0.0
			for j := 0; j < period; j++ {
				weight := float64(period - j)
				sum += closes[i-period+j+1] * weight
				weightSum += weight
			}
			values = append(values, sum/weightSum)
		}
		result = append(result, MAData{
			Period: period,
			Values: values,
		})
	}
	return result
}

// calculateOBV calculates On-Balance Volume
func calculateOBV(klines []KLine) OBVData {
	if len(klines) < 2 {
		return OBVData{}
	}

	values := make([]float64, 0, len(klines))
	obv := float64(0)
	values = append(values, obv)

	for i := 1; i < len(klines); i++ {
		if klines[i].Close > klines[i-1].Close {
			obv += float64(klines[i].Volume)
		} else if klines[i].Close < klines[i-1].Close {
			obv -= float64(klines[i].Volume)
		}
		values = append(values, obv)
	}

	trend := "neutral"
	if len(values) >= 5 {
		recentValues := values[len(values)-5:]
		first := recentValues[0]
		last := recentValues[len(recentValues)-1]
		if last > first*1.02 {
			trend = "rising"
		} else if last < first*0.98 {
			trend = "falling"
		}
	}

	return OBVData{
		Values: values,
		Trend:  trend,
	}
}

// calculateVWAP calculates Volume Weighted Average Price
func calculateVWAP(klines []KLine) VWAPData {
	if len(klines) < 2 {
		return VWAPData{}
	}

	values := make([]float64, 0, len(klines))
	var cumulativeTPV float64
	var cumulativeVolume int64

	for i := 0; i < len(klines); i++ {
		typicalPrice := (klines[i].High + klines[i].Low + klines[i].Close) / 3.0
		cumulativeTPV += typicalPrice * float64(klines[i].Volume)
		cumulativeVolume += klines[i].Volume
		if cumulativeVolume > 0 {
			values = append(values, cumulativeTPV/float64(cumulativeVolume))
		} else {
			values = append(values, typicalPrice)
		}
	}

	return VWAPData{Values: values}
}

// analyzePatterns performs comprehensive pattern analysis
func analyzePatterns(klines []KLine, closes []float64) PatternAnalysis {
	return PatternAnalysis{
		CandlestickPatterns: detectCandlestickPatterns(klines),
		TrendPatterns:       detectTrendPatterns(klines, closes),
	}
}

// detectCandlestickPatterns identifies K-line patterns
func detectCandlestickPatterns(klines []KLine) []CandlestickPattern {
	var patterns []CandlestickPattern

	for i := 2; i < len(klines); i++ {
		if isDoji(klines[i]) {
			patterns = append(patterns, CandlestickPattern{
				Type:     "doji",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  false,
				Meaning:  "犹豫信号 - 可能出现反转",
			})
		}

		if isHammer(klines[i]) {
			patterns = append(patterns, CandlestickPattern{
				Type:     "hammer",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  true,
				Meaning:  "看涨反转 - 买入压力显现",
			})
		}

		if isShootingStar(klines[i]) {
			patterns = append(patterns, CandlestickPattern{
				Type:     "shooting_star",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  false,
				Meaning:  "看跌反转 - 卖出压力显现",
			})
		}

		if isEngulfing(klines[i-2], klines[i-1], klines[i]) {
			bullish := klines[i].Close > klines[i].Open
			patterns = append(patterns, CandlestickPattern{
				Type:     "engulfing",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  bullish,
				Meaning:  "强势反转 - " + map[bool]string{true: "看多", false: "看空"}[bullish] + "吞没形态",
			})
		}

		if isHarami(klines[i-2], klines[i-1], klines[i]) {
			bullish := klines[i].Close > klines[i].Open
			patterns = append(patterns, CandlestickPattern{
				Type:     "harami",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  bullish,
				Meaning:  "趋势衰竭 - " + map[bool]string{true: "看多", false: "看空"}[bullish] + "孕育形态",
			})
		}

		if isPiercingLine(klines[i-2], klines[i-1], klines[i]) {
			patterns = append(patterns, CandlestickPattern{
				Type:     "piercing_line",
				Date:     klines[i].Date,
				Strength: assessPatternStrength(klines[i]),
				Bullish:  true,
				Meaning:  "看涨反转 - 从前期跌势跳空向上",
			})
		}
	}

	return patterns
}

func isDoji(k KLine) bool {
	body := abs(k.Close - k.Open)
	totalRange := k.High - k.Low
	if totalRange == 0 {
		return false
	}
	return body/totalRange < 0.1
}

func isHammer(k KLine) bool {
	body := abs(k.Close - k.Open)
	totalRange := k.High - k.Low
	if totalRange == 0 {
		return false
	}
	lowerShadow := math.Min(k.Close, k.Open) - k.Low
	upperShadow := k.High - math.Max(k.Close, k.Open)
	return lowerShadow > body*2 && upperShadow < body*0.5
}

func isShootingStar(k KLine) bool {
	body := abs(k.Close - k.Open)
	totalRange := k.High - k.Low
	if totalRange == 0 {
		return false
	}
	lowerShadow := math.Min(k.Close, k.Open) - k.Low
	upperShadow := k.High - math.Max(k.Close, k.Open)
	return upperShadow > body*2 && lowerShadow < body*0.5
}

func isEngulfing(prev2, prev1, current KLine) bool {
	prevBody := abs(prev1.Close - prev1.Open)
	currentBody := abs(current.Close - current.Open)
	if prevBody == 0 || currentBody == 0 {
		return false
	}
	if currentBody < prevBody*1.5 {
		return false
	}
	if prev1.Close < prev1.Open && current.Close > current.Open {
		return current.Open < prev1.Close && current.Close > prev1.Open
	}
	if prev1.Close > prev1.Open && current.Close < current.Open {
		return current.Open > prev1.Close && current.Close < prev1.Open
	}
	return false
}

func isHarami(prev2, prev1, current KLine) bool {
	prevBody := abs(prev1.Close - prev1.Open)
	currentBody := abs(current.Close - current.Open)
	if prevBody == 0 || currentBody == 0 {
		return false
	}
	if currentBody > prevBody*0.5 {
		return false
	}
	return current.High < prev1.High && current.Low > prev1.Low
}

func isPiercingLine(prev2, prev1, current KLine) bool {
	if prev1.Close < prev1.Open {
		if current.Close < current.Open {
			return false
		}
		gapDown := current.Open < prev1.Close
		closesAboveMid := current.Close > (prev1.Open+prev1.Close)/2
		return gapDown && closesAboveMid
	}
	return false
}

func assessPatternStrength(k KLine) string {
	body := abs(k.Close - k.Open)
	totalRange := k.High - k.Low
	if totalRange == 0 {
		return "weak"
	}
	ratio := body / totalRange
	if ratio > 0.7 {
		return "strong"
	} else if ratio > 0.4 {
		return "medium"
	}
	return "weak"
}

// detectTrendPatterns identifies chart patterns
func detectTrendPatterns(klines []KLine, closes []float64) []TrendPattern {
	var patterns []TrendPattern

	if dt := detectDoubleTop(klines); dt.Type != "" {
		patterns = append(patterns, dt)
	}
	if db := detectDoubleBottom(klines); db.Type != "" {
		patterns = append(patterns, db)
	}
	if hs := detectHeadShoulders(klines); hs.Type != "" {
		patterns = append(patterns, hs)
	}
	if ihs := detectInverseHeadShoulders(klines); ihs.Type != "" {
		patterns = append(patterns, ihs)
	}
	if t := detectTriangle(klines); t.Type != "" {
		patterns = append(patterns, t)
	}
	if f := detectFlag(klines); f.Type != "" {
		patterns = append(patterns, f)
	}
	if r := detectRectangle(klines); r.Type != "" {
		patterns = append(patterns, r)
	}

	return patterns
}

func detectDoubleTop(klines []KLine) TrendPattern {
	if len(klines) < 40 {
		return TrendPattern{}
	}

	var peaks []int
	for i := 10; i < len(klines)-10; i++ {
		if isPeak(klines, i) {
			peaks = append(peaks, i)
		}
	}

	for i := 0; i < len(peaks)-1; i++ {
		for j := i + 1; j < len(peaks); j++ {
			peak1 := klines[peaks[i]].High
			peak2 := klines[peaks[j]].High
			if peak1 == 0 {
				continue
			}
			diff := abs(peak1-peak2) / peak1
			if diff < 0.03 && peaks[j]-peaks[i] >= 20 {
				return TrendPattern{
					Type:      "double_top",
					StartDate: klines[peaks[i]].Date,
					EndDate:   klines[peaks[j]].Date,
					Strength:  "medium",
					Bullish:   false,
					Breakout:  findSupportBetween(klines, peaks[i], peaks[j]),
					Meaning:   "看跌反转 - 未能两次突破阻力位",
				}
			}
		}
	}
	return TrendPattern{}
}

func detectDoubleBottom(klines []KLine) TrendPattern {
	if len(klines) < 40 {
		return TrendPattern{}
	}

	var bottoms []int
	for i := 10; i < len(klines)-10; i++ {
		if isBottom(klines, i) {
			bottoms = append(bottoms, i)
		}
	}

	for i := 0; i < len(bottoms)-1; i++ {
		for j := i + 1; j < len(bottoms); j++ {
			bottom1 := klines[bottoms[i]].Low
			bottom2 := klines[bottoms[j]].Low
			if bottom1 == 0 {
				continue
			}
			diff := abs(bottom1-bottom2) / bottom1
			if diff < 0.03 && bottoms[j]-bottoms[i] >= 20 {
				return TrendPattern{
					Type:      "double_bottom",
					StartDate: klines[bottoms[i]].Date,
					EndDate:   klines[bottoms[j]].Date,
					Strength:  "medium",
					Bullish:   true,
					Breakout:  findResistanceBetween(klines, bottoms[i], bottoms[j]),
					Meaning:   "看涨反转 - 支撑位确认两次",
				}
			}
		}
	}
	return TrendPattern{}
}

func detectHeadShoulders(klines []KLine) TrendPattern {
	if len(klines) < 60 {
		return TrendPattern{}
	}

	var peaks []int
	for i := 20; i < len(klines)-20; i++ {
		if isPeak(klines, i) {
			peaks = append(peaks, i)
		}
	}

	for _, headIdx := range peaks {
		// Skip if not enough data on left side
		if headIdx < 30 {
			continue
		}
		// Skip if not enough data on right side
		if headIdx >= len(klines)-30 {
			continue
		}

		headHeight := klines[headIdx].High
		var leftShoulder, rightShoulder int
		for i := headIdx - 30; i < headIdx-10; i++ {
			if i >= 0 && i < len(klines) && klines[i].High > headHeight*0.95 {
				leftShoulder = i
				break
			}
		}
		for i := headIdx + 10; i < headIdx+30; i++ {
			if i >= 0 && i < len(klines) && klines[i].High > headHeight*0.95 {
				rightShoulder = i
				break
			}
		}
		if leftShoulder > 0 && rightShoulder > 0 {
			lsHeight := klines[leftShoulder].High
			rsHeight := klines[rightShoulder].High
			if lsHeight != 0 && abs(lsHeight-rsHeight)/lsHeight < 0.05 {
				neckline := (klines[leftShoulder].Low + klines[rightShoulder].Low) / 2
				if klines[headIdx].Low < neckline {
					return TrendPattern{
						Type:      "head_shoulders",
						StartDate: klines[leftShoulder].Date,
						EndDate:   klines[rightShoulder].Date,
						Strength:  "strong",
						Bullish:   false,
						Breakout:  neckline,
						Meaning:   "看跌反转 - 经典反转形态确认",
					}
				}
			}
		}
	}
	return TrendPattern{}
}

func detectInverseHeadShoulders(klines []KLine) TrendPattern {
	if len(klines) < 60 {
		return TrendPattern{}
	}

	var bottoms []int
	for i := 20; i < len(klines)-20; i++ {
		if isBottom(klines, i) {
			bottoms = append(bottoms, i)
		}
	}

	for _, headIdx := range bottoms {
		// Skip if not enough data on left side
		if headIdx < 30 {
			continue
		}
		// Skip if not enough data on right side
		if headIdx >= len(klines)-30 {
			continue
		}

		headDepth := klines[headIdx].Low
		var leftShoulder, rightShoulder int
		for i := headIdx - 30; i < headIdx-10; i++ {
			if i >= 0 && i < len(klines) && klines[i].Low < headDepth*1.05 {
				leftShoulder = i
				break
			}
		}
		for i := headIdx + 10; i < headIdx+30; i++ {
			if i >= 0 && i < len(klines) && klines[i].Low < headDepth*1.05 {
				rightShoulder = i
				break
			}
		}
		if leftShoulder > 0 && rightShoulder > 0 {
			lsDepth := klines[leftShoulder].Low
			rsDepth := klines[rightShoulder].Low
			if lsDepth != 0 && abs(lsDepth-rsDepth)/lsDepth < 0.05 {
				neckline := (klines[leftShoulder].High + klines[rightShoulder].High) / 2
				if klines[headIdx].High > neckline {
					return TrendPattern{
						Type:      "inverse_head_shoulders",
						StartDate: klines[leftShoulder].Date,
						EndDate:   klines[rightShoulder].Date,
						Strength:  "strong",
						Bullish:   true,
						Breakout:  neckline,
						Meaning:   "看涨反转 - 头肩底形态确认",
					}
				}
			}
		}
	}
	return TrendPattern{}
}

func detectTriangle(klines []KLine) TrendPattern {
	if len(klines) < 30 {
		return TrendPattern{}
	}

	recentBars := 20
	if len(klines) < recentBars {
		recentBars = len(klines)
	}

	var highs, lows []float64
	for i := len(klines) - recentBars; i < len(klines); i++ {
		highs = append(highs, klines[i].High)
		lows = append(lows, klines[i].Low)
	}

	highConsistent := true
	for i := 1; i < len(highs); i++ {
		if highs[i] > highs[0]*1.02 {
			highConsistent = false
			break
		}
	}
	lowsRising := true
	for i := 1; i < len(lows); i++ {
		if lows[i] <= lows[i-1] {
			lowsRising = false
			break
		}
	}
	if highConsistent && lowsRising {
		return TrendPattern{
			Type:      "ascending_triangle",
			StartDate: klines[len(klines)-recentBars].Date,
			EndDate:   klines[len(klines)-1].Date,
			Strength:  "medium",
			Bullish:   true,
			Breakout:  highs[0],
			Meaning:   "看涨延续 - 上升三角形有向上突破潜力",
		}
	}

	lowsConsistent := true
	for i := 1; i < len(lows); i++ {
		if lows[i] < lows[0]*0.98 {
			lowsConsistent = false
			break
		}
	}
	highsFalling := true
	for i := 1; i < len(highs); i++ {
		if highs[i] >= highs[i-1] {
			highsFalling = false
			break
		}
	}
	if lowsConsistent && highsFalling {
		return TrendPattern{
			Type:      "descending_triangle",
			StartDate: klines[len(klines)-recentBars].Date,
			EndDate:   klines[len(klines)-1].Date,
			Strength:  "medium",
			Bullish:   false,
			Breakout:  lows[0],
			Meaning:   "看跌延续 - 下降三角形有向下突破潜力",
		}
	}

	return TrendPattern{}
}

func detectFlag(klines []KLine) TrendPattern {
	if len(klines) < 15 {
		return TrendPattern{}
	}

	firstHalf := klines[:len(klines)/2]
	secondHalf := klines[len(klines)/2:]

	if firstHalf[0].Close == 0 || secondHalf[0].Close == 0 {
		return TrendPattern{}
	}
	firstTrend := (firstHalf[len(firstHalf)-1].Close - firstHalf[0].Close) / firstHalf[0].Close
	secondTrend := (secondHalf[len(secondHalf)-1].Close - secondHalf[0].Close) / secondHalf[0].Close

	if abs(firstTrend) > 0.05 && abs(secondTrend) < abs(firstTrend)*0.3 {
		bullish := firstTrend > 0
		return TrendPattern{
			Type:      "flag",
			StartDate: firstHalf[0].Date,
			EndDate:   secondHalf[len(secondHalf)-1].Date,
			Strength:  "medium",
			Bullish:   bullish,
			Breakout:  0,
			Meaning:   map[bool]string{true: "看涨延续 - 旗形暗示向上移动延续", false: "看跌延续 - 旗形暗示向下移动延续"}[bullish],
		}
	}
	return TrendPattern{}
}

func detectRectangle(klines []KLine) TrendPattern {
	if len(klines) < 30 {
		return TrendPattern{}
	}

	recentBars := 20
	startIdx := len(klines) - recentBars
	if startIdx < 0 {
		startIdx = 0
	}

	var highs, lows []float64
	for i := startIdx; i < len(klines); i++ {
		highs = append(highs, klines[i].High)
		lows = append(lows, klines[i].Low)
	}

	maxHigh := highs[0]
	minLow := lows[0]
	for _, h := range highs {
		if h > maxHigh {
			maxHigh = h
		}
	}
	for _, l := range lows {
		if l < minLow {
			minLow = l
		}
	}

	rangeRatio := (maxHigh - minLow) / maxHigh
	if maxHigh == 0 {
		return TrendPattern{}
	}
	if rangeRatio < 0.05 {
		return TrendPattern{
			Type:      "rectangle",
			StartDate: klines[startIdx].Date,
			EndDate:   klines[len(klines)-1].Date,
			Strength:  "medium",
			Bullish:   false,
			Breakout:  maxHigh,
			Meaning:   "整理形态 - 突破方向不确定",
		}
	}
	return TrendPattern{}
}

func isPeak(klines []KLine, idx int) bool {
	if idx < 2 || idx >= len(klines)-2 {
		return false
	}
	return klines[idx].High > klines[idx-1].High &&
		klines[idx].High > klines[idx-2].High &&
		klines[idx].High > klines[idx+1].High &&
		klines[idx].High > klines[idx+2].High
}

func isBottom(klines []KLine, idx int) bool {
	if idx < 2 || idx >= len(klines)-2 {
		return false
	}
	return klines[idx].Low < klines[idx-1].Low &&
		klines[idx].Low < klines[idx-2].Low &&
		klines[idx].Low < klines[idx+1].Low &&
		klines[idx].Low < klines[idx+2].Low
}

func findSupportBetween(klines []KLine, start, end int) float64 {
	minPrice := klines[start].Low
	for i := start; i <= end; i++ {
		if klines[i].Low < minPrice {
			minPrice = klines[i].Low
		}
	}
	return minPrice
}

func findResistanceBetween(klines []KLine, start, end int) float64 {
	maxPrice := klines[start].High
	for i := start; i <= end; i++ {
		if klines[i].High > maxPrice {
			maxPrice = klines[i].High
		}
	}
	return maxPrice
}

// detectSupportResistance identifies support and resistance levels
func detectSupportResistance(klines []KLine) PriceLevels {
	levels := PriceLevels{
		Resistance: []PriceLevel{},
		Support:    []PriceLevel{},
	}

	var resistancePrices, supportPrices []float64

	for i := 5; i < len(klines)-5; i++ {
		isLocalMax := true
		for j := i - 5; j <= i+5; j++ {
			if j != i && klines[j].High >= klines[i].High {
				isLocalMax = false
				break
			}
		}
		if isLocalMax {
			resistancePrices = append(resistancePrices, klines[i].High)
		}

		isLocalMin := true
		for j := i - 5; j <= i+5; j++ {
			if j != i && klines[j].Low <= klines[i].Low {
				isLocalMin = false
				break
			}
		}
		if isLocalMin {
			supportPrices = append(supportPrices, klines[i].Low)
		}
	}

	levels.Resistance = clusterPriceLevels(resistancePrices, "resistance")
	levels.Support = clusterPriceLevels(supportPrices, "support")
	levels.TrendLine = detectTrendLine(klines)
	levels.Channel = detectChannel(klines)

	return levels
}

func clusterPriceLevels(prices []float64, levelType string) []PriceLevel {
	if len(prices) == 0 {
		return []PriceLevel{}
	}

	sorted := make([]float64, len(prices))
	copy(sorted, prices)

	var clusters [][]float64
	for _, p := range sorted {
		found := false
		for i := range clusters {
			if len(clusters[i]) > 0 && abs(clusters[i][0]-p)/clusters[i][0] < 0.02 {
				clusters[i] = append(clusters[i], p)
				found = true
				break
			}
		}
		if !found {
			clusters = append(clusters, []float64{p})
		}
	}

	var levels []PriceLevel
	for _, cluster := range clusters {
		if len(cluster) == 0 {
			continue
		}
		avgPrice := 0.0
		for _, p := range cluster {
			avgPrice += p
		}
		avgPrice /= float64(len(cluster))

		strength := "weak"
		if len(cluster) >= 5 {
			strength = "strong"
		} else if len(cluster) >= 3 {
			strength = "medium"
		}

		levels = append(levels, PriceLevel{
			Price:    avgPrice,
			Strength: strength,
			Touches:  len(cluster),
			Type:     levelType,
		})
	}

	for i := 0; i < len(levels)-1; i++ {
		for j := i + 1; j < len(levels); j++ {
			if levels[j].Touches > levels[i].Touches {
				levels[i], levels[j] = levels[j], levels[i]
			}
		}
	}

	if len(levels) > 5 {
		levels = levels[:5]
	}
	return levels
}

func detectTrendLine(klines []KLine) *TrendLine {
	if len(klines) < 20 {
		return nil
	}

	startIdx := len(klines) - 30
	if startIdx < 0 {
		startIdx = 0
	}

	recentCloses := make([]float64, len(klines)-startIdx)
	for i := startIdx; i < len(klines); i++ {
		recentCloses[i-startIdx] = klines[i].Close
	}

	firstClose := recentCloses[0]
	lastClose := recentCloses[len(recentCloses)-1]
	direction := "up"
	if lastClose < firstClose {
		direction = "down"
	}

	x1, x2 := 0.0, float64(len(recentCloses)-1)
	y1, y2 := recentCloses[0], recentCloses[len(recentCloses)-1]
	if y2 == y1 {
		return nil
	}
	slope := (y2 - y1) / (x2 - x1)
	intercept := y1 - slope*x1

	var residuals float64
	for i, y := range recentCloses {
		predicted := slope*float64(i) + intercept
		residuals += (y - predicted) * (y - predicted)
	}
	meanY := 0.0
	for _, y := range recentCloses {
		meanY += y
	}
	meanY /= float64(len(recentCloses))
	var totalVar float64
	for _, y := range recentCloses {
		totalVar += (y - meanY) * (y - meanY)
	}
	rSquared := 1 - residuals/totalVar

	strength := "weak"
	if rSquared > 0.8 {
		strength = "strong"
	} else if rSquared > 0.5 {
		strength = "medium"
	}

	return &TrendLine{
		Slope:     slope,
		Intercept: intercept,
		StartDate: klines[startIdx].Date,
		EndDate:   klines[len(klines)-1].Date,
		Strength:  strength,
		Direction: direction,
	}
}

func detectChannel(klines []KLine) *Channel {
	if len(klines) < 20 {
		return nil
	}

	startIdx := len(klines) - 30
	if startIdx < 0 {
		startIdx = 0
	}

	midline := detectTrendLine(klines)
	if midline == nil {
		return nil
	}

	maxHighDev := 0.0
	maxLowDev := 0.0
	for i := startIdx; i < len(klines); i++ {
		predicted := midline.Slope*float64(i-startIdx) + midline.Intercept
		highDev := klines[i].High - predicted
		lowDev := predicted - klines[i].Low
		if highDev > maxHighDev {
			maxHighDev = highDev
		}
		if lowDev > maxLowDev {
			maxLowDev = lowDev
		}
	}

	upperLine := TrendLine{
		Slope:     midline.Slope,
		Intercept: midline.Intercept + maxHighDev,
		StartDate: klines[startIdx].Date,
		EndDate:   klines[len(klines)-1].Date,
		Strength:  midline.Strength,
		Direction: midline.Direction,
	}

	lowerLine := TrendLine{
		Slope:     midline.Slope,
		Intercept: midline.Intercept - maxLowDev,
		StartDate: klines[startIdx].Date,
		EndDate:   klines[len(klines)-1].Date,
		Strength:  midline.Strength,
		Direction: midline.Direction,
	}

	return &Channel{
		UpperLine:  upperLine,
		LowerLine:  lowerLine,
		StartDate:  klines[startIdx].Date,
		EndDate:    klines[len(klines)-1].Date,
		Strength:   midline.Strength,
	}
}

// generateRecommendation generates trading recommendation
func generateRecommendation(analysis *TechnicalAnalysis) Recommendation {
	var reasons []string
	buySignals, holdSignals, sellSignals := 0, 0, 0

	for _, ma := range analysis.MA {
		if len(ma.Values) >= 2 {
			current := ma.Values[len(ma.Values)-1]
			prev := ma.Values[len(ma.Values)-2]
			if current > prev {
				buySignals++
			} else if current < prev {
				sellSignals++
			}
		}
	}

	for _, rsi := range analysis.RSI {
		if len(rsi.Values) > 0 {
			current := rsi.Values[len(rsi.Values)-1]
			if current < 30 {
				buySignals += 2
				reasons = append(reasons, fmt.Sprintf("RSI(%d) 超卖 at %.2f", rsi.Period, current))
			} else if current > 70 {
				sellSignals += 2
				reasons = append(reasons, fmt.Sprintf("RSI(%d) 超买 at %.2f", rsi.Period, current))
			} else {
				holdSignals++
			}
		}
	}

	if len(analysis.MACD.DIF) > 0 && len(analysis.MACD.DEA) > 0 {
		dif := analysis.MACD.DIF[len(analysis.MACD.DIF)-1]
		dea := analysis.MACD.DEA[len(analysis.MACD.DEA)-1]

		if dif > dea {
			buySignals += 2
			reasons = append(reasons, "MACD 金叉")
		} else if dif < dea {
			sellSignals += 2
			reasons = append(reasons, "MACD 死叉")
		}
	}

	if len(analysis.KDJ.K) > 0 {
		k := analysis.KDJ.K[len(analysis.KDJ.K)-1]
		d := analysis.KDJ.D[len(analysis.KDJ.D)-1]

		if k < 20 && d < 20 {
			buySignals += 2
			reasons = append(reasons, fmt.Sprintf("KDJ 超卖 (K=%.2f, D=%.2f)", k, d))
		} else if k > 80 && d > 80 {
			sellSignals += 2
			reasons = append(reasons, fmt.Sprintf("KDJ 超买 (K=%.2f, D=%.2f)", k, d))
		}
	}

	if len(analysis.RecentPrices) > 0 && len(analysis.BOLL.Upper) > 0 && len(analysis.BOLL.Lower) > 0 {
		lastPrice := analysis.RecentPrices[len(analysis.RecentPrices)-1].Close
		upper := analysis.BOLL.Upper[len(analysis.BOLL.Upper)-1]
		lower := analysis.BOLL.Lower[len(analysis.BOLL.Lower)-1]
		mid := analysis.BOLL.Mid[len(analysis.BOLL.Mid)-1]

		bandWidth := (upper - lower) / mid * 100
		position := (lastPrice - lower) / (upper - lower) * 100

		if position < 20 {
			buySignals++
			reasons = append(reasons, "价格接近布林带下轨")
		} else if position > 80 {
			sellSignals++
			reasons = append(reasons, "价格接近布林带上轨")
		}

		if bandWidth < 3 {
			holdSignals++
			reasons = append(reasons, "布林带收窄 - 可能即将突破")
		}
	}

	if analysis.DMI.ADX > 25 {
		if analysis.DMI.PlusDI > analysis.DMI.MinusDI {
			buySignals++
			reasons = append(reasons, fmt.Sprintf("DMI 看涨趋势 (ADX=%.2f, +DI=%.2f)", analysis.DMI.ADX, analysis.DMI.PlusDI))
		} else {
			sellSignals++
			reasons = append(reasons, fmt.Sprintf("DMI 看跌趋势 (ADX=%.2f, -DI=%.2f)", analysis.DMI.ADX, analysis.DMI.MinusDI))
		}
	}

	if analysis.OBV.Trend == "rising" {
		buySignals++
		reasons = append(reasons, "OBV 上升趋势 - 成交量确认价格")
	} else if analysis.OBV.Trend == "falling" {
		sellSignals++
		reasons = append(reasons, "OBV 下降趋势 - 成交量与价格背离")
	}

	if len(analysis.RecentPrices) > 0 && len(analysis.VWAP.Values) > 0 {
		lastPrice := analysis.RecentPrices[len(analysis.RecentPrices)-1].Close
		vwap := analysis.VWAP.Values[len(analysis.VWAP.Values)-1]
		if lastPrice > vwap {
			buySignals++
			reasons = append(reasons, "价格高于VWAP - 日内看涨")
		} else {
			sellSignals++
			reasons = append(reasons, "价格低于VWAP - 日内看跌")
		}
	}

	// 去重：每个形态类型只记录一次
	seenCandlestick := make(map[string]bool)
	for _, pattern := range analysis.Patterns.CandlestickPatterns {
		if len(pattern.Date) > 0 && pattern.Type != "" && !seenCandlestick[pattern.Type] {
			seenCandlestick[pattern.Type] = true
			if pattern.Bullish {
				buySignals++
			} else {
				sellSignals++
			}
			reasons = append(reasons, pattern.Meaning)
		}
	}

	seenTrend := make(map[string]bool)
	for _, pattern := range analysis.Patterns.TrendPatterns {
		if pattern.Type != "" && !seenTrend[pattern.Type] {
			seenTrend[pattern.Type] = true
			if pattern.Bullish {
				buySignals++
			} else {
				sellSignals++
			}
			reasons = append(reasons, pattern.Meaning)
		}
	}

	totalSignals := buySignals + holdSignals + sellSignals
	if totalSignals == 0 {
		totalSignals = 1
	}

	var action, summary, riskLevel string
	var confidence float64

	buyRatio := float64(buySignals) / float64(totalSignals) * 100
	sellRatio := float64(sellSignals) / float64(totalSignals) * 100

	if buyRatio > 50 {
		action = "buy"
		confidence = buyRatio
		riskLevel = "medium"
		summary = "多个指标显示看涨动能，买入机会已出现"
	} else if sellRatio > 50 {
		action = "watch"
		confidence = sellRatio
		riskLevel = "high"
		summary = "多个指标显示看跌压力，建议谨慎观察"
	} else {
		action = "hold"
		confidence = 100 - abs(buyRatio-sellRatio)
		riskLevel = "low"
		summary = "出现混合信号，保持当前持仓同时关注更清晰的趋势"
	}

	for _, pattern := range analysis.Patterns.TrendPatterns {
		if pattern.Strength == "strong" {
			confidence += 10
			if confidence > 100 {
				confidence = 100
			}
		}
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "数据不足，无法进行详细分析")
	}

	return Recommendation{
		Action:     action,
		Confidence: confidence,
		Reasons:    reasons,
		Summary:    summary,
		RiskLevel:  riskLevel,
	}
}
