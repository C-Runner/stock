package services

import (
	"fmt"
	"log"
	"time"

	"backend/config"
	"backend/models"
)

// BackupStockDaily backs up daily data for a single stock with cross-validation
func BackupStockDaily(code string) error {
	sinaCode := convertToSinaCode(code)
	klines, err := fetchKLineData(sinaCode)
	if err != nil {
		return fmt.Errorf("failed to fetch K-line data: %w", err)
	}

	if len(klines) < 60 {
		return fmt.Errorf("insufficient K-line data for %s", code)
	}

	// Get the latest K-line (today's or last trading day)
	latest := klines[len(klines)-1]

	// Get quote from Sina (primary source)
	sinaQuote, err := SinaFinanceAPI(code)
	if err != nil {
		return fmt.Errorf("failed to get Sina quote: %w", err)
	}

	// Cross-validate with Tencent Finance (secondary source)
	validationResult, tencentQuote, err := ValidateWithSecondSource(code, sinaQuote)
	if err != nil {
		log.Printf("Warning: cross-validation failed for %s: %v", code, err)
	} else if validationResult != nil {
		if !validationResult.IsValid {
			log.Printf("Warning: data mismatch detected for %s, max diff: %.2f%%", code, validationResult.MaxDiff*100)
			// Log detailed differences
			for field, diff := range validationResult.FieldsDiff {
				log.Printf("  %s difference: %.2f%%", field, diff*100)
			}
		} else {
			log.Printf("Cross-validation passed for %s (max diff: %.2f%%)", code, validationResult.MaxDiff*100)
		}
		_ = tencentQuote // tencentQuote available if needed for further processing
	}

	// Calculate indicators using all available data
	closes := make([]float64, len(klines))
	for i, k := range klines {
		closes[i] = k.Close
	}

	ma := calculateMA(closes, []int{5, 10, 20, 60})
	ema := calculateEMA(closes, []int{12, 26})
	rsi := calculateRSI(closes, []int{6, 12, 24})
	macd := calculateMACD(closes)
	kdj := calculateKDJ(klines)
	boll := calculateBOLL(closes)

	// Find MA values
	var ma5, ma10, ma20, ma60 float64
	for _, m := range ma {
		switch m.Period {
		case 5:
			ma5 = m.Value
		case 10:
			ma10 = m.Value
		case 20:
			ma20 = m.Value
		case 60:
			ma60 = m.Value
		}
	}

	var ema12, ema26 float64
	for _, e := range ema {
		switch e.Period {
		case 12:
			ema12 = e.Value
		case 26:
			ema26 = e.Value
		}
	}

	var rsi6, rsi12, rsi24 float64
	for _, r := range rsi {
		switch r.Period {
		case 6:
			rsi6 = r.Value
		case 12:
			rsi12 = r.Value
		case 24:
			rsi24 = r.Value
		}
	}

	// Calculate turnover rate (volume / total shares outstanding)
	// Simplified: use volume ratio compared to average
	var turnoverRate float64
	if latest.Volume > 0 && sinaQuote.Volume > 0 {
		// Estimate turnover rate as percentage
		avgVolume := float64(latest.Volume)
		if len(klines) >= 20 {
			sum := int64(0)
			for i := len(klines) - 20; i < len(klines)-1; i++ {
				sum += klines[i].Volume
			}
			avgVolume = float64(sum) / 20
		}
		if avgVolume > 0 {
			turnoverRate = (float64(latest.Volume) / avgVolume) * 100
		}
	}

	snapshot := models.StockDailySnapshot{
		Code:         code,
		Date:         latest.Date,
		Name:         sinaQuote.Name,
		Open:         latest.Open,
		High:         latest.High,
		Low:          latest.Low,
		Close:        latest.Close,
		Volume:       latest.Volume,
		Amount:       latest.Amount,
		TurnoverRate: turnoverRate,
		MA5:          ma5,
		MA10:         ma10,
		MA20:         ma20,
		MA60:         ma60,
		EMA12:        ema12,
		EMA26:        ema26,
		RSI6:         rsi6,
		RSI12:        rsi12,
		RSI24:        rsi24,
		DIF:          macd.DIF,
		DEA:          macd.DEA,
		MACD:         macd.MACD,
		KDJK:         kdj.K,
		KDJD:         kdj.D,
		KDJJ:         kdj.J,
		BOLLUpper:    boll.Upper,
		BOLLMid:      boll.Mid,
		BOLLLower:    boll.Lower,
	}

	// Upsert: update if exists, insert if not
	result := config.DB.Save(&snapshot)
	if result.Error != nil {
		return fmt.Errorf("failed to save snapshot: %w", result.Error)
	}

	log.Printf("Backed up daily data for %s on %s", code, latest.Date)
	return nil
}

// BackupAllWatchlist backs up daily data for all stocks in watchlist
func BackupAllWatchlist() error {
	var watchlist []models.WatchlistItem
	if err := config.DB.Find(&watchlist).Error; err != nil {
		return fmt.Errorf("failed to get watchlist: %w", err)
	}

	log.Printf("Starting daily backup for %d watchlist stocks", len(watchlist))
	successCount := 0
	failCount := 0

	for _, item := range watchlist {
		if err := BackupStockDaily(item.Code); err != nil {
			log.Printf("Failed to backup %s: %v", item.Code, err)
			failCount++
		} else {
			successCount++
		}
		// Rate limit to avoid API restrictions
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Daily backup completed: %d succeeded, %d failed", successCount, failCount)
	if failCount > 0 {
		return fmt.Errorf("%d stocks failed to backup", failCount)
	}
	return nil
}

// GetDailySnapshots retrieves historical daily snapshots for a stock
func GetDailySnapshots(code string, limit int) ([]models.StockDailySnapshot, error) {
	var snapshots []models.StockDailySnapshot
	query := config.DB.Where("code = ?", code).Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&snapshots).Error; err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	return snapshots, nil
}

// GetDailySnapshot retrieves a specific day's snapshot
func GetDailySnapshot(code string, date string) (*models.StockDailySnapshot, error) {
	var snapshot models.StockDailySnapshot
	if err := config.DB.Where("code = ? AND date = ?", code, date).First(&snapshot).Error; err != nil {
		return nil, fmt.Errorf("snapshot not found: %w", err)
	}
	return &snapshot, nil
}