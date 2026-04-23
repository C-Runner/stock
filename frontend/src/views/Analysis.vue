<template>
  <div class="analysis">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
    </div>

    <n-space vertical :size="12" class="content">
      <div class="stock-header">
        <div class="stock-info" v-if="quote">
          <div class="logo">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 3v18h18"/>
              <path d="M18 9l-5 5-4-4-3 3"/>
            </svg>
          </div>
          <div class="stock-text">
            <h2 class="stock-name">{{ quote.name }}</h2>
            <span class="stock-code">{{ stockCode }}</span>
          </div>
        </div>
        <div class="stock-price" v-if="quote">
          <span class="price">¥{{ quote.current.toFixed(2) }}</span>
          <span :class="['change', quote.current > quote.prevClose ? 'up' : 'down']">
            {{ quote.current >= quote.prevClose ? '+' : '' }}{{ (((quote.current - quote.prevClose) / quote.prevClose) * 100).toFixed(2) }}%
          </span>
        </div>
      </div>

      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><IconWallet /></n-icon>
            <span>Position Analysis</span>
          </div>
        </template>
        <n-spin v-if="loading" show description="Loading..." />
        <template v-else>
          <n-alert v-if="error" type="error">{{ error }}</n-alert>
          <template v-else-if="analysis">
            <div class="stats-grid">
              <div class="stat-item">
                <span class="stat-label">Holdings</span>
                <span class="stat-value">{{ analysis.quantity }} shares</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Avg Cost</span>
                <span class="stat-value">¥{{ analysis.buyPrice.toFixed(2) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Total Cost</span>
                <span class="stat-value">¥{{ analysis.cost.toFixed(2) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Market Value</span>
                <span class="stat-value highlight">¥{{ analysis.marketValue.toFixed(2) }}</span>
              </div>
            </div>
            <n-divider />
            <div class="profit-section">
              <div class="profit-item">
                <span class="profit-label">Profit/Loss</span>
                <span :class="['profit-value', analysis.profitLoss >= 0 ? 'up' : 'down']">
                  {{ analysis.profitLoss >= 0 ? '+' : '' }}¥{{ analysis.profitLoss.toFixed(2) }}
                </span>
              </div>
              <div class="profit-item">
                <span class="profit-label">Return Rate</span>
                <span :class="['profit-value', analysis.profitRate >= 0 ? 'up' : 'down']">
                  {{ analysis.profitRate >= 0 ? '+' : '' }}{{ analysis.profitRate.toFixed(2) }}%
                </span>
              </div>
            </div>
          </template>
          <template v-else-if="!hasPosition">
            <div class="no-position">
              <div class="no-position-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <circle cx="12" cy="12" r="9"/>
                  <path d="M12 8v4"/>
                  <circle cx="12" cy="16" r="1"/>
                </svg>
              </div>
              <div class="no-position-text">
                <span class="no-position-title">No Position Held</span>
                <span class="no-position-desc">You don't own any shares of this stock yet</span>
              </div>
            </div>
          </template>
        </template>
      </n-card>

      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><IconTrend /></n-icon>
            <span>Technical Analysis</span>
          </div>
        </template>
        <n-spin v-if="techLoading" show description="Loading..." />
        <template v-else>
          <n-alert v-if="techError" type="warning">{{ techError }}</n-alert>
          <template v-else-if="technical">
            <div class="tech-section">
              <h4 class="tech-title">Moving Averages</h4>
              <div class="tech-grid">
                <div class="tech-item" v-for="ma in technical.ma" :key="'ma'+ma.period">
                  <span class="tech-label">MA{{ ma.period }}</span>
                  <span class="tech-value">¥{{ ma.value.toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <n-divider />

            <div class="tech-section">
              <h4 class="tech-title">RSI</h4>
              <div class="tech-grid">
                <div class="tech-item" v-for="rsi in technical.rsi" :key="'rsi'+rsi.period">
                  <span class="tech-label">RSI({{ rsi.period }})</span>
                  <span :class="['tech-value', getRSIClass(rsi.value)]">{{ rsi.value.toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <n-divider />

            <div class="tech-section">
              <h4 class="tech-title">MACD</h4>
              <div class="tech-grid three-col">
                <div class="tech-item">
                  <span class="tech-label">DIF</span>
                  <span class="tech-value">{{ technical.macd.dif.toFixed(4) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">DEA</span>
                  <span class="tech-value">{{ technical.macd.dea.toFixed(4) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">MACD</span>
                  <span :class="['tech-value', technical.macd.macd >= 0 ? 'up' : 'down']">
                    {{ technical.macd.macd >= 0 ? '+' : '' }}{{ technical.macd.macd.toFixed(4) }}
                  </span>
                </div>
              </div>
            </div>

            <n-divider />

            <div class="tech-section">
              <h4 class="tech-title">KDJ</h4>
              <div class="tech-grid three-col">
                <div class="tech-item">
                  <span class="tech-label">K</span>
                  <span :class="['tech-value', getKDJClass(technical.kdj.k)]">{{ technical.kdj.k.toFixed(2) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">D</span>
                  <span :class="['tech-value', getKDJClass(technical.kdj.d)]">{{ technical.kdj.d.toFixed(2) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">J</span>
                  <span :class="['tech-value', getKDJClass(technical.kdj.j)]">{{ technical.kdj.j.toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <n-divider />

            <div class="tech-section">
              <h4 class="tech-title">Bollinger Bands</h4>
              <div class="tech-grid three-col">
                <div class="tech-item">
                  <span class="tech-label">Upper</span>
                  <span class="tech-value up">¥{{ technical.boll.upper.toFixed(2) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">Middle</span>
                  <span class="tech-value">¥{{ technical.boll.mid.toFixed(2) }}</span>
                </div>
                <div class="tech-item">
                  <span class="tech-label">Lower</span>
                  <span class="tech-value down">¥{{ technical.boll.lower.toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </template>
        </template>
      </n-card>

      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><IconClock /></n-icon>
            <span>Real-time Quote</span>
          </div>
        </template>
        <template v-if="quote">
          <div class="quote-grid">
            <div class="quote-item">
              <span class="quote-label">Open</span>
              <span class="quote-value">¥{{ quote.open.toFixed(2) }}</span>
            </div>
            <div class="quote-item">
              <span class="quote-label">High</span>
              <span class="quote-value up">¥{{ quote.high.toFixed(2) }}</span>
            </div>
            <div class="quote-item">
              <span class="quote-label">Low</span>
              <span class="quote-value down">¥{{ quote.low.toFixed(2) }}</span>
            </div>
            <div class="quote-item">
              <span class="quote-label">Volume</span>
              <span class="quote-value">{{ formatVolume(quote.volume) }}</span>
            </div>
            <div class="quote-item">
              <span class="quote-label">Amount</span>
              <span class="quote-value">¥{{ formatAmount(quote.amount) }}</span>
            </div>
            <div class="quote-item">
              <span class="quote-label">Update</span>
              <span class="quote-value time">{{ quote.updateTime }}</span>
            </div>
          </div>
        </template>
      </n-card>

      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><IconChart /></n-icon>
            <span>K-Line Chart</span>
          </div>
        </template>
        <n-spin v-if="techLoading" show description="Loading..." />
        <template v-else>
          <n-alert v-if="techError" type="warning">{{ techError }}</n-alert>
          <KLineChart
            v-if="technical"
            :price-data="technical.recentPrices"
            :ma-data="technical.ma"
            height="350px"
          />
        </template>
      </n-card>
    </n-space>

    <div class="bottom-tabs">
      <div class="tab-item" @click="router.push('/')">
        <div class="tab-icon">
          <n-icon size="22"><IconHome /></n-icon>
        </div>
        <span class="tab-label">Portfolio</span>
      </div>
      <div class="tab-item primary" @click="router.push('/watchlist')">
        <div class="tab-icon">
          <n-icon size="22"><IconStar /></n-icon>
        </div>
        <span class="tab-label">Watchlist</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard, NSpace, NSpin, NAlert,
  NDivider, NIcon
} from 'naive-ui'
import { stockApi, type StockAnalysis, type StockQuote, type TechnicalAnalysis } from '../api'
import { IconWallet, IconTrend, IconClock, IconChart, IconHome, IconStar } from '../components/icons'
import { formatVolume, formatAmount } from '../utils/format'
import KLineChart from '../components/KLineChart.vue'

const route = useRoute()
const router = useRouter()

const stockCode = ref(route.params.code as string)
const loading = ref(false)
const error = ref('')
const hasPosition = ref(true)
const analysis = ref<StockAnalysis | null>(null)
const quote = ref<StockQuote | null>(null)
const technical = ref<TechnicalAnalysis | null>(null)
const techLoading = ref(false)
const techError = ref('')

const getRSIClass = (rsi: number): string => {
  if (rsi > 70) return 'overbought'
  if (rsi < 30) return 'oversold'
  return ''
}

const getKDJClass = (value: number): string => {
  if (value > 80) return 'overbought'
  if (value < 20) return 'oversold'
  return ''
}

const fetchAnalysis = async () => {
  loading.value = true
  error.value = ''
  try {
    const result = await stockApi.getAnalysis(stockCode.value)
    if ('hasPosition' in result && !result.hasPosition) {
      hasPosition.value = false
      analysis.value = null
    } else {
      hasPosition.value = true
      analysis.value = result as StockAnalysis
    }
  } catch (err) {
    error.value = (err as Error).message
  } finally {
    loading.value = false
  }
}

const fetchQuote = async () => {
  try {
    quote.value = await stockApi.getQuote(stockCode.value)
  } catch (err) {
    console.error(err)
  }
}

const fetchTechnical = async () => {
  techLoading.value = true
  techError.value = ''
  try {
    technical.value = await stockApi.getTechnical(stockCode.value)
  } catch (err) {
    techError.value = (err as Error).message
  } finally {
    techLoading.value = false
  }
}

onMounted(() => {
  fetchAnalysis()
  fetchQuote()
  fetchTechnical()
})
</script>

<style scoped>
.analysis {
  width: 100%;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px;
  padding-bottom: calc(80px + env(safe-area-inset-bottom));
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

.background {
  position: fixed;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: #6366f1;
  top: -150px;
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: #8b5cf6;
  bottom: -100px;
  right: -100px;
  animation: float 10s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}

.content {
  position: relative;
  align-items: stretch;
}

.content > * {
  flex-shrink: 0;
}

.stock-header {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(139, 92, 246, 0.2));
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 24px;
  border-radius: 20px;
  backdrop-filter: blur(20px);
  color: white;
  transition: all 0.3s ease;
}

.stock-header:hover {
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 0 8px 32px rgba(99, 102, 241, 0.15);
}

.back-btn {
  color: rgba(255, 255, 255, 0.8) !important;
}

.stock-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 44px;
  height: 44px;
  padding: 8px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo svg {
  width: 100%;
  height: 100%;
  color: #fff;
}

.stock-text {
  display: flex;
  flex-direction: column;
}

.stock-name {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  color: #fff;
}

.stock-code {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.stock-price {
  text-align: right;
}

.stock-price .price {
  display: block;
  font-size: 26px;
  font-weight: bold;
  color: #fff;
}

.stock-price .change {
  display: inline-block;
  font-size: 13px;
  padding: 6px 12px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.08);
  font-weight: 500;
  backdrop-filter: blur(4px);
}

.analysis-card {
  --n-border-radius: 20px !important;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
  overflow: hidden !important;
}

.analysis-card:last-of-type {
  margin-bottom: calc(80px + env(safe-area-inset-bottom)) !important;
}

.analysis-card :deep(.n-card) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  border-radius: 20px !important;
  overflow: hidden !important;
  padding: 0 !important;
  margin: 0 !important;
  --n-padding-top: 0 !important;
  --n-padding-bottom: 0 !important;
  --n-padding-left: 0 !important;
  --n-padding-right: 0 !important;
}

.analysis-card :deep(.n-card__content),
.analysis-card :deep(.n-card-body),
.analysis-card :deep(.n-card-content) {
  padding: 16px !important;
  margin: 0 !important;
  border-radius: inherit !important;
}

.analysis-card :deep(.n-card__footer),
.analysis-card :deep(.n-card-footer) {
  display: none !important;
  padding: 0 !important;
  margin: 0 !important;
  min-height: 0 !important;
  height: 0 !important;
  border: none !important;
}

.analysis-card :deep(.n-base-segments),
.analysis-card :deep(.n-card__border) {
  display: none !important;
}

.analysis-card :deep(.n-card-header) {
  background: transparent !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06) !important;
  padding: 16px 20px !important;
  margin: 0 !important;
  border-radius: inherit !important;
}

.analysis-card :deep(.n-card-wrapper) {
  margin: 0 !important;
  padding: 0 !important;
  border-radius: inherit !important;
}

.analysis-card :deep(.n-base-card) {
  border-radius: inherit !important;
  overflow: hidden !important;
}

.analysis-card :deep(.n-card),
.analysis-card :deep(.n-base-segment) {
  border-radius: 20px !important;
  overflow: hidden !important;
}

.analysis-card:hover {
  border-color: rgba(255, 255, 255, 0.15) !important;
  box-shadow: 0 8px 32px rgba(99, 102, 241, 0.15) !important;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 16px;
  color: #fff;
}

.card-header :deep(.n-icon) {
  color: #6366f1 !important;
}

.no-position {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px 20px;
  text-align: left;
}

.no-position-icon {
  width: 36px;
  height: 36px;
  padding: 8px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 50%;
  color: #6366f1;
  flex-shrink: 0;
}

.no-position-icon svg {
  width: 100%;
  height: 100%;
}

.no-position-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.no-position-title {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.no-position-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  transition: all 0.2s ease;
}

.stat-item:hover {
  background: rgba(99, 102, 241, 0.06);
  border-color: rgba(99, 102, 241, 0.15);
}

.stat-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}

.stat-value {
  font-size: 17px;
  font-weight: 600;
  color: #fff;
}

.stat-value.highlight {
  color: #6366f1;
  font-weight: 700;
}

.profit-section {
  display: flex;
  justify-content: space-around;
  padding: 16px 0;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 14px;
  margin-top: 8px;
}

.profit-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 8px 24px;
}

.profit-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.profit-value {
  font-size: 26px;
  font-weight: bold;
}

.tech-section {
  margin-bottom: 8px;
}

.tech-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 500;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.tech-grid.three-col {
  grid-template-columns: repeat(3, 1fr);
}

.tech-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 14px 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  transition: all 0.2s ease;
}

.tech-item:hover {
  background: rgba(99, 102, 241, 0.08);
  border-color: rgba(99, 102, 241, 0.2);
  transform: translateY(-2px);
}

.tech-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}

.tech-value {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.quote-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.quote-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 14px 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  transition: all 0.2s ease;
}

.quote-item:hover {
  background: rgba(99, 102, 241, 0.08);
  border-color: rgba(99, 102, 241, 0.2);
}

.quote-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}

.quote-value {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.quote-value.time {
  font-size: 12px;
}

.up {
  color: #ff6b6b;
}

.down {
  color: #38ef7d;
}

.overbought {
  color: #ff6b6b;
  font-weight: bold;
}

.oversold {
  color: #38ef7d;
  font-weight: bold;
}

@media (max-width: 768px) {
  .analysis {
    padding: 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .tech-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .quote-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stock-header {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }

  .stock-info, .stock-price {
    text-align: center;
    justify-content: center;
  }

  .stock-price .price {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .analysis {
    padding: 16px 12px;
    padding-bottom: calc(16px + env(safe-area-inset-bottom));
  }

  .stock-header {
    padding: 16px 12px;
  }

  .stock-name {
    font-size: 18px;
  }

  .stock-code {
    font-size: 12px;
  }

  .stock-price .price {
    font-size: 22px;
  }

  .stock-price .change {
    font-size: 12px;
    padding: 2px 6px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  .stat-item {
    padding: 8px;
  }

  .stat-label {
    font-size: 10px;
  }

  .stat-value {
    font-size: 12px;
  }

  .profit-section {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 16px;
    justify-content: center;
  }

  .profit-item {
    min-width: 100px;
  }

  .profit-value {
    font-size: 18px;
  }

  .tech-section {
    margin-bottom: 4px;
  }

  .tech-title {
    font-size: 12px;
    margin-bottom: 8px;
  }

  .tech-grid, .tech-grid.three-col {
    grid-template-columns: repeat(3, 1fr);
    gap: 6px;
  }

  .tech-item {
    padding: 8px 4px;
  }

  .tech-label {
    font-size: 10px;
  }

  .tech-value {
    font-size: 12px;
  }

  .quote-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  .quote-item {
    padding: 8px;
  }

  .quote-label {
    font-size: 10px;
  }

  .quote-value {
    font-size: 12px;
  }

  .card-header {
    font-size: 14px;
  }

  :deep(.n-card__content) {
    padding: 12px 8px;
  }

:deep(.n-divider) {
    margin: 12px 0;
  }
}

.bottom-tabs {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 70px;
  background: rgba(20, 19, 60, 0.4);
  backdrop-filter: blur(32px) saturate(180%);
  -webkit-backdrop-filter: blur(32px) saturate(180%);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.15);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  padding: 0 24px;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
  touch-action: none;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 10px 20px;
  border-radius: 16px;
  transition: all 0.2s ease;
  color: rgba(255, 255, 255, 0.6);
}

.tab-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
}

.tab-item.primary {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.4);
}

.tab-item.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.5);
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-label {
  font-size: 12px;
  font-weight: 600;
}
</style>
