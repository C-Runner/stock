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
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Current    float64 `json:"current"`
	Volume     int64   `json:"volume"`
	Amount     float64 `json:"amount"`
	UpdateTime string  `json:"updateTime"`
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

// WatchlistItem represents a stock in user's watchlist
type WatchlistItem struct {
	Code    string    `json:"code" gorm:"primaryKey;column:code"`
	Name    string    `json:"name" gorm:"column:name"`
	AddedAt time.Time `json:"addedAt" gorm:"column:added_at"`
}

func (WatchlistItem) TableName() string {
	return "watchlist"
}

// WatchlistQuote combines watchlist item with real-time quote
type WatchlistQuote struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	AddedAt    string  `json:"addedAt"`
	Current    float64 `json:"current"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Change     float64 `json:"change"`
	ChangeRate float64 `json:"changeRate"`
	Volume     int64   `json:"volume"`
	UpdateTime string  `json:"updateTime"`
}
