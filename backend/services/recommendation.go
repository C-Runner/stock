package services

import "fmt"

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
