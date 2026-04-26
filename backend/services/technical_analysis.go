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
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	MA           []MAData     `json:"ma"`           // Moving averages
	EMA          []MAData     `json:"ema"`           // Exponential moving averages
	RSI          []RSIData    `json:"rsi"`           // Relative Strength Index
	MACD         MACDData     `json:"macd"`          // MACD indicator
	KDJ          KDJData      `json:"kdj"`           // KDJ indicator
	BOLL         BOLLData     `json:"boll"`          // Bollinger Bands
	WR           []WRData     `json:"wr"`            // Williams %R
	DMI          DMIData      `json:"dmi"`           // Directional Movement Index
	RecentPrices []PricePoint `json:"recentPrices"`  // Recent price data for chart
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
