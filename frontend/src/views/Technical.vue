<template>
  <div class="technical">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
    </div>

    <n-space vertical :size="12" class="content" :style="{ width: '100%' }">
      <div class="page-header">
        <div class="back-btn" @click="router.back()">
          <n-icon size="24"><IconArrowBack /></n-icon>
        </div>
        <div class="header-info">
          <h1 class="page-title">Technical Analysis</h1>
        </div>
      </div>

      <n-spin v-if="loading" show description="Loading..." />

      <template v-else-if="technical">
        <n-card class="analysis-card recommendation-card" :class="technical.recommendation.action">
          <div class="recommendation-top">
            <div class="stock-info-section">
              <span class="stock-name" v-if="quote">{{ quote.name }}</span>
              <span class="stock-code">({{ stockCode }})</span>
              <span class="stock-price" v-if="quote">¥{{ quote.current.toFixed(2) }}</span>
            </div>
            <div class="action-section">
              <div class="action-badge-large" :class="technical.recommendation.action">
                <span class="action-icon" v-if="technical.recommendation.action === 'buy'">↑</span>
                <span class="action-icon" v-else-if="technical.recommendation.action === 'watch'">↓</span>
                <span class="action-icon" v-else>→</span>
                {{ technical.recommendation.action.toUpperCase() }}
              </div>
              <div class="confidence-ring">
                <svg viewBox="0 0 36 36" class="confidence-ring-svg">
                  <path class="confidence-ring-bg" d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"/>
                  <path class="confidence-ring-fill" :class="technical.recommendation.action" :stroke-dasharray="`${technical.recommendation.confidence}, 100`" d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"/>
                </svg>
                <span class="confidence-value">{{ technical.recommendation.confidence.toFixed(0) }}%</span>
              </div>
            </div>
          </div>
          <div class="risk-indicator" :class="technical.recommendation.riskLevel">
            <span class="risk-label">Risk:</span>
            <span class="risk-level">{{ technical.recommendation.riskLevel.toUpperCase() }}</span>
            <div class="risk-bar">
              <div class="risk-bar-fill" :class="technical.recommendation.riskLevel"></div>
            </div>
          </div>
          <p class="recommendation-summary">{{ technical.recommendation.summary }}</p>
        </n-card>

        <!-- Technical Analysis by Period + Reasons -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconTrend /></n-icon>
              <span>Technical Analysis by Period</span>
            </div>
          </template>

          <div class="period-vertical">
            <!-- Short-term -->
            <div class="period-row short">
              <div class="period-row-header">
                <span class="period-label">Short-term</span>
                <span class="period-badge short">1-10 Days</span>
              </div>
              <p class="period-desc">{{ getShortTermSummary() }}</p>
            </div>

            <!-- Medium-term -->
            <div class="period-row medium">
              <div class="period-row-header">
                <span class="period-label">Medium-term</span>
                <span class="period-badge medium">10-30 Days</span>
              </div>
              <p class="period-desc">{{ getMediumTermSummary() }}</p>
            </div>

            <!-- Long-term -->
            <div class="period-row long">
              <div class="period-row-header">
                <span class="period-label">Long-term</span>
                <span class="period-badge long">30-180 Days</span>
              </div>
              <p class="period-desc">{{ getLongTermSummary() }}</p>
            </div>
          </div>
        </n-card>

        <!-- Support & Resistance -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconTrend /></n-icon>
              <span>Support & Resistance</span>
            </div>
          </template>
          <p class="period-desc">{{ getSupportResistanceSummary() }}</p>
        </n-card>

        <!-- Pattern Recognition -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconChart /></n-icon>
              <span>Pattern Recognition</span>
            </div>
          </template>
          <p class="period-desc">{{ getPatternSummary() }}</p>
        </n-card>

        <!-- Support & Resistance -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconChart /></n-icon>
              <span>K-Line Chart</span>
            </div>
          </template>
          <KLineChart
            v-if="technical"
            :price-data="technical.recentPrices"
            :ma-data="technical.ma"
            :rsi-data="technical.rsi"
            :kdj-data="technical.kdj"
            height="550px"
          />
        </n-card>
      </template>

      <n-alert v-else-if="error" type="error">{{ error }}</n-alert>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard, NSpace, NSpin, NAlert, NIcon
} from 'naive-ui'
import { stockApi, type TechnicalAnalysis, type StockQuote } from '../api'
import { IconTrend, IconChart, IconArrowBack } from '../components/icons'
import KLineChart from '../components/KLineChart.vue'

const route = useRoute()
const router = useRouter()

