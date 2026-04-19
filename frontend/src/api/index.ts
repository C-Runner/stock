const API_BASE_URL = 'http://localhost:8080'

export interface HealthResponse {
  status: string
  timestamp: string
  service: string
}

// Watchlist types
export interface WatchlistItem {
  code: string
  name: string
  addedAt: string
}

export interface WatchlistQuote {
  code: string
  name: string
  addedAt: string
  current: number
  open: number
  high: number
  low: number
  change: number
  changeRate: number
  volume: number
  updateTime: string
}

export interface Stock {
  code: string
  name: string
  currentPrice: number
  quantity: number
  buyPrice: number
  createdAt: string
  updatedAt: string
}

export interface StockRequest {
  code: string
  name: string
  currentPrice: number
  quantity: number
  buyPrice: number
}

export interface StockQuote {
  code: string
  name: string
  open: number
  high: number
  low: number
  current: number
  volume: number
  amount: number
  updateTime: string
}

export interface StockAnalysis {
  code: string
  name: string
  currentPrice: number
  quantity: number
  buyPrice: number
  marketValue: number
  cost: number
  profitLoss: number
  profitRate: number
  change: number
  changeAmount: number
}

export interface MAData {
  period: number
  value: number
}

export interface RSIData {
  period: number
  value: number
}

export interface MACDData {
  dif: number
  dea: number
  macd: number
}

export interface KDJData {
  k: number
  d: number
  j: number
}

export interface BOLLData {
  upper: number
  mid: number
  lower: number
}

export interface PricePoint {
  date: string
  open: number
  close: number
  high: number
  low: number
  volume: number
}

export interface TechnicalAnalysis {
  code: string
  name: string
  ma: MAData[]
  ema: MAData[]
  rsi: RSIData[]
  macd: MACDData
  kdj: KDJData
  boll: BOLLData
  recentPrices: PricePoint[]
}

const getAuthHeaders = () => {
  const token = localStorage.getItem('token')
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

export const api = {
  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    if (response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpiry')
      window.location.href = '/login'
      throw new Error('Unauthorized')
    }
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`)
    }
    return response.json()
  },

  async post<T>(path: string, data: unknown): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })
    if (response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpiry')
      window.location.href = '/login'
      throw new Error('Unauthorized')
    }
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`)
    }
    return response.json()
  },

  async delete<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })
    if (response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpiry')
      window.location.href = '/login'
      throw new Error('Unauthorized')
    }
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`)
    }
    return response.json()
  }
}

export const healthApi = {
  check: () => api.get<HealthResponse>('/health')
}

export const stockApi = {
  getStocks: () => api.get<Stock[]>('/api/stocks'),
  createStock: (data: StockRequest) => api.post<Stock>('/api/stocks', data),
  deleteStock: (code: string) => api.delete<{ message: string }>(`/api/stocks/${code}`),
  getQuote: (code: string) => api.get<StockQuote>(`/api/stocks/quote/${code}`),
  getAnalysis: (code: string) => api.get<StockAnalysis>(`/api/stocks/analysis/${code}`),
  getTechnical: (code: string) => api.get<TechnicalAnalysis>(`/api/stocks/technical/${code}`),
  searchStocks: (q: string) => api.get<Stock[]>(`/api/stocks/search?q=${encodeURIComponent(q)}`)
}

export const watchlistApi = {
  getWatchlist: () => api.get<WatchlistItem[]>('/api/watchlist'),
  addToWatchlist: (code: string, name: string) =>
    api.post<WatchlistItem>('/api/watchlist', { code, name }),
  removeFromWatchlist: (code: string) =>
    api.delete<{ message: string }>(`/api/watchlist/${code}`)
}
