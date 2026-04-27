package services

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
	Code            string              `json:"code"`
	Name            string              `json:"name"`
	MA              []MAData            `json:"ma"`              // Moving averages
	EMA             []MAData            `json:"ema"`             // Exponential moving averages
	WMA             []MAData            `json:"wma"`             // Weighted moving averages
	RSI             []RSIData           `json:"rsi"`             // Relative Strength Index
	MACD            MACDData            `json:"macd"`            // MACD indicator
	KDJ             KDJData             `json:"kdj"`             // KDJ indicator
	BOLL            BOLLData            `json:"boll"`            // Bollinger Bands
	WR              []WRData            `json:"wr"`              // Williams %R
	DMI             DMIData             `json:"dmi"`             // Directional Movement Index
	OBV             OBVData             `json:"obv"`             // On-Balance Volume
	VWAP            VWAPData            `json:"vwap"`            // Volume Weighted Average Price
	Patterns        PatternAnalysis     `json:"patterns"`        // Pattern recognition results
	Levels          PriceLevels        `json:"levels"`          // Support/Resistance levels
	Recommendation  Recommendation     `json:"recommendation"` // Buy/Hold/Watch recommendation
	RecentPrices    []PricePoint       `json:"recentPrices"`    // Recent price data for chart
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

// OBVData represents On-Balance Volume
type OBVData struct {
	Values []float64 `json:"values"` // OBV values time series
	Trend  string    `json:"trend"`  // "rising", "falling", "neutral"
}

// VWAPData represents Volume Weighted Average Price
type VWAPData struct {
	Values []float64 `json:"values"` // VWAP values time series
}

// PatternAnalysis holds all pattern recognition results
type PatternAnalysis struct {
	CandlestickPatterns []CandlestickPattern `json:"candlestickPatterns"` // K-line patterns
	TrendPatterns       []TrendPattern      `json:"trendPatterns"`      // Trend/chart patterns
}

// CandlestickPattern represents a recognized candlestick pattern
type CandlestickPattern struct {
	Type     string `json:"type"`     // Pattern name
	Date     string `json:"date"`     // Date when pattern was detected
	Strength string `json:"strength"` // "strong", "medium", "weak"
	Bullish  bool   `json:"bullish"`  // True if bullish pattern
	Meaning  string `json:"meaning"`  // Pattern meaning in English
}

// TrendPattern represents a recognized chart pattern
type TrendPattern struct {
	Type      string  `json:"type"`      // "double_top", "head_shoulders", etc.
	StartDate string  `json:"startDate"` // Pattern start date
	EndDate   string  `json:"endDate"`   // Pattern end date
	Strength  string  `json:"strength"`  // "strong", "medium", "weak"
	Bullish   bool    `json:"bullish"`   // True if bullish pattern
	Breakout  float64 `json:"breakout"` // Breakout price level
	Meaning   string  `json:"meaning"`   // Pattern meaning in English
}

// PriceLevels represents support and resistance levels
type PriceLevels struct {
	Resistance []PriceLevel `json:"resistance"` // Resistance levels
	Support    []PriceLevel `json:"support"`    // Support levels
	TrendLine  *TrendLine   `json:"trendLine"`  // Current trend line
	Channel    *Channel     `json:"channel"`    // Price channel if detected
}

// PriceLevel represents a single price level
type PriceLevel struct {
	Price    float64 `json:"price"`    // Price level
	Strength string  `json:"strength"` // "strong", "medium", "weak"
	Touches  int     `json:"touches"`  // Number of times price touched
	Type     string  `json:"type"`     // "resistance" or "support"
}

// TrendLine represents a detected trend line
type TrendLine struct {
	Slope     float64 `json:"slope"`     // Slope of the trend line
	Intercept float64 `json:"intercept"` // Y-intercept
	StartDate string  `json:"startDate"` // Trend line start date
	EndDate   string  `json:"endDate"`   // Trend line end date
	Strength  string  `json:"strength"`  // "strong", "medium", "weak"
	Direction string  `json:"direction"` // "up" or "down"
}

// Channel represents a price channel
type Channel struct {
	UpperLine TrendLine `json:"upperLine"` // Upper boundary
	LowerLine TrendLine `json:"lowerLine"` // Lower boundary
	StartDate string    `json:"startDate"` // Channel start date
	EndDate   string    `json:"endDate"`   // Channel end date
	Strength  string    `json:"strength"`  // "strong", "medium", "weak"
}

// Recommendation represents trading recommendation
type Recommendation struct {
	Action     string   `json:"action"`     // "buy", "hold", "watch"
	Confidence float64  `json:"confidence"` // 0-100 confidence score
	Reasons    []string `json:"reasons"`   // Reasons for the recommendation
	Summary    string   `json:"summary"`   // English summary of the analysis
	RiskLevel  string   `json:"riskLevel"`  // "low", "medium", "high"
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
