package models

import "time"

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Database  string    `json:"database"`
}

type Stock struct {
	Code         string    `json:"code" gorm:"primaryKey;column:code"`
	UserID       string    `json:"userId" gorm:"primaryKey;column:user_id;index"`
	Name         string    `json:"name" gorm:"column:name"`
	CurrentPrice float64   `json:"currentPrice" gorm:"column:current_price"`
	Quantity     int       `json:"quantity" gorm:"column:quantity"`
	BuyPrice     float64   `json:"buyPrice" gorm:"column:buy_price"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Stock) TableName() string {
	return "stocks"
}

type StockRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	CurrentPrice float64 `json:"currentPrice" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	BuyPrice     float64 `json:"buyPrice" binding:"required"`
}

// StockQuote represents real-time stock quote from external API
type StockQuote struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Open        float64 `json:"open"`
	PrevClose   float64 `json:"prevClose"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Current     float64 `json:"current"`
	Volume      int64   `json:"volume"`
	Amount      float64 `json:"amount"`
	UpdateTime  string  `json:"updateTime"`
}

// StockAnalysis represents stock analysis result with position metrics
type StockAnalysis struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"currentPrice"`
	Quantity     int     `json:"quantity"`
	BuyPrice     float64 `json:"buyPrice"`
	MarketValue  float64 `json:"marketValue"`
	Cost         float64 `json:"cost"`
	ProfitLoss   float64 `json:"profitLoss"`
	ProfitRate   float64 `json:"profitRate"`
	Change       float64 `json:"change"`
	ChangeAmount float64 `json:"changeAmount"`
}

// WatchlistItem represents a stock in watchlist (backup table, no user_id)
type WatchlistItem struct {
	Code    string    `json:"code" gorm:"primaryKey;column:code"`
	Name    string    `json:"name" gorm:"column:name"`
	AddedAt time.Time `json:"addedAt" gorm:"column:added_at"`
}

func (WatchlistItem) TableName() string {
	return "watchlist"
}

// UserWatchlist associates users with their watchlist stocks (composite PK: user_id, code)
type UserWatchlist struct {
	UserID  string `json:"userId" gorm:"column:user_id;primaryKey"`
	Code    string `json:"code" gorm:"column:code;primaryKey"`
	AddedAt time.Time `json:"addedAt" gorm:"column:added_at"`
}

func (UserWatchlist) TableName() string {
	return "user_watchlist"
}

// StockDailySnapshot stores daily stock data for historical analysis and backup
type StockDailySnapshot struct {
	Code          string    `json:"code" gorm:"primaryKey;column:code"`
	Date          string    `json:"date" gorm:"primaryKey;column:date"` // Format: 2024-01-01
	Name          string    `json:"name" gorm:"column:name"`
	Open          float64   `json:"open" gorm:"column:open"`
	High          float64   `json:"high" gorm:"column:high"`
	Low           float64   `json:"low" gorm:"column:low"`
	Close         float64   `json:"close" gorm:"column:close"`
	Volume        int64     `json:"volume" gorm:"column:volume"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	TurnoverRate  float64   `json:"turnoverRate" gorm:"column:turnover_rate"` //换手率 %
	MA5           float64   `json:"ma5" gorm:"column:ma5"`
	MA10          float64   `json:"ma10" gorm:"column:ma10"`
	MA20          float64   `json:"ma20" gorm:"column:ma20"`
	MA60          float64   `json:"ma60" gorm:"column:ma60"`
	EMA12         float64   `json:"ema12" gorm:"column:ema12"`
	EMA26         float64   `json:"ema26" gorm:"column:ema26"`
	RSI6          float64   `json:"rsi6" gorm:"column:rsi6"`
	RSI12         float64   `json:"rsi12" gorm:"column:rsi12"`
	RSI24         float64   `json:"rsi24" gorm:"column:rsi24"`
	DIF           float64   `json:"dif" gorm:"column:dif"`
	DEA           float64   `json:"dea" gorm:"column:dea"`
	MACD          float64   `json:"macd" gorm:"column:macd"`
	KDJK          float64   `json:"kdjk" gorm:"column:kdjk"`
	KDJD          float64   `json:"kdjd" gorm:"column:kdjd"`
	KDJJ          float64   `json:"kdjj" gorm:"column:kdjj"`
	BOLLUpper     float64   `json:"bollUpper" gorm:"column:boll_upper"`
	BOLLMid       float64   `json:"bollMid" gorm:"column:boll_mid"`
	BOLLLower     float64   `json:"bollLower" gorm:"column:boll_lower"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (StockDailySnapshot) TableName() string {
	return "stock_daily_snapshots"
}

// User represents a user account
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;column:username;size:100"`
	Password  string    `json:"-" gorm:"column:password;size:255"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

// AISettings stores AI LLM configuration
type AISettings struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"userId" gorm:"column:user_id;index;size:100"`
	APIKey    string    `json:"apiKey" gorm:"column:api_key;size:500"`
	APIURL    string    `json:"apiUrl" gorm:"column:api_url;size:500"`
	Model     string    `json:"model" gorm:"column:model;size:100"`
	GroupID   string    `json:"groupId" gorm:"column:group_id;size:100"`
	Enabled   bool      `json:"enabled" gorm:"column:enabled;default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (AISettings) TableName() string {
	return "ai_settings"
}
