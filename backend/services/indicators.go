package services

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

		for i := period - 1; i < len(closes); i++ {
			if i == period-1 {
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

	for i := 8; i < len(dif); i++ {
		window := dif[i-8 : i+1]
		sum := 0.0
		for _, v := range window {
			sum += v
		}
		dea = append(dea, sum/9.0)
	}

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
		sum := 0.0
		for j := i - 19; j <= i; j++ {
			sum += closes[j]
		}
		ma20 := sum / 20.0

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

		tr := high - low
		if abs(high-prevClose) > tr {
			tr = abs(high - prevClose)
		}
		if abs(low-prevClose) > tr {
			tr = abs(low - prevClose)
		}
		trList = append(trList, tr)

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

	var adx float64
	for i := len(dxValues) - period; i < len(dxValues); i++ {
		adx += dxValues[i]
	}
	adx /= float64(period)

	return adx
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
