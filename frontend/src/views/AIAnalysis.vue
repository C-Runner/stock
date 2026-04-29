<template>
  <div class="ai-analysis">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
    </div>

    <div class="page-header">
      <button class="back-btn" @click="goBack">
        <n-icon size="20"><IconArrowLeft /></n-icon>
      </button>
      <div class="header-info">
        <h1 class="page-title">AI Analysis</h1>
        <p class="subtitle">{{ report?.name || code }} ({{ code }})</p>
      </div>
    </div>

    <n-space vertical :size="12" class="content" :style="{ width: '100%' }">
      <n-spin v-if="loading" show description="Generating AI analysis..." />
      <n-alert v-else-if="error" type="error">{{ error }}</n-alert>

      <template v-else-if="report">
        <!-- Heuristic Mode Notice -->
        <div v-if="report.analysisMethod === 'heuristic'" class="heuristic-notice">
          <n-icon size="16"><IconLight /></n-icon>
          <span>Algorithm analysis (AI not configured)</span>
        </div>

        <!-- Cache Notice -->
        <div v-else-if="report.fromCache" class="cache-notice">
          <n-icon size="16"><IconClock /></n-icon>
          <span>Cached result (refreshes every 5 minutes)</span>
        </div>

        <!-- Raw AI Analysis Output (if available) -->
        <n-card v-if="report.rawAnalysis" class="analysis-card raw-analysis-card" :bordered="false" :content-style="{ padding: '20px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconRobot /></n-icon>
              <span>AI Analysis</span>
            </div>
          </template>
          <div class="raw-analysis-content">
            <pre>{{ report.rawAnalysis }}</pre>
          </div>
        </n-card>

        <!-- AI Analysis Results - Top Section (Description, Evaluation, Suggestions) -->

        <!-- Analysis Summary -->
        <n-card class="analysis-card summary-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconChart /></n-icon>
              <span>Analysis Summary</span>
            </div>
          </template>
          <p class="summary-text">{{ report.summary }}</p>
        </n-card>

        <!-- Investment Advice -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconLight /></n-icon>
              <span>Investment Advice</span>
            </div>
          </template>

          <div class="advice-overall" :class="getAdviceClass(report.investmentAdvice.riskLevel)">
            {{ report.investmentAdvice.overallAdvice }}
          </div>

          <div class="advice-grid">
            <div class="advice-section">
              <h4 class="advice-title">Entry Points</h4>
              <ul class="advice-list">
                <li v-for="(point, idx) in report.investmentAdvice.entryPoints" :key="idx">{{ point }}</li>
              </ul>
            </div>
            <div class="advice-section">
              <h4 class="advice-title">Exit Points</h4>
              <ul class="advice-list">
                <li v-for="(point, idx) in report.investmentAdvice.exitPoints" :key="idx">{{ point }}</li>
              </ul>
            </div>
          </div>

          <div class="advice-details">
            <div class="advice-detail">
              <span class="detail-label">Stop Loss</span>
              <span class="detail-value risk">{{ report.investmentAdvice.stopLoss }}</span>
            </div>
            <div class="advice-detail">
              <span class="detail-label">Risk Level</span>
              <span class="detail-value" :class="'risk-' + report.investmentAdvice.riskLevel">{{ report.investmentAdvice.riskLevel.toUpperCase() }}</span>
            </div>
            <div class="advice-detail">
              <span class="detail-label">Time Horizon</span>
              <span class="detail-value">{{ report.investmentAdvice.timeHorizon }}</span>
            </div>
            <div class="advice-detail full-width">
              <span class="detail-label">Position Sizing</span>
              <span class="detail-value">{{ report.investmentAdvice.positionSizing }}</span>
            </div>
          </div>

          <div class="risk-warnings">
            <h4 class="warning-title">Risk Warnings</h4>
            <div class="warning-list">
              <div v-for="(warning, idx) in report.investmentAdvice.riskWarnings" :key="idx" class="warning-item">
                <span class="warning-icon">!</span>
                <span>{{ warning }}</span>
              </div>
            </div>
          </div>
        </n-card>

        <!-- Key Findings -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconLight /></n-icon>
              <span>Key Findings</span>
            </div>
          </template>

          <div class="findings-grid">
            <div class="findings-column highlights">
              <h3 class="column-title">Highlights</h3>
              <div v-for="(highlight, idx) in report.keyFindings.highlights" :key="idx" class="finding-item highlight">
                <div class="finding-icon">✦</div>
                <div class="finding-content">
                  <div class="finding-title">{{ highlight.title }}</div>
                  <div class="finding-context">{{ highlight.context }}</div>
                </div>
              </div>
            </div>

            <div class="findings-column risks">
              <h3 class="column-title">Risks</h3>
              <div v-for="(risk, idx) in report.keyFindings.risks" :key="idx" class="finding-item risk">
                <div class="finding-icon">⚠</div>
                <div class="finding-content">
                  <div class="finding-title">{{ risk.title }}</div>
                  <div class="finding-context">{{ risk.context }}</div>
                </div>
              </div>
            </div>
          </div>
        </n-card>

        <!-- Chart Analysis -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconChart /></n-icon>
              <span>Chart Analysis</span>
            </div>
          </template>

          <div class="chart-analysis-section">
            <div class="analysis-item">
              <div class="analysis-label">Trend</div>
              <div class="analysis-text">{{ report.chartAnalysis.trendInterpretation }}</div>
            </div>
            <div class="analysis-item">
              <div class="analysis-label">Support & Resistance</div>
              <div class="analysis-text">{{ report.chartAnalysis.supportResistance }}</div>
            </div>
            <div class="analysis-item">
              <div class="analysis-label">Volume</div>
              <div class="analysis-text">{{ report.chartAnalysis.volumeAnalysis }}</div>
            </div>
            <div class="analysis-item">
              <div class="analysis-label">Indicators</div>
              <div class="analysis-text">{{ report.chartAnalysis.indicatorSummary }}</div>
            </div>
            <div class="analysis-item">
              <div class="analysis-label">Patterns</div>
              <div class="analysis-text">{{ report.chartAnalysis.patternDescription }}</div>
            </div>
            <div class="analysis-item">
              <div class="analysis-label">Momentum</div>
              <div class="analysis-text">{{ report.chartAnalysis.momentumAnalysis }}</div>
            </div>
          </div>
        </n-card>

        <!-- Data Section - Below -->

        <!-- Section Divider -->
        <div class="section-divider">
          <span>Analysis Data</span>
        </div>

        <!-- Multi-dimensional Scores -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconChart /></n-icon>
              <span>Multi-dimensional Analysis</span>
            </div>
          </template>

          <div class="score-cards">
            <div class="score-card technical">
              <div class="score-header">
                <span class="score-label">Technical (30%)</span>
                <span :class="['trend-badge', report.scores.technical.trend]">{{ report.scores.technical.trend }}</span>
              </div>
              <div class="score-value">{{ report.scores.technical.score }}</div>
              <div class="score-bar">
                <div class="score-fill" :style="{ width: report.scores.technical.score + '%' }"></div>
              </div>
              <p class="score-summary">{{ report.scores.technical.summary }}</p>
            </div>

            <div class="score-card fundamental">
              <div class="score-header">
                <span class="score-label">Fundamental (30%)</span>
                <span :class="['trend-badge', report.scores.fundamental.trend]">{{ report.scores.fundamental.trend }}</span>
              </div>
              <div class="score-value">{{ report.scores.fundamental.score }}</div>
              <div class="score-bar">
                <div class="score-fill fundamental" :style="{ width: report.scores.fundamental.score + '%' }"></div>
              </div>
              <p class="score-summary">{{ report.scores.fundamental.summary }}</p>
            </div>

            <div class="score-card money-flow">
              <div class="score-header">
                <span class="score-label">Money Flow (25%)</span>
                <span :class="['trend-badge', report.scores.moneyFlow.trend]">{{ report.scores.moneyFlow.trend }}</span>
              </div>
              <div class="score-value">{{ report.scores.moneyFlow.score }}</div>
              <div class="score-bar">
                <div class="score-fill money-flow" :style="{ width: report.scores.moneyFlow.score + '%' }"></div>
              </div>
              <p class="score-summary">{{ report.scores.moneyFlow.summary }}</p>
            </div>

            <div class="score-card news-sentiment">
              <div class="score-header">
                <span class="score-label">News Sentiment (15%)</span>
                <span :class="['trend-badge', report.scores.newsSentiment.trend]">{{ report.scores.newsSentiment.trend }}</span>
              </div>
              <div class="score-value">{{ report.scores.newsSentiment.score }}</div>
              <div class="score-bar">
                <div class="score-fill news-sentiment" :style="{ width: report.scores.newsSentiment.score + '%' }"></div>
              </div>
              <p class="score-summary">{{ report.scores.newsSentiment.summary }}</p>
            </div>
          </div>

          <n-divider />

          <div class="composite-section">
            <div class="composite-card">
              <div class="composite-label">Composite Score</div>
              <div class="composite-value">{{ report.scores.compositeScore }}</div>
              <div class="composite-bar">
                <div class="composite-fill" :style="{ width: report.scores.compositeScore + '%' }"></div>
              </div>
            </div>
            <div class="composite-card anxiety">
              <div class="composite-label">Anxiety Index</div>
              <div class="composite-value">{{ report.scores.anxietyIndex }}</div>
              <div class="anxiety-badge" :class="getAnxietyClass(report.scores.anxietyIndex)">
                {{ getAnxietyLabel(report.scores.anxietyIndex) }}
              </div>
            </div>
            <div class="composite-card attention">
              <div class="composite-label">Attention Level</div>
              <div class="attention-value" :class="report.scores.attentionLevel">
                {{ report.scores.attentionLevel.toUpperCase() }}
              </div>
            </div>
          </div>
        </n-card>

        <!-- Similar Patterns -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconTrend /></n-icon>
              <span>Similar K-line Patterns</span>
            </div>
          </template>

          <p class="section-desc">Historical patterns with most similar price movement trajectory</p>

          <div v-if="report.similarPatterns.length === 0" class="no-data">
            Insufficient historical data for pattern matching
          </div>

          <div v-else class="patterns-list">
            <div v-for="pattern in report.similarPatterns" :key="pattern.patternId" class="pattern-card">
              <div class="pattern-header">
                <span class="pattern-date">{{ pattern.startDate }} ~ {{ pattern.endDate }}</span>
                <span class="pattern-similarity">{{ pattern.similarity.toFixed(1) }}% Match</span>
              </div>
              <div class="pattern-stats">
                <div class="stat-item">
                  <span class="stat-label">Price Change</span>
                  <span class="stat-value" :class="pattern.priceChange >= 0 ? 'up' : 'down'">
                    {{ pattern.priceChange >= 0 ? '+' : '' }}{{ pattern.priceChange.toFixed(2) }}%
                  </span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Next 5D Win Rate</span>
                  <span class="stat-value">{{ pattern.next5DayWinRate.toFixed(1) }}%</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">Next 20D Win Rate</span>
                  <span class="stat-value">{{ pattern.next20DayWinRate.toFixed(1) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </n-card>

        <!-- Institutional Sentiment -->
        <n-card class="analysis-card" :bordered="false" :content-style="{ padding: '16px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconRobot /></n-icon>
              <span>Institutional Research Sentiment</span>
            </div>
          </template>

          <p class="section-desc">Last {{ report.institutionalSentiment.period }} aggregated from {{ report.institutionalSentiment.totalReports }} research reports</p>

          <div class="sentiment-grid">
            <div class="sentiment-chart">
              <div class="sentiment-bar">
                <div class="bar-segment buy" :style="{ width: report.institutionalSentiment.buyRatio + '%' }"></div>
                <div class="bar-segment hold" :style="{ width: report.institutionalSentiment.holdRatio + '%' }"></div>
                <div class="bar-segment sell" :style="{ width: report.institutionalSentiment.sellRatio + '%' }"></div>
              </div>
              <div class="sentiment-legend">
                <span class="legend-item"><span class="dot buy"></span>Buy {{ report.institutionalSentiment.buyRatio.toFixed(0) }}%</span>
                <span class="legend-item"><span class="dot hold"></span>Hold {{ report.institutionalSentiment.holdRatio.toFixed(0) }}%</span>
                <span class="legend-item"><span class="dot sell"></span>Sell {{ report.institutionalSentiment.sellRatio.toFixed(0) }}%</span>
              </div>
            </div>

            <div class="sentiment-stats">
              <div class="sentiment-stat">
                <span class="stat-label">Trend</span>
                <span class="stat-value" :class="report.institutionalSentiment.ratioChange">
                  {{ report.institutionalSentiment.ratioChange }}
                </span>
              </div>
              <div class="sentiment-stat">
                <span class="stat-label">Consensus Rating</span>
                <span class="stat-value rating" :class="report.institutionalSentiment.consensusRating">
                  {{ report.institutionalSentiment.consensusRating }}
                </span>
              </div>
            </div>
          </div>
        </n-card>

        <div class="generated-at">
          Generated at: {{ report.generatedAt }}
        </div>
      </template>
    </n-space>

    <div class="bottom-tabs">
      <div class="tab-item" @click="router.push('/')">
        <div class="tab-icon">
          <n-icon size="22"><IconHome /></n-icon>
        </div>
        <span class="tab-label">Portfolio</span>
      </div>
      <div class="tab-item" @click="router.push(`/analysis/${code}`)">
        <div class="tab-icon">
          <n-icon size="22"><IconChart /></n-icon>
        </div>
        <span class="tab-label">Analysis</span>
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
  NCard, NSpace, NSpin, NAlert, NDivider, NIcon
} from 'naive-ui'
import { aiApi, type AIAnalysisReport } from '../api'
import { IconArrowLeft, IconChart, IconClock, IconLight, IconTrend, IconRobot, IconHome, IconStar } from '../components/icons'

const route = useRoute()
const router = useRouter()

const code = route.params.code as string
const report = ref<AIAnalysisReport | null>(null)
const loading = ref(true)
const error = ref('')

const fetchAnalysis = async () => {
  loading.value = true
  error.value = ''
  try {
    report.value = await aiApi.getAnalysis(code)
  } catch (e: any) {
    error.value = e.message || 'Failed to load AI analysis'
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const getAnxietyClass = (index: number) => {
  if (index < 30) return 'low'
  if (index < 70) return 'medium'
  return 'high'
}

const getAnxietyLabel = (index: number) => {
  if (index < 30) return 'Low Anxiety'
  if (index < 70) return 'Moderate'
  return 'High Anxiety'
}

const getAdviceClass = (riskLevel: string) => {
  if (riskLevel === 'low' || riskLevel === 'medium') return 'advice-bullish'
  if (riskLevel === 'high') return 'advice-bearish'
  return 'advice-neutral'
}

onMounted(fetchAnalysis)
</script>

<style scoped>
.ai-analysis {
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
  padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
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
  width: 100%;
}

.content > * {
  flex-shrink: 0;
  margin-bottom: 0 !important;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
  position: relative;
  flex: 0 0 auto;
}

.back-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.8);
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.15);
}

