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

# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

## 项目定制指南 
- 本项目要求所有面向用户的界面元素必须使用**英文**。
