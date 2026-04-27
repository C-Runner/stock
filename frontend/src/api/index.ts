const API_BASE_URL = ''

export interface HealthResponse {
  status: string
  timestamp: string
  service: string
}

// Stock types
export interface Stock {
  code: string
  name: string
  currentPrice: number
  quantity: number
  buyPrice: number
  createdAt: string
  updatedAt: string
}

export interface WatchlistItem {
  code: string
  name: string
  addedAt: string
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
  prevClose: number
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

export interface NoPositionResponse {
  hasPosition: false
}

export type StockAnalysisResponse = StockAnalysis | NoPositionResponse

export interface MAData {
  period: number
  values: number[]
}

export interface RSIData {
  period: number
  values: number[]
}

export interface MACDData {
  dif: number[]
  dea: number[]
  macd: number[]
}

export interface KDJData {
  k: number[]
  d: number[]
  j: number[]
}

export interface BOLLData {
  upper: number[]
  mid: number[]
  lower: number[]
}

export interface WRData {
  period: number
  values: number[]
}

export interface DMIData {
  plusDI: number
  minusDI: number
  adx: number
}

export interface OBVData {
  values: number[]
  trend: string
}

export interface VWAPData {
  values: number[]
}

export interface PricePoint {
  date: string
  open: number
  close: number
  high: number
  low: number
  volume: number
}

export interface CandlestickPattern {
  type: string
  date: string
  strength: string
  bullish: boolean
  meaning: string
}

export interface TrendPattern {
  type: string
  startDate: string
  endDate: string
  strength: string
  bullish: boolean
  breakout: number
  meaning: string
}

export interface PatternAnalysis {
  candlestickPatterns: CandlestickPattern[]
  trendPatterns: TrendPattern[]
}

export interface PriceLevel {
  price: number
  strength: string
  touches: number
  type: string
}

export interface TrendLine {
  slope: number
  intercept: number
  startDate: string
  endDate: string
  strength: string
  direction: string
}

export interface Channel {
  upperLine: TrendLine
  lowerLine: TrendLine
  startDate: string
  endDate: string
  strength: string
}

export interface PriceLevels {
  resistance: PriceLevel[]
  support: PriceLevel[]
  trendLine: TrendLine | null
  channel: Channel | null
}

export interface Recommendation {
  action: string
  confidence: number
  reasons: string[]
  summary: string
  riskLevel: string
}

export interface TechnicalAnalysis {
  code: string
  name: string
  ma: MAData[]
  ema: MAData[]
  wma: MAData[]
  rsi: RSIData[]
  macd: MACDData
  kdj: KDJData
  boll: BOLLData
  wr: WRData[]
  dmi: DMIData
  obv: OBVData
  vwap: VWAPData
  patterns: PatternAnalysis
  levels: PriceLevels
  recommendation: Recommendation
  recentPrices: PricePoint[]
}

export interface StockDailySnapshot {
  code: string
  date: string
  name: string
  open: number
  high: number
  low: number
  close: number
  volume: number
  amount: number
  turnoverRate: number
  ma5: number
  ma10: number
  ma20: number
  ma60: number
  ema12: number
  ema26: number
  rsi6: number
  rsi12: number
  rsi24: number
  dif: number
  dea: number
  macd: number
  kdjk: number
  kdjd: number
  kdjj: number
  bollUpper: number
  bollMid: number
  bollLower: number
  createdAt: string
}

export interface BackupResponse {
  message: string
  status: string
  started_at?: string
  code?: string
  backupDate?: string
}

export interface DailySnapshotsResponse {
  code: string
  count: number
  snapshots: StockDailySnapshot[]
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
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 10000)
    try {
      const response = await fetch(`${API_BASE_URL}${path}`, {
        method: 'GET',
        headers: getAuthHeaders(),
        signal: controller.signal
      })
      clearTimeout(timeout)
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
    } catch (e) {
      clearTimeout(timeout)
      throw e
    }
  },

  async post<T>(path: string, data: unknown): Promise<T> {
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 10000)
    try {
      const response = await fetch(`${API_BASE_URL}${path}`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(data),
        signal: controller.signal
      })
      clearTimeout(timeout)
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
    } catch (e) {
      clearTimeout(timeout)
      throw e
    }
  },

  async delete<T>(path: string): Promise<T> {
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 10000)
    try {
      const response = await fetch(`${API_BASE_URL}${path}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
        signal: controller.signal
      })
      clearTimeout(timeout)
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
    } catch (e) {
      clearTimeout(timeout)
      throw e
    }
  }
}

export const healthApi = {
  check: () => api.get<{ status: string; timestamp: string; service: string }>('/health')
}

export const authApi = {
  login: (username: string, password: string) =>
    api.post<{ token: string; expiresAt: number; user: { id: number; username: string; createdAt: string } }>('/api/login', { username, password }),
  logout: () => api.post<{ message: string }>('/api/logout', {})
}

export const stockApi = {
  getStocks: () => api.get<Stock[]>('/api/stocks'),
  createStock: (data: StockRequest) => api.post<Stock>('/api/stocks', data),
  deleteStock: (code: string) => api.delete<{ message: string }>(`/api/stocks/${code}`),
  getQuote: (code: string) => api.get<StockQuote>(`/api/stocks/quote/${code}`),
  getAnalysis: (code: string) => api.get<StockAnalysisResponse>(`/api/stocks/analysis/${code}`),
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

export const backupApi = {
  // Manual backup trigger (async - returns immediately)
  triggerBackup: () => api.post<BackupResponse>('/api/stocks/backup', {}),
  // Backup single stock
  backupStock: (code: string) =>
    api.post<BackupResponse>(`/api/stocks/backup/${code}`, {}),
  // Get historical daily snapshots
  getDailySnapshots: (code: string, limit = 60) =>
    api.get<DailySnapshotsResponse>(`/api/stocks/daily/${code}?limit=${limit}`),
  // Get specific date snapshot
  getDailySnapshot: (code: string, date: string) =>
    api.get<StockDailySnapshot>(`/api/stocks/daily/${code}/${date}`)
}
