package services

import "math"

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
