# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stock analysis system with Go/Gin backend and Vue 3/TypeScript frontend. Uses Sina Finance API for real-time stock data and technical indicators.

## Architecture

```
backend/
├── main.go           # Entry point, Gin router setup, JWT auth middleware
├── config/            # Configuration and PostgreSQL connection
├── handlers/         # HTTP handlers (auth, stocks, watchlist, analysis)
├── models/           # GORM models (Stock, WatchlistItem)
├── services/        # Business logic (stock_api.go, technical_analysis.go)

frontend/
├── src/
│   ├── api/          # TypeScript API client
│   ├── views/        # Page components (Login, Home, Analysis, Watchlist)
│   ├── components/   # Reusable components (KLineChart, StockSearch)
│   └── router/       # Vue Router with auth guards
└── vite.config.ts
```

## Commands

### Backend
```bash
cd backend && go run main.go          # Start server (port 8080)
```

### Frontend
```bash
cd frontend
npm install                            # Install dependencies
npm run dev                            # Start dev server (port 5173)
npm run build                          # Production build
```

## API Design

- **Authentication**: JWT-based, token sent as `Authorization: Bearer <token>`
- **Protected routes**: All `/api/*` routes require valid JWT
- **Public routes**: `/health`, `/api/login`

### Key Endpoints
- `GET /api/stocks` - List user's stocks
- `GET /api/stocks/technical/:code` - K-line + technical indicators (MA, EMA, RSI, MACD, KDJ, BOLL)
- `GET /api/watchlist` - User's watchlist
- `GET /api/stocks/search?q=` - Search stocks

## Data Models

- **Stock**: User's holdings with code, name, quantity, buy price
- **WatchlistItem**: Simplified watchlist with code, name, addedAt
- **StockQuote**: Real-time quote from Sina (open, high, low, current, volume)
- **TechnicalAnalysis**: K-line data + calculated indicators (MA, RSI, MACD, KDJ, BOLL)

## Technical Indicators (backend/services/technical_analysis.go)

K-line data fetched from Sina Finance API, then calculated client-side:
- **MA/EMA**: Moving averages (5, 10, 20, 60 periods)
- **RSI**: Relative Strength Index (6, 12, 24 periods)
- **MACD**: 12/26-day EMA with 9-day signal line
- **KDJ**: Stochastic oscillator with K, D, J lines
- **BOLL**: Bollinger Bands (20-day MA ± 2 std dev)

## Frontend Routing

- `/login` - Login page (public)
- `/` - Home page (protected)
- `/watchlist` - Watchlist (protected)
- `/analysis/:code` - Stock analysis with K-line chart (protected)

Router guard validates JWT token and expiry from localStorage before allowing access to protected routes.
