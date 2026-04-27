package services

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
