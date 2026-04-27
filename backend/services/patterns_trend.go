package services

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
		if headIdx < 30 {
			continue
		}
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
		if headIdx < 30 {
			continue
		}
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