.header-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  margin: 0;
  font-size: 26px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.5px;
}

.subtitle {
  margin: 0;
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
}

.cache-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(251, 191, 36, 0.1);
  border: 1px solid rgba(251, 191, 36, 0.2);
  border-radius: 12px;
  font-size: 13px;
  color: #fbbf24;
}

.cache-notice :deep(.n-icon) {
  color: #fbbf24;
}

.heuristic-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 12px;
  font-size: 13px;
  color: #818cf8;
}

.heuristic-notice :deep(.n-icon) {
  color: #818cf8;
}

.summary-card .summary-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.85);
  padding: 8px 0;
}

.analysis-card {
  --n-border-radius: 20px !important;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
  padding: 0 !important;
  margin: 0 !important;
  position: relative !important;
  overflow: visible !important;
  display: block !important;
}

.analysis-card::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: 21px;
  padding: 1px;
  background: linear-gradient(
    135deg,
    rgba(99, 102, 241, 0.3),
    rgba(139, 92, 246, 0.1) 40%,
    rgba(139, 92, 246, 0.1) 60%,
    rgba(99, 102, 241, 0.3)
  );
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.analysis-card:hover::before {
  opacity: 1;
}

.analysis-card > * {
  box-sizing: border-box;
}

.analysis-card :deep(.n-base-card) {
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  border-radius: inherit !important;
  overflow: visible !important;
  padding: 0 !important;
  margin: 0 !important;
}

