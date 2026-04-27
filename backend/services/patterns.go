package services

// analyzePatterns performs comprehensive pattern analysis
func analyzePatterns(klines []KLine, closes []float64) PatternAnalysis {
	return PatternAnalysis{
		CandlestickPatterns: detectCandlestickPatterns(klines),
		TrendPatterns:       detectTrendPatterns(klines, closes),
	}
}
