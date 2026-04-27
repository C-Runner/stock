<template>
  <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
    <template #header>
      <div class="card-header">
        <n-icon size="20"><IconTrend /></n-icon>
        <span>Support & Resistance</span>
      </div>
    </template>
    <p class="period-desc">{{ getSupportResistanceSummary() }}</p>
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

const getSupportResistanceSummary = () => {
  const levels = props.technical.levels
  if (!levels) return 'No significant support or resistance levels detected.'

  const topResistanceItem = levels.resistance[0]
  const topSupportItem = levels.support[0]
  const topResistance = topResistanceItem ? `¥${topResistanceItem.price.toFixed(2)} (${topResistanceItem.strength})` : 'none'
  const topSupport = topSupportItem ? `¥${topSupportItem.price.toFixed(2)} (${topSupportItem.strength})` : 'none'
  const trendLine = levels.trendLine ? `${levels.trendLine.direction.toUpperCase()} trend (${levels.trendLine.strength})` : 'no clear trend line'

  return `Resistance: ${topResistance}\nSupport: ${topSupport}\nTrend: ${trendLine}`
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