const stockCode = ref(route.params.code as string)
const loading = ref(false)
const error = ref('')
const technical = ref<TechnicalAnalysis | null>(null)
const quote = ref<StockQuote | null>(null)

const getLatestValue = (values: number[]): number => {
  if (!values || values.length === 0) return 0
  return values[values.length - 1] ?? 0
}

// Short-term: MA5, RSI6, KDJ, WR
const getShortTermMA = () => technical.value?.ma.filter(m => m.period === 5) || []
const getShortTermRSI = () => technical.value?.rsi.filter(r => r.period === 6) || []

// Medium-term: MA10/20, EMA12, RSI12
const getMediumTermMA = () => technical.value?.ma.filter(m => m.period === 10 || m.period === 20) || []
const getMediumTermEMA = () => technical.value?.ema.filter(e => e.period === 12) || []
const getMediumTermRSI = () => technical.value?.rsi.filter(r => r.period === 12) || []

// Long-term: MA60, EMA26, RSI24
const getLongTermMA = () => technical.value?.ma.filter(m => m.period === 60) || []
const getLongTermEMA = () => technical.value?.ema.filter(e => e.period === 26) || []
const getLongTermRSI = () => technical.value?.rsi.filter(r => r.period === 24) || []

// Summary text for each period
const getShortTermSummary = () => {
  const ma5 = getLatestValue(getShortTermMA()[0]?.values || [])
  const rsi6 = getLatestValue(getShortTermRSI()[0]?.values || [])
  const k = getLatestValue(technical.value?.kdj.k || [])
  const d = getLatestValue(technical.value?.kdj.d || [])
  const dif = getLatestValue(technical.value?.macd.dif || [])
  const dea = getLatestValue(technical.value?.macd.dea || [])
  const prices = technical.value?.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const wr = getLatestValue(technical.value?.wr.find(w => w.period === 10)?.values || [])

  const signals: string[] = []

  if (rsi6 < 30) signals.push('RSI 超卖区域，短期超跌反弹概率增加')
  else if (rsi6 > 70) signals.push('RSI 超买区域，短期注意回调风险')
  else if (rsi6 < 40) signals.push('RSI 偏弱，卖方力量占优')
  else if (rsi6 > 60) signals.push('RSI 偏强，买方力量占优')

  if (k < 20 && d < 20) signals.push('KDJ 超卖区域，低位金叉信号')
  else if (k > 80 && d > 80) signals.push('KDJ 超买区域，高位死叉风险')
  else if (k < d && k < 30) signals.push('KDJ 低位钝化，跌势动能减弱')
  else if (k > d) signals.push('KDJ 形成金叉，短期看涨')
  else signals.push('KDJ 中性区域，等待方向确认')

  if (dif > dea) signals.push('MACD 在零轴上方形成金叉，上涨动能持续')
  else if (dif < 0 && dea < 0) signals.push('MACD 在零轴下方运行，空头趋势')
  else if (dif > dea) signals.push('MACD 温和金叉，短线看涨')
  else signals.push('MACD 死叉形成，短期调整信号')

  if (currentPrice > ma5) signals.push('价格站稳 MA5，短线支撑有效')
  else signals.push('价格跌破 MA5，短线压力较大')

  if (wr < -80) signals.push('WR 接近极值超卖区域')
  else if (wr > -20) signals.push('WR 显示超买区域')

  if (technical.value?.obv.trend === 'rising') signals.push('OBV 上升，成交量配合上涨')
  else if (technical.value?.obv.trend === 'falling') signals.push('OBV 下降，量价背离风险')

  return signals.length > 0 ? signals.join('；') : '短期趋势不明，建议观望'
}

