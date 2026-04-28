package services

import (
	"fmt"
	"log"
	"sync"

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
			ma5 = m.Values[len(m.Values)-1]
		case 10:
			ma10 = m.Values[len(m.Values)-1]
		case 20:
			ma20 = m.Values[len(m.Values)-1]
		case 60:
			ma60 = m.Values[len(m.Values)-1]
		}
	}

	var ema12, ema26 float64
	for _, e := range ema {
		switch e.Period {
		case 12:
			ema12 = e.Values[len(e.Values)-1]
		case 26:
			ema26 = e.Values[len(e.Values)-1]
		}
	}

	var rsi6, rsi12, rsi24 float64
	for _, r := range rsi {
		switch r.Period {
		case 6:
			rsi6 = r.Values[len(r.Values)-1]
		case 12:
			rsi12 = r.Values[len(r.Values)-1]
		case 24:
			rsi24 = r.Values[len(r.Values)-1]
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
		DIF:          macd.DIF[len(macd.DIF)-1],
		DEA:          macd.DEA[len(macd.DEA)-1],
		MACD:         macd.MACD[len(macd.MACD)-1],
		KDJK:         kdj.K[len(kdj.K)-1],
		KDJD:         kdj.D[len(kdj.D)-1],
		KDJJ:         kdj.J[len(kdj.J)-1],
		BOLLUpper:    boll.Upper[len(boll.Upper)-1],
		BOLLMid:      boll.Mid[len(boll.Mid)-1],
		BOLLLower:    boll.Lower[len(boll.Lower)-1],
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

	const maxConcurrent = 5
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	failCount := 0

	for _, item := range watchlist {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := BackupStockDaily(code); err != nil {
				log.Printf("Failed to backup %s: %v", code, err)
				mu.Lock()
				failCount++
				mu.Unlock()
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(item.Code)
	}

	wg.Wait()

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

// FetchHistoricalDataForWatchlist fetches up to 180 days of historical data
// for all stocks in a user's watchlist, using two-source validation.
// Only fetches dates not already in the database.
func FetchHistoricalDataForWatchlist(userID string) (int, int, error) {
	var watchlist []models.UserWatchlist
	if err := config.DB.Where("user_id = ?", userID).Find(&watchlist).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to get watchlist: %w", err)
	}

	if len(watchlist) == 0 {
		return 0, 0, nil
	}

	log.Printf("Starting historical data fetch for %d watchlist stocks (user: %s)", len(watchlist), userID)

	const maxConcurrent = 5
	const maxDays = 180
	const tolerance = 0.01 // 1% cross-validation tolerance

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	totalNew := 0
	totalSkipped := 0
	totalFailed := 0

	for _, item := range watchlist {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			inserted, skipped, err := fetchAndSaveHistoricalForStock(code, maxDays, tolerance)
			mu.Lock()
			if err != nil {
				log.Printf("Failed to fetch historical data for %s: %v", code, err)
				totalFailed++
			} else {
				log.Printf("Fetched %d new days for %s (%d skipped, %d failed total)", inserted, code, skipped, totalFailed)
				totalNew += inserted
				totalSkipped += skipped
			}
			mu.Unlock()
		}(item.Code)
	}

	wg.Wait()

	log.Printf("Historical fetch completed: %d new records, %d skipped (already exist), %d stocks failed", totalNew, totalSkipped, totalFailed)
	return totalNew, totalFailed, nil
}

// fetchAndSaveHistoricalForStock fetches historical K-line data for a single stock
// and saves new dates to the database after cross-validation.
func fetchAndSaveHistoricalForStock(code string, maxDays int, tolerance float64) (inserted int, skipped int, err error) {
	sinaCode := convertToSinaCode(code)

	// Get existing dates in DB for this stock
	existingDates := make(map[string]bool)
	var existing []models.StockDailySnapshot
	if err := config.DB.Select("date").Where("code = ?", code).Find(&existing).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to get existing dates: %w", err)
	}
	for _, s := range existing {
		existingDates[s.Date] = true
	}

	// Fetch from primary source (Sina)
	sinaKlines, err := fetchKLineData(sinaCode)
	if err != nil {
		return 0, 0, fmt.Errorf("Sina fetch failed: %w", err)
	}

	// Fetch from secondary source (Tencent)
	tencentKlines, err := fetchTencentKLineData(sinaCode)
	if err != nil {
		log.Printf("Warning: Tencent fetch failed for %s: %v", code, err)
		tencentKlines = nil
	}

	// Cross-validate
	var validDates map[string]bool
	if len(tencentKlines) > 0 {
		validDates = CrossValidateKLines(sinaKlines, tencentKlines, tolerance)
	} else {
		// No Tencent data - use all Sina dates with data
		validDates = make(map[string]bool)
		for _, k := range sinaKlines {
			validDates[k.Date] = true
		}
	}

	// Get stock name
	stockName := code
	if quote, err := SinaFinanceAPI(code); err == nil && quote.Name != "" {
		stockName = quote.Name
	}

	// Process dates from Sina (primary source)
	// Need at least 60 data points for indicators to be meaningful
	minDataPoints := 60
	if len(sinaKlines) < minDataPoints {
		return 0, 0, fmt.Errorf("insufficient k-line data: got %d, need at least %d", len(sinaKlines), minDataPoints)
	}

	for i, kline := range sinaKlines {
		// Check date limit
		if i >= maxDays {
			break
		}

		// Skip if already exists
		if existingDates[kline.Date] {
			skipped++
			continue
		}

		// Skip if not validated (when Tencent data exists)
		if validDates != nil {
			if ok, exists := validDates[kline.Date]; !exists || !ok {
				log.Printf("Skipping %s on %s: cross-validation failed or no Tencent data", code, kline.Date)
				skipped++
				continue
			}
		}

		// Build historical data using all prior klines for indicator calculation
		klinesForCalc := sinaKlines[:i+1]
		closes := make([]float64, len(klinesForCalc))
		for j, k := range klinesForCalc {
			closes[j] = k.Close
		}

		// Skip if not enough data for indicator calculations
		if len(closes) < minDataPoints {
			skipped++
			continue
		}

		ma := calculateMA(closes, []int{5, 10, 20, 60})
		ema := calculateEMA(closes, []int{12, 26})
		rsi := calculateRSI(closes, []int{6, 12, 24})
		macd := calculateMACD(closes)
		kdj := calculateKDJ(klinesForCalc)
		boll := calculateBOLL(closes)

		// Get latest indicator values
		ma5, ma10, ma20, ma60 := getLatestMA(ma)
		ema12, ema26 := getLatestEMA(ema)
		rsi6, rsi12, rsi24 := getLatestRSI(rsi)

		snapshot := models.StockDailySnapshot{
			Code:         code,
			Date:         kline.Date,
			Name:         stockName,
			Open:         kline.Open,
			High:         kline.High,
			Low:          kline.Low,
			Close:        kline.Close,
			Volume:       kline.Volume,
			Amount:       kline.Amount,
			TurnoverRate: 0,
			MA5:          ma5,
			MA10:         ma10,
			MA20:         ma20,
			MA60:         ma60,
			EMA12:        ema12,
			EMA26:        ema26,
			RSI6:         rsi6,
			RSI12:        rsi12,
			RSI24:        rsi24,
			DIF:          macd.DIF[len(macd.DIF)-1],
			DEA:          macd.DEA[len(macd.DEA)-1],
			MACD:         macd.MACD[len(macd.MACD)-1],
			KDJK:         kdj.K[len(kdj.K)-1],
			KDJD:         kdj.D[len(kdj.D)-1],
			KDJJ:         kdj.J[len(kdj.J)-1],
			BOLLUpper:    boll.Upper[len(boll.Upper)-1],
			BOLLMid:      boll.Mid[len(boll.Mid)-1],
			BOLLLower:    boll.Lower[len(boll.Lower)-1],
		}

		if err := config.DB.Save(&snapshot).Error; err != nil {
			log.Printf("Failed to save snapshot for %s on %s: %v", code, kline.Date, err)
			continue
		}
		inserted++
	}

	return inserted, skipped, nil
}

