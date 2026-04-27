<template>
  <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
    <template #header>
      <div class="card-header">
        <n-icon size="20"><IconChart /></n-icon>
        <span>Pattern Recognition</span>
      </div>
    </template>
    <p class="period-desc">{{ getPatternSummary() }}</p>
  </n-card>
</template>

<script setup lang="ts">
import { NCard, NIcon } from 'naive-ui'
import { IconChart } from '../icons'
import type { TechnicalAnalysis } from '../../api'

interface Props {
  technical: TechnicalAnalysis
}

const props = defineProps<Props>()

const getPatternSummary = () => {
  const patterns = props.technical.patterns
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

.period-desc {
  margin: 0;
  font-size: 12px;
  line-height: 1.9;
  color: rgba(255, 255, 255, 0.8);
  white-space: pre-line;
}
</style>
