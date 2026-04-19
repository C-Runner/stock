<template>
  <div class="analysis">
    <n-space vertical :size="20">
      <!-- Header -->
      <div class="stock-header">
        <n-button @click="goBack" quaternary>
          <template #icon>
            <n-icon><ArrowBack /></n-icon>
          </template>
        </n-button>
        <div class="stock-info" v-if="quote">
          <h2 class="stock-name">{{ quote.name }}</h2>
          <span class="stock-code">{{ stockCode }}</span>
        </div>
        <div class="stock-price" v-if="quote">
          <span class="price">¥{{ quote.current.toFixed(2) }}</span>
          <span :class="['change', quote.current > quote.open ? 'up' : 'down']">
            {{ quote.current >= quote.open ? '+' : '' }}{{ (((quote.current - quote.open) / quote.open) * 100).toFixed(2) }}%
          </span>
        </div>
      </div>

      <!-- Position Analysis -->
      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><Wallet /></n-icon>
            <span>Position Analysis</span>
          </div>
        </template>
        <n-spin :show="loading">
          <n-alert v-if="error" type="error">{{ error }}</n-alert>
          <template v-if="analysis">
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
        </n-spin>
      </n-card>

      <!-- Technical Analysis -->
      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><TrendCharts /></n-icon>
            <span>Technical Analysis</span>
          </div>
        </template>
        <n-spin :show="techLoading">
          <n-alert v-if="techError" type="warning">{{ techError }}</n-alert>
          <template v-if="technical">
            <!-- MA/EMA Section -->
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

            <!-- RSI Section -->
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

            <!-- MACD Section -->
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

            <!-- KDJ Section -->
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

            <!-- BOLL Section -->
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
        </n-spin>
      </n-card>

      <!-- Real-time Quote -->
      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><Clock /></n-icon>
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

      <!-- K-Line Chart -->
      <n-card class="analysis-card" :bordered="false" shadow>
        <template #header>
          <div class="card-header">
            <n-icon size="20"><ChartIcon /></n-icon>
            <span>K-Line Chart</span>
          </div>
        </template>
        <n-spin :show="techLoading">
          <n-alert v-if="techError" type="warning">{{ techError }}</n-alert>
          <KLineChart
            v-if="technical"
            :price-data="technical.recentPrices"
            :ma-data="technical.ma"
            height="450px"
          />
        </n-spin>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton, NCard, NSpace, NSpin, NAlert,
  NDivider, NIcon
} from 'naive-ui'
import { stockApi, type StockAnalysis, type StockQuote, type TechnicalAnalysis } from '../api'
import KLineChart from '../components/KLineChart.vue'

// Icons as functional components
const ArrowBack = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z' })
])
const Wallet = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M21 18v1c0 1.1-.9 2-2 2H5c-1.11 0-2-.9-2-2V5c0-1.1.89-2 2-2h14c1.1 0 2 .9 2 2v1h-9c-1.11 0-2 .9-2 2v8c0 1.1.89 2 2 2h9zm-9-2h10V8H12v8zm4-2.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z' })
])
const TrendCharts = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z' })
])
const Clock = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z' })
])
const ChartIcon = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z' })
])

const route = useRoute()
const router = useRouter()

const stockCode = ref(route.params.code as string)
const loading = ref(false)
const error = ref('')
const analysis = ref<StockAnalysis | null>(null)
const quote = ref<StockQuote | null>(null)
const technical = ref<TechnicalAnalysis | null>(null)
const techLoading = ref(false)
const techError = ref('')

const goBack = () => router.push('/')

const formatVolume = (vol: number): string => {
  if (vol >= 100000000) return (vol / 100000000).toFixed(2) + ' 亿'
  if (vol >= 10000) return (vol / 10000).toFixed(2) + ' 万'
  return vol.toString()
}

const formatAmount = (amt: number): string => {
  if (amt >= 100000000) return (amt / 100000000).toFixed(2) + ' 亿'
  if (amt >= 10000) return (amt / 10000).toFixed(2) + ' 万'
  return amt.toFixed(2)
}

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
    analysis.value = await stockApi.getAnalysis(stockCode.value)
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
  max-width: 1200px;
  margin: 0 auto;
  background: #000;
  min-height: calc(100vh - 60px);
  box-sizing: border-box;
}

/* Header */
.stock-header {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 24px;
  border-radius: 12px;
  color: white;
}

.stock-info {
  flex: 1;
}

.stock-name {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.stock-code {
  font-size: 14px;
  opacity: 0.8;
}

.stock-price {
  text-align: right;
}

.stock-price .price {
  display: block;
  font-size: 28px;
  font-weight: bold;
}

.stock-price .change {
  display: inline-block;
  font-size: 14px;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(255,255,255,0.2);
}

/* Cards */
.analysis-card {
  border-radius: 12px;
  background: #1a1a1a !important;
}

.analysis-card :deep(.n-card) {
  background: #1a1a1a !important;
  border: 1px solid #333 !important;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 16px;
  color: #fff;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.stat-value {
  font-size: 16px;
  font-weight: 500;
  color: #fff;
}

.stat-value.highlight {
  color: #667eea;
  font-weight: 600;
}

/* Profit Section */
.profit-section {
  display: flex;
  justify-content: space-around;
}

.profit-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.profit-label {
  font-size: 14px;
  color: #999;
}

.profit-value {
  font-size: 24px;
  font-weight: bold;
}

/* Tech Section */
.tech-section {
  margin-bottom: 8px;
}

.tech-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #999;
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
  padding: 12px;
  background: #1a1a1a;
  border-radius: 8px;
}

.tech-label {
  font-size: 12px;
  color: #888;
  margin-bottom: 4px;
}

.tech-value {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
}

/* Quote Grid */
.quote-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.quote-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  background: #1a1a1a;
  border-radius: 8px;
}

.quote-label {
  font-size: 12px;
  color: #888;
  margin-bottom: 4px;
}

.quote-value {
  font-size: 14px;
  font-weight: 500;
}

.quote-value.time {
  font-size: 12px;
}

/* Color Classes */
.up {
  color: #e74c3c;
}

.down {
  color: #27ae60;
}

.overbought {
  color: #e74c3c;
  font-weight: bold;
}

.oversold {
  color: #27ae60;
  font-weight: bold;
}

@media (max-width: 768px) {
  .analysis {
    padding: 12px;
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
  }

  .stock-price .price {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .analysis {
    padding: 8px;
    padding-bottom: calc(8px + env(safe-area-inset-bottom));
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
</style>