const getMediumTermSummary = () => {
  const ma10 = getLatestValue(getMediumTermMA()[0]?.values || [])
  const ma20 = getLatestValue(getMediumTermMA()[1]?.values || [])
  const ema12 = getLatestValue(getMediumTermEMA()[0]?.values || [])
  const rsi12 = getLatestValue(getMediumTermRSI()[0]?.values || [])
  const upper = getLatestValue(technical.value?.boll.upper || [])
  const lower = getLatestValue(technical.value?.boll.lower || [])
  const mid = getLatestValue(technical.value?.boll.mid || [])
  const prices = technical.value?.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const vwap = getLatestValue(technical.value?.vwap.values || [])

  const signals: string[] = []

  if (ma10 > ma20) signals.push('MA10 > MA20 多头排列，中期上升趋势')
  else if (ma10 < ma20) signals.push('MA10 < MA20 空头排列，中期下降趋势')
  else signals.push('MA10 与 MA20 粘合，中期横盘整理')

  if (currentPrice > ma20) signals.push('价格站稳 MA20，中期支撑有效')
  else signals.push('价格跌破 MA20，中期压力较大')

  if (currentPrice > ema12) signals.push('价格站上 EMA12，均线呈多头')
  else signals.push('价格跌破 EMA12，均线转空')

  if (rsi12 < 30) signals.push('RSI12 超卖，趋势可能反转')
  else if (rsi12 > 70) signals.push('RSI12 超买，注意中期回调风险')
  else if (rsi12 > 55) signals.push('RSI12 偏强，中期趋势向上')
  else if (rsi12 < 45) signals.push('RSI12 偏弱，中期趋势向下')

  if (currentPrice > upper) {
    const bandwidth = (upper - lower) / mid
    if (bandwidth > 0.1) signals.push('股价突破 BOLL 上轨，动能强劲但警惕回调')
    else signals.push('股价触及 BOLL 上轨，注意压力')
  } else if (currentPrice < lower) {
    signals.push('股价跌破 BOLL 下轨，超卖信号，可能反弹')
  } else {
    signals.push('股价在 BOLL 中轨运行，中性震荡')
  }

  if (currentPrice > vwap) signals.push('价格高于 VWAP，日内做多信号')
  else signals.push('价格低于 VWAP，日内做空信号')

  return signals.length > 0 ? signals.join('；') : '中期趋势不明，建议观望'
}

const getLongTermSummary = () => {
  const ma60 = getLatestValue(getLongTermMA()[0]?.values || [])
  const ema26 = getLatestValue(getLongTermEMA()[0]?.values || [])
  const rsi24 = getLatestValue(getLongTermRSI()[0]?.values || [])
  const prices = technical.value?.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const plusDI = technical.value?.dmi.plusDI || 0
  const minusDI = technical.value?.dmi.minusDI || 0
  const adx = technical.value?.dmi.adx || 0

  const signals: string[] = []

  if (currentPrice > ma60) {
    signals.push('价格位于 MA60 上方，长期趋势向上')
  } else {
    signals.push('价格位于 MA60 下方，长期趋势向下')
  }

  if (currentPrice > ema26) signals.push('EMA26 提供有效支撑，多头格局')
  else signals.push('EMA26 成为压力位，空头格局')

  if (rsi24 < 30) signals.push('RSI24 严重超卖，长线入场机会')
  else if (rsi24 > 70) signals.push('RSI24 超买区域，长线注意风险')
  else if (rsi24 > 55) signals.push('RSI24 偏强，长线趋势向上')
  else if (rsi24 < 45) signals.push('RSI24 偏弱，长线趋势向下')
  else signals.push('RSI24 中性区域，长线等待方向')

  if (adx > 25) {
    if (plusDI > minusDI) signals.push('DMI +DI 上穿，长期上升趋势确认')
    else signals.push('DMI -DI 上穿，长期下降趋势确认')
  } else if (adx < 20) {
    signals.push('ADX 低于 20，长期趋势不明显，横盘整理')
  } else {
    signals.push('ADX 趋势信号不强，等待趋势形成')
  }

  const priceDiff = ((currentPrice - ma60) / ma60 * 100).toFixed(2)
  if (Math.abs(parseFloat(priceDiff)) < 2) signals.push(`价格偏离 MA60 仅 ${priceDiff}%，趋势待确认`)

  return signals.length > 0 ? signals.join('；') : '长期趋势不明，建议观望'
}

const getSupportResistanceSummary = () => {
  const levels = technical.value?.levels
  if (!levels) return 'No significant support or resistance levels detected.'

  const topResistanceItem = levels.resistance[0]
  const topSupportItem = levels.support[0]
  const topResistance = topResistanceItem ? `¥${topResistanceItem.price.toFixed(2)} (${topResistanceItem.strength})` : 'none'
  const topSupport = topSupportItem ? `¥${topSupportItem.price.toFixed(2)} (${topSupportItem.strength})` : 'none'
  const trendLine = levels.trendLine ? `${levels.trendLine.direction.toUpperCase()} trend (${levels.trendLine.strength})` : 'no clear trend line'

  return `Resistance: ${topResistance}\nSupport: ${topSupport}\nTrend: ${trendLine}`
}

