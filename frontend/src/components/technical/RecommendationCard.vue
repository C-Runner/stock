<template>
  <n-card class="analysis-card recommendation-card" :class="technical.recommendation.action">
    <div class="recommendation-top">
      <div class="stock-info-section">
        <span class="stock-name" v-if="quote">{{ quote.name }}</span>
        <span class="stock-code">({{ technical.code }})</span>
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
</template>

<script setup lang="ts">
import { NCard } from 'naive-ui'
import type { TechnicalAnalysis } from '../../api'
import type { StockQuote } from '../../api'

interface Props {
  technical: TechnicalAnalysis
  quote: StockQuote | null
}

defineProps<Props>()
</script>

<style scoped>
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
</style>
