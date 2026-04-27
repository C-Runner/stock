<template>
  <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
    <template #header>
      <div class="card-header">
        <n-icon size="20"><IconTrend /></n-icon>
        <span>Technical Analysis by Period</span>
      </div>
    </template>

    <div class="period-vertical">
      <div class="period-row short">
        <div class="period-row-header">
          <span class="period-label">Short-term</span>
          <span class="period-badge short">1-10 Days</span>
        </div>
        <p class="period-desc">{{ getShortTermSummary() }}</p>
      </div>

      <div class="period-row medium">
        <div class="period-row-header">
          <span class="period-label">Medium-term</span>
          <span class="period-badge medium">10-30 Days</span>
        </div>
        <p class="period-desc">{{ getMediumTermSummary() }}</p>
      </div>

      <div class="period-row long">
        <div class="period-row-header">
          <span class="period-label">Long-term</span>
          <span class="period-badge long">30-180 Days</span>
        </div>
        <p class="period-desc">{{ getLongTermSummary() }}</p>
      </div>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { NCard, NIcon } from 'naive-ui'
import { IconTrend } from '../icons'
import type { TechnicalAnalysis } from '../../api'

interface Props {
  technical: TechnicalAnalysis
}

const props = defineProps<Props>()

const getLatestValue = (values: number[]): number => {
  if (!values || values.length === 0) return 0
  return values[values.length - 1] ?? 0
}

const getShortTermMA = () => props.technical.ma.filter(m => m.period === 5) || []
const getShortTermRSI = () => props.technical.rsi.filter(r => r.period === 6) || []
const getMediumTermMA = () => props.technical.ma.filter(m => m.period === 10 || m.period === 20) || []
const getMediumTermEMA = () => props.technical.ema.filter(e => e.period === 12) || []
const getMediumTermRSI = () => props.technical.rsi.filter(r => r.period === 12) || []
const getLongTermMA = () => props.technical.ma.filter(m => m.period === 60) || []
const getLongTermEMA = () => props.technical.ema.filter(e => e.period === 26) || []
const getLongTermRSI = () => props.technical.rsi.filter(r => r.period === 24) || []

const getShortTermSummary = () => {
  const ma5 = getLatestValue(getShortTermMA()[0]?.values || [])
  const rsi6 = getLatestValue(getShortTermRSI()[0]?.values || [])
  const k = getLatestValue(props.technical.kdj.k || [])
  const d = getLatestValue(props.technical.kdj.d || [])
  const dif = getLatestValue(props.technical.macd.dif || [])
  const dea = getLatestValue(props.technical.macd.dea || [])
  const prices = props.technical.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const wr = getLatestValue(props.technical.wr.find(w => w.period === 10)?.values || [])

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

  if (props.technical.obv.trend === 'rising') signals.push('OBV 上升，成交量配合上涨')
  else if (props.technical.obv.trend === 'falling') signals.push('OBV 下降，量价背离风险')

  return signals.length > 0 ? signals.join('；') : '短期趋势不明，建议观望'
}

const getMediumTermSummary = () => {
  const ma10 = getLatestValue(getMediumTermMA()[0]?.values || [])
  const ma20 = getLatestValue(getMediumTermMA()[1]?.values || [])
  const ema12 = getLatestValue(getMediumTermEMA()[0]?.values || [])
  const rsi12 = getLatestValue(getMediumTermRSI()[0]?.values || [])
  const upper = getLatestValue(props.technical.boll.upper || [])
  const lower = getLatestValue(props.technical.boll.lower || [])
  const mid = getLatestValue(props.technical.boll.mid || [])
  const prices = props.technical.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const vwap = getLatestValue(props.technical.vwap.values || [])

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
  const prices = props.technical.recentPrices || []
  const lastPrice = prices[prices.length - 1]
  const currentPrice = lastPrice ? lastPrice.close : 0
  const plusDI = props.technical.dmi.plusDI || 0
  const minusDI = props.technical.dmi.minusDI || 0
  const adx = props.technical.dmi.adx || 0

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
</script>

<style scoped>
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

.period-desc {
  margin: 0;
  font-size: 12px;
  line-height: 1.9;
  color: rgba(255, 255, 255, 0.8);
  white-space: pre-line;
}
</style>
