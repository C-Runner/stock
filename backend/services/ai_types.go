package services

import "backend/models"

// AIAnalysisReport represents the structured AI analysis report
type AIAnalysisReport struct {
	Code         string          `json:"code"`
	Name         string          `json:"name"`
	GeneratedAt  string          `json:"generatedAt"`
	Cached       bool            `json:"cached"`
	FromCache    bool            `json:"fromCache"`

	// Multi-dimensional scoring
	Scores       Scores          `json:"scores"`

	// Key findings
	KeyFindings  KeyFindings     `json:"keyFindings"`

	// Similar patterns
	SimilarPatterns []SimilarPattern `json:"similarPatterns"`

	// Institutional sentiment
	InstitutionalSentiment InstitutionalSentiment `json:"institutionalSentiment"`
}

// Scores represents multi-dimensional scoring
type Scores struct {
	Technical     DimensionScore `json:"technical"`     // 30%
	Fundamental   DimensionScore `json:"fundamental"`   // 30%
	MoneyFlow     DimensionScore `json:"moneyFlow"`     // 25%
	NewsSentiment DimensionScore `json:"newsSentiment"` // 15%

	CompositeScore float64 `json:"compositeScore"` // Weighted average
	AnxietyIndex   float64 `json:"anxietyIndex"`  // 0-100
	AttentionLevel string  `json:"attentionLevel"` // "high", "medium", "low"
}

// DimensionScore represents a single dimension score
type DimensionScore struct {
	Score     float64 `json:"score"`     // 0-100
	Trend     string  `json:"trend"`     // "improving", "stable", "declining"
	Summary   string  `json:"summary"`   // Brief English summary
	Factors   []string `json:"factors"`  // Contributing factors
}

// KeyFindings represents key findings from AI analysis
type KeyFindings struct {
	Highlights []Highlight `json:"highlights"` // 3 positive highlights
	Risks      []Risk      `json:"risks"`      // 2 risk factors
}

// Highlight represents a positive finding
type Highlight struct {
	Title   string `json:"title"`   // e.g., "ROE continuous growth"
	Context string `json:"context"` // e.g., "ROE连续3年下滑，但北向资金最近5日逆势加仓"
}

// Risk represents a risk factor
type Risk struct {
	Title   string `json:"title"`   // e.g., "Margin pressure"
	Context string `json:"context"` // Detailed risk description
}

// SimilarPattern represents a similar historical K-line pattern
type SimilarPattern struct {
	PatternID      string  `json:"patternId"`      // Unique pattern identifier
	StartDate      string  `json:"startDate"`      // Pattern start date
	EndDate        string  `json:"endDate"`        // Pattern end date
	Similarity     float64 `json:"similarity"`     // 0-100 similarity percentage
	PriceChange    float64 `json:"priceChange"`    // Price change during pattern
	Next5DayWinRate float64 `json:"next5DayWinRate"` // Win rate in next 5 trading days
	Next20DayWinRate float64 `json:"next20DayWinRate"` // Win rate in next 20 trading days
}

// InstitutionalSentiment represents institutional research sentiment
type InstitutionalSentiment struct {
	Period          string  `json:"period"`          // "30d" for last 30 days
	TotalReports    int     `json:"totalReports"`   // Total report count
	BuyRatio        float64 `json:"buyRatio"`       // Percentage of "buy" ratings
	HoldRatio       float64 `json:"holdRatio"`      // Percentage of "hold" ratings
	SellRatio       float64 `json:"sellRatio"`      // Percentage of "sell" ratings
	RatioChange     string  `json:"ratioChange"`    // "improving", "stable", "deteriorating"
	ConsensusTarget float64 `json:"consensusTarget"` // Average price target
	ConsensusRating string  `json:"consensusRating"` // "buy", "hold", "sell"
}

// AIAnalysisInput is the input data for AI analysis
type AIAnalysisInput struct {
	Code           string
	Name           string
	Quote          *StockQuote
	Technical      *TechnicalAnalysis
	HistoricalData []KLine // Full historical data for pattern matching
}

// StockQuote is an alias for models.StockQuote for convenience
type StockQuote = models.StockQuote