.analysis-card :deep(.n-card) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  border-radius: 20px !important;
  overflow: visible !important;
  padding: 0 !important;
  margin: 0 !important;
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
  border: none !important;
  padding: 16px 20px !important;
  margin: 0 !important;
  border-radius: inherit !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06) !important;
}

.analysis-card :deep(.n-card-wrapper) {
  margin: 0 !important;
  padding: 0 !important;
  border-radius: inherit !important;
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

.score-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.score-card {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  border-left: 3px solid;
  transition: all 0.2s ease;
}

.score-card:hover {
  background: rgba(255, 255, 255, 0.05);
  transform: translateY(-2px);
}

.score-card.technical { border-left-color: #3b82f6; }
.score-card.fundamental { border-left-color: #10b981; }
.score-card.money-flow { border-left-color: #f59e0b; }
.score-card.news-sentiment { border-left-color: #8b5cf6; }

.score-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.score-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 500;
}

.trend-badge {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 8px;
  text-transform: capitalize;
}

.trend-badge.improving { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.trend-badge.stable { background: rgba(255, 255, 255, 0.1); color: rgba(255, 255, 255, 0.5); }
.trend-badge.declining { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

.score-value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
}

.score-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 8px;
}

.score-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #6366f1);
  border-radius: 2px;
  transition: width 0.5s ease;
}

.score-fill.fundamental { background: linear-gradient(90deg, #10b981, #34d399); }
.score-fill.money-flow { background: linear-gradient(90deg, #f59e0b, #fbbf24); }
.score-fill.news-sentiment { background: linear-gradient(90deg, #8b5cf6, #a78bfa); }

.score-summary {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
  line-height: 1.4;
}

.composite-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-top: 16px;
}

.composite-card {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  text-align: center;
}

.composite-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.composite-value {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
}

.composite-bar {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  overflow: hidden;
}

.composite-fill {
  height: 100%;
  background: linear-gradient(90deg, #10b981, #34d399);
  border-radius: 3px;
  transition: width 0.5s ease;
}

.composite-card.anxiety .composite-value {
  color: #ef4444;
}

.anxiety-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
}

.anxiety-badge.low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.anxiety-badge.medium { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
.anxiety-badge.high { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

.attention-value {
  font-size: 18px;
  font-weight: 700;
  padding: 6px 16px;
  border-radius: 10px;
  display: inline-block;
}

.attention-value.high { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.attention-value.medium { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
.attention-value.low { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

.section-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0 0 16px 0;
}

.findings-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.findings-column h3.column-title {
  font-size: 14px;
  margin: 0 0 12px 0;
  font-weight: 600;
}

.findings-column.highlights h3.column-title { color: #10b981; }
.findings-column.risks h3.column-title { color: #ef4444; }

.finding-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 10px;
}

.finding-item.highlight {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.15);
}

.finding-item.risk {
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.15);
}

.finding-icon {
  font-size: 18px;
  line-height: 1;
  flex-shrink: 0;
}

.finding-title {
  font-weight: 600;
  font-size: 13px;
  color: #fff;
  margin-bottom: 4px;
}

.finding-context {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.4;
}

.no-data {
  padding: 40px;
  text-align: center;
  color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.patterns-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pattern-card {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  transition: all 0.2s ease;
}

.pattern-card:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(99, 102, 241, 0.2);
}

.pattern-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.pattern-date {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.pattern-similarity {
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
  padding: 4px 10px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
}

.pattern-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.stat-value.up { color: #ff6b6b; }
.stat-value.down { color: #38ef7d; }

.sentiment-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.sentiment-chart {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  padding: 16px;
}

.sentiment-bar {
  display: flex;
  height: 20px;
  border-radius: 10px;
  overflow: hidden;
}

.bar-segment { transition: width 0.5s ease; }
.bar-segment.buy { background: #10b981; }
.bar-segment.hold { background: #f59e0b; }
.bar-segment.sell { background: #ef4444; }

.sentiment-legend {
  display: flex;
  gap: 16px;
  margin-top: 12px;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.6);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dot.buy { background: #10b981; }
.dot.hold { background: #f59e0b; }
.dot.sell { background: #ef4444; }

.sentiment-stats {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sentiment-stat {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sentiment-stat .stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.sentiment-stat .stat-value {
  font-size: 13px;
  font-weight: 600;
  text-transform: capitalize;
}

.sentiment-stat .stat-value.improving { color: #10b981; }
.sentiment-stat .stat-value.stable { color: rgba(255, 255, 255, 0.6); }
.sentiment-stat .stat-value.deteriorating { color: #ef4444; }

.sentiment-stat .stat-value.rating.buy { background: rgba(16, 185, 129, 0.15); color: #10b981; padding: 4px 12px; border-radius: 8px; }
.sentiment-stat .stat-value.rating.hold { background: rgba(251, 191, 36, 0.15); color: #fbbf24; padding: 4px 12px; border-radius: 8px; }
.sentiment-stat .stat-value.rating.sell { background: rgba(239, 68, 68, 0.15); color: #ef4444; padding: 4px 12px; border-radius: 8px; }

.generated-at {
  text-align: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
  margin-top: 8px;
  padding: 16px;
}

.section-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin: 8px 0;
}

.section-divider::before,
.section-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(99, 102, 241, 0.3), transparent);
}

.section-divider span {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 500;
}

.raw-analysis-card {
  margin-bottom: 16px !important;
}

.raw-analysis-content {
  font-size: 14px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.85);
}

.raw-analysis-content pre {
  margin: 0;
  padding: 0;
  font-family: inherit;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* Chart Analysis Styles */
.chart-analysis-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.analysis-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.analysis-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 500;
}

.analysis-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
}

/* Investment Advice Styles */
.advice-overall {
  padding: 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  text-align: center;
  margin-bottom: 16px;
}

.advice-overall.advice-bullish {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.advice-overall.advice-neutral {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(251, 191, 36, 0.3);
}

.advice-overall.advice-bearish {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.advice-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}

.advice-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 14px;
}

.advice-title {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 10px 0;
  font-weight: 500;
}

.advice-list {
  margin: 0;
  padding-left: 16px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.6;
}

.advice-list li {
  margin-bottom: 4px;
}

.advice-details {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 16px;
}

.advice-detail {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.advice-detail.full-width {
  grid-column: span 3;
}

.detail-label {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-value {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.detail-value.risk {
  color: #ef4444;
}

.detail-value.risk-low { color: #10b981; }
.detail-value.risk-medium { color: #fbbf24; }
.detail-value.risk-high { color: #ef4444; }

.risk-warnings {
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.15);
  border-radius: 12px;
  padding: 14px;
}

.warning-title {
  font-size: 12px;
  color: #ef4444;
  margin: 0 0 10px 0;
  font-weight: 600;
}

.warning-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.warning-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.5;
}

.warning-icon {
  width: 18px;
  height: 18px;
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  flex-shrink: 0;
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

@media (max-width: 768px) {
  .ai-analysis {
    padding: 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }

  .score-cards {
    grid-template-columns: 1fr;
  }

  .composite-section {
    grid-template-columns: 1fr;
  }

  .findings-grid {
    grid-template-columns: 1fr;
  }

  .sentiment-grid {
    grid-template-columns: 1fr;
  }

  .pattern-stats {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}
</style>