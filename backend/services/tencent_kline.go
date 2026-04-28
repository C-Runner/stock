package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// fetchTencentKLineData fetches historical daily K-line data from Tencent Finance API
func fetchTencentKLineData(sinaCode string) ([]KLine, error) {
	// Tencent uses same sh/sz prefix format as Sina
	url := fmt.Sprintf("https://web.ifzq.com/app/json/app_kline.php?param=DAY,,%s,180,,,", sinaCode)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://finance.qq.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseTencentKLineData(string(body))
}

// parseTencentKLineData parses Tencent K-line JSON response
// Response format: [["2024-01-02", open, close, high, low, volume], ...]
func parseTencentKLineData(body string) ([]KLine, error) {
	body = strings.TrimSpace(body)
	if body == "" || body == "null" {
		return nil, fmt.Errorf("empty or null response")
	}

	// Parse JSON array
	var rawData [][]interface{}
	if err := json.Unmarshal([]byte(body), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var klines []KLine
	for _, entry := range rawData {
		if len(entry) < 6 {
			continue
		}

		// entry: [date, open, close, high, low, volume]
		kline := KLine{}

		// Date is at index 0
		if dateStr, ok := entry[0].(string); ok {
			kline.Date = dateStr
		} else {
			continue
		}

		// Open at index 1
		if v, ok := entry[1].(float64); ok {
			kline.Open = v
		} else if v, ok := entry[1].(string); ok {
			kline.Open, _ = strconv.ParseFloat(v, 64)
		}

		// Close at index 2
		if v, ok := entry[2].(float64); ok {
			kline.Close = v
		} else if v, ok := entry[2].(string); ok {
			kline.Close, _ = strconv.ParseFloat(v, 64)
		}

		// High at index 3
		if v, ok := entry[3].(float64); ok {
			kline.High = v
		} else if v, ok := entry[3].(string); ok {
			kline.High, _ = strconv.ParseFloat(v, 64)
		}

		// Low at index 4
		if v, ok := entry[4].(float64); ok {
			kline.Low = v
		} else if v, ok := entry[4].(string); ok {
			kline.Low, _ = strconv.ParseFloat(v, 64)
		}

		// Volume at index 5
		if v, ok := entry[5].(float64); ok {
			kline.Volume = int64(v)
		} else if v, ok := entry[5].(string); ok {
			if vol, err := strconv.ParseFloat(v, 64); err == nil {
				kline.Volume = int64(vol)
			}
		}

		if kline.Date != "" {
			klines = append(klines, kline)
		}
	}

	if len(klines) == 0 {
		return nil, fmt.Errorf("no valid k-line data found")
	}

	return klines, nil
}

// CrossValidateKLines validates Sina K-line against Tencent K-line
// Returns a map of date -> validation result for dates that have significant differences
func CrossValidateKLines(sinaKlines, tencentKlines []KLine, tolerance float64) map[string]bool {
	// Build Tencent lookup by date
	tencentByDate := make(map[string]KLine)
	for _, k := range tencentKlines {
		tencentByDate[k.Date] = k
	}

	// Map to track which dates pass validation
	validDates := make(map[string]bool)

	for _, sina := range sinaKlines {
		tencent, exists := tencentByDate[sina.Date]
		if !exists {
			// No Tencent data for this date - skip
			continue
		}

		// Compare close prices (most reliable field)
		if sina.Close > 0 && tencent.Close > 0 {
			diff := abs(sina.Close - tencent.Close)
			avg := (sina.Close + tencent.Close) / 2
			if avg > 0 {
				pctDiff := diff / avg
				if pctDiff <= tolerance {
					validDates[sina.Date] = true
				} else {
					// Significant mismatch
					validDates[sina.Date] = false
				}
				continue
			}
		}

		// Compare high prices
		if sina.High > 0 && tencent.High > 0 {
			diff := abs(sina.High - tencent.High)
			avg := (sina.High + tencent.High) / 2
			if avg > 0 {
				pctDiff := diff / avg
				if pctDiff <= tolerance {
					validDates[sina.Date] = true
				} else {
					validDates[sina.Date] = false
				}
				continue
			}
		}

		// Compare low prices
		if sina.Low > 0 && tencent.Low > 0 {
			diff := abs(sina.Low - tencent.Low)
			avg := (sina.Low + tencent.Low) / 2
			if avg > 0 {
				pctDiff := diff / avg
				if pctDiff <= tolerance {
					validDates[sina.Date] = true
				} else {
					validDates[sina.Date] = false
				}
				continue
			}
		}

		// Compare open prices
		if sina.Open > 0 && tencent.Open > 0 {
			diff := abs(sina.Open - tencent.Open)
			avg := (sina.Open + tencent.Open) / 2
			if avg > 0 {
				pctDiff := diff / avg
				if pctDiff <= tolerance {
					validDates[sina.Date] = true
				} else {
					validDates[sina.Date] = false
				}
				continue
			}
		}

		// If all prices are 0, skip
		validDates[sina.Date] = false
	}

	return validDates
}
