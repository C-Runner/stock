package models

import "time"

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

type Stock struct {
	Code         string    `json:"code" bson:"code"`
	Name         string    `json:"name" bson:"name"`
	CurrentPrice float64   `json:"currentPrice" bson:"currentPrice"`
	Quantity     int       `json:"quantity" bson:"quantity"`
	BuyPrice     float64   `json:"buyPrice" bson:"buyPrice"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
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
