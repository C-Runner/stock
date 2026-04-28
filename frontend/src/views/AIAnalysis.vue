<template>
  <div class="ai-analysis-container">
    <div class="header">
      <button class="back-btn" @click="goBack">
        <span class="icon">←</span> Back
      </button>
      <div class="stock-info">
        <h1>{{ report?.name || code }} ({{ code }})</h1>
        <span class="ai-badge">AI Analysis</span>
      </div>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Generating AI analysis...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchAnalysis">Retry</button>
    </div>

    <div v-else-if="report" class="analysis-content">
      <!-- Cache Status -->
      <div v-if="report.fromCache" class="cache-notice">
        <span class="icon">📦</span> Cached result (refreshes every 5 minutes)
      </div>

      <!-- Multi-dimensional Scores -->
      <section class="scores-section">
        <h2>Multi-dimensional Analysis</h2>
        <div class="score-cards">
          <div class="score-card technical">
            <div class="score-header">
              <span class="label">Technical (30%)</span>
              <span class="trend" :class="report.scores.technical.trend">{{ report.scores.technical.trend }}</span>
            </div>
            <div class="score-value">{{ report.scores.technical.score }}</div>
            <div class="score-bar">
              <div class="score-fill" :style="{ width: report.scores.technical.score + '%' }"></div>
            </div>
            <p class="score-summary">{{ report.scores.technical.summary }}</p>
          </div>

          <div class="score-card fundamental">
            <div class="score-header">
              <span class="label">Fundamental (30%)</span>
              <span class="trend" :class="report.scores.fundamental.trend">{{ report.scores.fundamental.trend }}</span>
            </div>
            <div class="score-value">{{ report.scores.fundamental.score }}</div>
            <div class="score-bar">
              <div class="score-fill" :style="{ width: report.scores.fundamental.score + '%' }"></div>
            </div>
            <p class="score-summary">{{ report.scores.fundamental.summary }}</p>
          </div>

          <div class="score-card money-flow">
            <div class="score-header">
              <span class="label">Money Flow (25%)</span>
              <span class="trend" :class="report.scores.moneyFlow.trend">{{ report.scores.moneyFlow.trend }}</span>
            </div>
            <div class="score-value">{{ report.scores.moneyFlow.score }}</div>
            <div class="score-bar">
              <div class="score-fill" :style="{ width: report.scores.moneyFlow.score + '%' }"></div>
            </div>
            <p class="score-summary">{{ report.scores.moneyFlow.summary }}</p>
          </div>

          <div class="score-card news-sentiment">
            <div class="score-header">
              <span class="label">News Sentiment (15%)</span>
              <span class="trend" :class="report.scores.newsSentiment.trend">{{ report.scores.newsSentiment.trend }}</span>
            </div>
            <div class="score-value">{{ report.scores.newsSentiment.score }}</div>
            <div class="score-bar">
              <div class="score-fill" :style="{ width: report.scores.newsSentiment.score + '%' }"></div>
            </div>
            <p class="score-summary">{{ report.scores.newsSentiment.summary }}</p>
          </div>
        </div>

        <!-- Composite Score & Anxiety Index -->
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
            <div class="anxiety-indicator" :class="getAnxietyClass(report.scores.anxietyIndex)">
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
      </section>

      <!-- Key Findings -->
      <section class="findings-section">
        <h2>Key Findings</h2>

        <div class="findings-grid">
          <div class="findings-column highlights">
            <h3>Highlights</h3>
            <div v-for="(highlight, idx) in report.keyFindings.highlights" :key="idx" class="finding-item highlight">
              <div class="finding-icon">✦</div>
              <div class="finding-content">
                <div class="finding-title">{{ highlight.title }}</div>
                <div class="finding-context">{{ highlight.context }}</div>
              </div>
            </div>
          </div>

          <div class="findings-column risks">
            <h3>Risks</h3>
            <div v-for="(risk, idx) in report.keyFindings.risks" :key="idx" class="finding-item risk">
              <div class="finding-icon">⚠</div>
              <div class="finding-content">
                <div class="finding-title">{{ risk.title }}</div>
                <div class="finding-context">{{ risk.context }}</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Similar Patterns -->
      <section class="patterns-section">
        <h2>Similar K-line Patterns</h2>
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
              <div class="stat">
                <span class="stat-label">Price Change</span>
                <span class="stat-value" :class="pattern.priceChange >= 0 ? 'positive' : 'negative'">
                  {{ pattern.priceChange >= 0 ? '+' : '' }}{{ pattern.priceChange.toFixed(2) }}%
                </span>
              </div>
              <div class="stat">
                <span class="stat-label">Next 5D Win Rate</span>
                <span class="stat-value">{{ pattern.next5DayWinRate.toFixed(1) }}%</span>
              </div>
              <div class="stat">
                <span class="stat-label">Next 20D Win Rate</span>
                <span class="stat-value">{{ pattern.next20DayWinRate.toFixed(1) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Institutional Sentiment -->
      <section class="institutional-section">
        <h2>Institutional Research Sentiment</h2>
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
      </section>

      <div class="generated-at">
        Generated at: {{ report.generatedAt }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { aiApi, type AIAnalysisReport } from '../api'

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

onMounted(fetchAnalysis)
</script>

<style scoped>
.ai-analysis-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.back-btn:hover {
  background: #e0e0e0;
}

.stock-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stock-info h1 {
  font-size: 24px;
  margin: 0;
}

.ai-badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.loading, .error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.cache-notice {
  background: #fffbeb;
  border: 1px solid #fbbf24;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

section {
  margin-bottom: 32px;
}

section h2 {
  font-size: 20px;
  margin-bottom: 16px;
}

.score-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.score-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
}

.score-card.technical { border-left: 4px solid #3b82f6; }
.score-card.fundamental { border-left: 4px solid #10b981; }
.score-card.money-flow { border-left: 4px solid #f59e0b; }
.score-card.news-sentiment { border-left: 4px solid #8b5cf6; }

.score-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.score-header .label {
  font-size: 13px;
  color: #64748b;
}

.trend {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: capitalize;
}

.trend.improving { background: #d1fae5; color: #059669; }
.trend.stable { background: #e2e8f0; color: #64748b; }
.trend.declining { background: #fee2e2; color: #dc2626; }

.score-value {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
}

.score-bar {
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  margin: 8px 0;
  overflow: hidden;
}

.score-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #6366f1);
  border-radius: 3px;
  transition: width 0.5s ease;
}

.score-summary {
  font-size: 12px;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

.composite-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 20px;
}

.composite-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
  text-align: center;
}

.composite-label {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
}

.composite-value {
  font-size: 36px;
  font-weight: 700;
  color: #1e293b;
}

.composite-bar {
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  margin: 12px 0;
  overflow: hidden;
}

.composite-fill {
  height: 100%;
  background: linear-gradient(90deg, #10b981, #34d399);
  border-radius: 4px;
}

.anxiety .composite-value { color: #dc2626; }

.anxiety-indicator {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.anxiety-indicator.low { background: #d1fae5; color: #059669; }
.anxiety-indicator.medium { background: #fef3c7; color: #d97706; }
.anxiety-indicator.high { background: #fee2e2; color: #dc2626; }

.attention-value {
  font-size: 18px;
  font-weight: 700;
  padding: 4px 16px;
  border-radius: 8px;
  display: inline-block;
}

.attention-value.high { background: #d1fae5; color: #059669; }
.attention-value.medium { background: #fef3c7; color: #d97706; }
.attention-value.low { background: #fee2e2; color: #dc2626; }

.findings-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.findings-column h3 {
  font-size: 16px;
  margin-bottom: 12px;
}

.findings-column.highlights h3 { color: #059669; }
.findings-column.risks h3 { color: #dc2626; }

.finding-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.finding-item.highlight {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
}

.finding-item.risk {
  background: #fef2f2;
  border: 1px solid #fecaca;
}

.finding-icon {
  font-size: 20px;
  line-height: 1;
}

.finding-title {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.finding-context {
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
}

.patterns-section .section-desc,
.institutional-section .section-desc {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 16px;
}

.no-data {
  padding: 40px;
  text-align: center;
  color: #64748b;
  background: #f8fafc;
  border-radius: 12px;
}

.patterns-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pattern-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
}

.pattern-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.pattern-date {
  font-size: 13px;
  color: #64748b;
}

.pattern-similarity {
  background: #e0e7ff;
  color: #4f46e5;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.pattern-stats {
  display: flex;
  gap: 24px;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 11px;
  color: #94a3b8;
  text-transform: uppercase;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
}

.stat-value.positive { color: #059669; }
.stat-value.negative { color: #dc2626; }

.sentiment-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.sentiment-chart {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
}

.sentiment-bar {
  display: flex;
  height: 24px;
  border-radius: 12px;
  overflow: hidden;
}

.bar-segment {
  transition: width 0.5s ease;
}

.bar-segment.buy { background: #10b981; }
.bar-segment.hold { background: #f59e0b; }
.bar-segment.sell { background: #ef4444; }

.sentiment-legend {
  display: flex;
  gap: 16px;
  margin-top: 12px;
  font-size: 13px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot.buy { background: #10b981; }
.dot.hold { background: #f59e0b; }
.dot.sell { background: #ef4444; }

.sentiment-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.sentiment-stat {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sentiment-stat .stat-label {
  font-size: 13px;
  color: #64748b;
}

.sentiment-stat .stat-value {
  font-size: 14px;
  font-weight: 600;
  text-transform: capitalize;
}

.sentiment-stat .stat-value.improving { color: #059669; }
.sentiment-stat .stat-value.stable { color: #64748b; }
.sentiment-stat .stat-value.deteriorating { color: #dc2626; }

.sentiment-stat .stat-value.rating.buy { background: #d1fae5; color: #059669; padding: 4px 12px; border-radius: 8px; }
.sentiment-stat .stat-value.rating.hold { background: #fef3c7; color: #d97706; padding: 4px 12px; border-radius: 8px; }
.sentiment-stat .stat-value.rating.sell { background: #fee2e2; color: #dc2626; padding: 4px 12px; border-radius: 8px; }

.generated-at {
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
  margin-top: 32px;
}
</style>