const getPatternSummary = () => {
  const patterns = technical.value?.patterns
  if (!patterns || (patterns.candlestickPatterns.length === 0 && patterns.trendPatterns.length === 0)) {
    return 'No significant patterns detected.'
  }

  const recentCandles = patterns.candlestickPatterns.slice(-3)
  const recentTrends = patterns.trendPatterns.slice(-2)

  let summary = ''

  if (recentCandles.length > 0) {
    const latest = recentCandles[recentCandles.length - 1]
    if (latest) {
      summary += `Recent candle: ${latest.type.replace('_', ' ')} (${latest.bullish ? 'bullish' : 'bearish'}, ${latest.strength}). `
    }
  }

  if (recentTrends.length > 0) {
    const latest = recentTrends[recentTrends.length - 1]
    if (latest) {
      summary += `Chart pattern: ${latest.type.replace('_', ' ')} (${latest.bullish ? 'bullish' : 'bearish'}, ${latest.strength}).`
    }
  }

  return summary || 'No significant patterns detected.'
}

const fetchData = async () => {
  loading.value = true
  error.value = ''
  try {
    const [techData, quoteData] = await Promise.all([
      stockApi.getTechnical(stockCode.value),
      stockApi.getQuote(stockCode.value)
    ])
    technical.value = techData
    quote.value = quoteData
  } catch (err) {
    error.value = (err as Error).message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.technical {
  width: 100%;
  height: 100%;
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
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-y;
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
  background: #10b981;
  top: -150px;
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: #059669;
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
  width: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 8px;
}

.back-btn {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: rgba(255, 255, 255, 0.8);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateX(-2px);
}

.header-info {
  display: flex;
  flex-direction: column;
}

.page-title {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  color: #fff;
}

.stock-code {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.analysis-card {
  --n-border-radius: 20px !important;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
}

.analysis-card:hover {
  border-color: rgba(255, 255, 255, 0.15) !important;
  box-shadow: 0 8px 32px rgba(16, 185, 129, 0.1) !important;
}

.recommendation-card {
  padding: 24px;
}

.recommendation-card.buy {
  background: linear-gradient(145deg, rgba(16, 185, 129, 0.2), rgba(5, 150, 105, 0.15), rgba(16, 185, 129, 0.05)) !important;
  border: 1px solid rgba(16, 185, 129, 0.4) !important;
}

.recommendation-card.hold {
  background: linear-gradient(145deg, rgba(245, 158, 11, 0.2), rgba(217, 119, 6, 0.15), rgba(245, 158, 11, 0.05)) !important;
  border: 1px solid rgba(245, 158, 11, 0.4) !important;
}

.recommendation-card.watch {
  background: linear-gradient(145deg, rgba(239, 68, 68, 0.2), rgba(220, 38, 38, 0.15), rgba(239, 68, 68, 0.05)) !important;
  border: 1px solid rgba(239, 68, 68, 0.4) !important;
}

.recommendation-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.stock-info-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stock-info-section .stock-name {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
}

.stock-info-section .stock-code {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

.stock-info-section .stock-price {
  font-size: 18px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4px;
}

.action-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.action-badge-large {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 22px;
  font-weight: 800;
  padding: 12px 24px;
  border-radius: 14px;
  letter-spacing: 2px;
}

.action-badge-large.buy {
  background: linear-gradient(135deg, #10b981, #059669);
  color: #fff;
  box-shadow: 0 4px 20px rgba(16, 185, 129, 0.4);
}

.action-badge-large.hold {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
  box-shadow: 0 4px 20px rgba(245, 158, 11, 0.4);
}

.action-badge-large.watch {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: #fff;
  box-shadow: 0 4px 20px rgba(239, 68, 68, 0.4);
}

.action-icon {
  font-size: 20px;
}

.confidence-ring {
  position: relative;
  width: 56px;
  height: 56px;
}

.confidence-ring-svg {
  transform: rotate(-90deg);
  width: 100%;
  height: 100%;
}

.confidence-ring-bg {
  fill: none;
  stroke: rgba(255, 255, 255, 0.1);
  stroke-width: 3;
}

.confidence-ring-fill {
  fill: none;
  stroke-width: 3;
  stroke-linecap: round;
}

.confidence-ring-fill.buy { stroke: #10b981; }
.confidence-ring-fill.hold { stroke: #f59e0b; }
.confidence-ring-fill.watch { stroke: #ef4444; }

.confidence-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 12px;
  font-weight: 700;
  color: #fff;
}

.risk-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}

.risk-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  text-transform: uppercase;
}

.risk-level {
  font-size: 12px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 6px;
}

.risk-level.low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.risk-level.medium { background: rgba(245, 158, 11, 0.2); color: #f59e0b; }
.risk-level.high { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

.risk-bar {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  overflow: hidden;
}

.risk-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s ease;
}

.risk-bar-fill.low { width: 33%; background: #10b981; }
.risk-bar-fill.medium { width: 66%; background: #f59e0b; }
.risk-bar-fill.high { width: 100%; background: #ef4444; }

.recommendation-summary {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
  line-height: 1.6;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  border-left: 3px solid rgba(255, 255, 255, 0.2);
}

.recommendation-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.stock-info-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.stock-name {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.stock-info-header .stock-code {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.action-badge {
  font-size: 18px;
  font-weight: 700;
  padding: 8px 18px;
  border-radius: 10px;
  letter-spacing: 1px;
}

.action-badge.buy {
  background: linear-gradient(135deg, #10b981, #059669);
  color: #fff;
}

.action-badge.hold {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
}

.action-badge.watch {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: #fff;
}

.confidence-badge {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
}

.risk-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  text-transform: uppercase;
}

.risk-badge.low {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.risk-badge.medium {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.risk-badge.high {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.recommendation-summary {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  line-height: 1.5;
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
  color: #10b981 !important;
}

.period-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 16px;
  color: #fff;
}

.period-header.short :deep(.n-icon) {
  color: #f59e0b !important;
}

.period-header.medium :deep(.n-icon) {
  color: #10b981 !important;
}

.period-header.long :deep(.n-icon) {
  color: #3b82f6 !important;
}

.period-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-left: auto;
}

.period-badge.short {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.period-badge.medium {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.period-badge.long {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

/* Period Vertical Layout */
.period-vertical {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.period-row {
  padding: 14px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.period-row-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.period-label {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.period-row.short .period-label {
  color: #f59e0b;
}

.period-row.medium .period-label {
  color: #10b981;
}

.period-row.long .period-label {
  color: #3b82f6;
}

.indicator-text {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.period-desc {
  margin: 0;
  font-size: 12px;
  line-height: 1.9;
  color: rgba(255, 255, 255, 0.8);
  white-space: pre-line;
}

.tech-section {
  margin-bottom: 12px;
}

.tech-title {
  margin: 0 0 14px 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tech-title::before {
  content: '';
  width: 4px;
  height: 14px;
  background: linear-gradient(180deg, #10b981, #059669);
  border-radius: 2px;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

.tech-grid.three-col {
  grid-template-columns: repeat(3, 1fr);
}

.tech-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.tech-item.full {
  flex-direction: row;
  justify-content: space-between;
}

.tech-label {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}

.tech-value {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.up { color: #ff6b6b; }
.down { color: #38ef7d; }
.overbought { color: #ff6b6b; font-weight: bold; }
.oversold { color: #38ef7d; font-weight: bold; }

.volume-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
}

.obv-trend {
  font-size: 12px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 6px;
}

.obv-trend.rising {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.obv-trend.falling {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.obv-trend.neutral {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.5);
}

.obv-meaning {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.levels-container {
  display: flex;
  gap: 16px;
}

.level-group {
  flex: 1;
}

.level-header {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 6px 12px;
  border-radius: 8px;
  text-align: center;
  margin-bottom: 8px;
}

.level-header.resistance {
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
}

.level-header.support {
  color: #38ef7d;
  background: rgba(56, 239, 125, 0.1);
}

.level-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.level-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.level-price {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  flex: 1;
}

.level-touches {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
}

.level-strength {
  font-size: 9px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.level-strength.strong {
  color: #10b981;
  background: rgba(16, 185, 129, 0.15);
}

.level-strength.medium {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.15);
}

.level-strength.weak {
  color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.05);
}

.no-levels {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
  padding: 12px;
}

.trend-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.trend-direction {
  font-size: 12px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 6px;
}

.trend-direction.up {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.trend-direction.down {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.trend-strength {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.patterns-section {
  margin-bottom: 16px;
}

.patterns-section:last-child {
  margin-bottom: 0;
}

.patterns-title {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 10px 0;
}

.pattern-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.pattern-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 10px;
  flex-wrap: wrap;
}

.pattern-badge {
  font-size: 9px;
  font-weight: 700;
  padding: 4px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.pattern-badge.bullish {
  background: rgba(56, 239, 125, 0.15);
  color: #38ef7d;
}

.pattern-badge.bearish {
  background: rgba(255, 107, 107, 0.15);
  color: #ff6b6b;
}

.pattern-date {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
}

.pattern-strength {
  font-size: 9px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.5);
}

.pattern-meaning {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  width: 100%;
  margin-top: 4px;
}

.no-patterns {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
  padding: 20px;
}

@media (max-width: 768px) {
  .technical {
    padding: 12px;
  }

  .tech-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .tech-grid.three-col {
    grid-template-columns: repeat(3, 1fr);
  }

  .levels-container {
    flex-direction: column;
  }

  .pattern-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