// getLatestMA extracts the latest MA values from MAData slice
func getLatestMA(ma []MAData) (ma5, ma10, ma20, ma60 float64) {
	for _, m := range ma {
		if len(m.Values) == 0 {
			continue
		}
		val := m.Values[len(m.Values)-1]
		switch m.Period {
		case 5:
			ma5 = val
		case 10:
			ma10 = val
		case 20:
			ma20 = val
		case 60:
			ma60 = val
		}
	}
	return
}

// getLatestEMA extracts the latest EMA values from MAData slice
func getLatestEMA(ema []MAData) (ema12, ema26 float64) {
	for _, e := range ema {
		if len(e.Values) == 0 {
			continue
		}
		val := e.Values[len(e.Values)-1]
		switch e.Period {
		case 12:
			ema12 = val
		case 26:
			ema26 = val
		}
	}
	return
}

// getLatestRSI extracts the latest RSI values from RSIData slice
func getLatestRSI(rsi []RSIData) (rsi6, rsi12, rsi24 float64) {
	for _, r := range rsi {
		if len(r.Values) == 0 {
			continue
		}
		val := r.Values[len(r.Values)-1]
		switch r.Period {
		case 6:
			rsi6 = val
		case 12:
			rsi12 = val
		case 24:
			rsi24 = val
		}
	}
	return
}