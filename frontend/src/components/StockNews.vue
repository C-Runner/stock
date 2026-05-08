<template>
  <div class="stock-news">
    <n-spin v-if="loading" />
    <template v-else>
      <!-- Sentiment Summary Card -->
      <n-card v-if="sentiment" class="sentiment-card" :bordered="false">
        <div class="sentiment-header">
          <span class="sentiment-title">News Sentiment Analysis</span>
          <span class="news-count">{{ sentiment.newsCount }} news articles</span>
        </div>
        <div class="sentiment-body">
          <div class="sentiment-score-section">
            <div class="sentiment-score" :class="sentimentClass">
              {{ sentiment.overallScore > 0 ? '+' : '' }}{{ sentiment.overallScore.toFixed(0) }}
            </div>
            <div class="sentiment-label" :class="sentimentClass">
              {{ sentimentLabel }}
            </div>
          </div>
          <div class="sentiment-bar">
            <div class="bar-positive" :style="{ width: positivePercent + '%' }"></div>
            <div class="bar-negative" :style="{ width: negativePercent + '%' }"></div>
          </div>
          <div class="sentiment-stats">
            <div class="stat positive">
              <span class="stat-value">{{ sentiment.positiveCount }}</span>
              <span class="stat-label">Positive</span>
            </div>
            <div class="stat neutral">
              <span class="stat-value">{{ sentiment.neutralCount }}</span>
              <span class="stat-label">Neutral</span>
            </div>
            <div class="stat negative">
              <span class="stat-value">{{ sentiment.negativeCount }}</span>
              <span class="stat-label">Negative</span>
            </div>
          </div>
        </div>
        <div class="sentiment-time" v-if="sentiment.latestNewsTime">
          Latest: {{ formatTime(sentiment.latestNewsTime) }}
        </div>
      </n-card>

      <!-- Error State -->
      <n-result v-if="error" status="error" title="Failed to load news" :description="error">
      </n-result>

      <!-- Empty State -->
      <div v-else-if="newsItems.length === 0 && !error" class="no-news">
        No news available for this stock
      </div>

      <!-- News List -->
      <div v-if="newsItems.length > 0" class="news-list">
        <h3 class="list-title">Recent News</h3>
        <n-card
          v-for="item in newsItems"
          :key="item.id"
          class="news-item"
          :bordered="false"
          hoverable
          @click="openNews(item.sourceUrl)"
        >
          <div class="news-header">
            <n-tag
              :type="item.sentiment === 'positive' ? 'success' : item.sentiment === 'negative' ? 'error' : 'default'"
              size="small"
              round
            >
              {{ item.sentiment === 'positive' ? 'Positive' : item.sentiment === 'negative' ? 'Negative' : 'Neutral' }}
            </n-tag>
            <span class="news-source">{{ item.source }}</span>
            <span class="news-time">{{ formatTime(item.publishTime) }}</span>
          </div>
          <h4 class="news-title">{{ decodeTitle(item.title) }}</h4>
        </n-card>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { NCard, NTag, NSpin, NResult } from 'naive-ui'
import { newsApi, type NewsItem, type NewsSentimentResult } from '../api'

const props = defineProps<{
  code: string
}>()

const loading = ref(true)
const error = ref('')
const newsItems = ref<NewsItem[]>([])
const sentiment = ref<NewsSentimentResult | null>(null)

const sentimentClass = computed(() => {
  if (!sentiment.value) return ''
  if (sentiment.value.overallScore > 15) return 'positive'
  if (sentiment.value.overallScore < -15) return 'negative'
  return 'neutral'
})

const sentimentLabel = computed(() => {
  if (!sentiment.value) return ''
  if (sentiment.value.overallScore > 30) return 'Very Bullish'
  if (sentiment.value.overallScore > 15) return 'Bullish'
  if (sentiment.value.overallScore > -15 && sentiment.value.overallScore <= 15) return 'Neutral'
  if (sentiment.value.overallScore > -30) return 'Bearish'
  return 'Very Bearish'
})

const positivePercent = computed(() => {
  if (!sentiment.value || sentiment.value.newsCount === 0) return 0
  return (sentiment.value.positiveCount / sentiment.value.newsCount) * 100
})

const negativePercent = computed(() => {
  if (!sentiment.value || sentiment.value.newsCount === 0) return 0
  return (sentiment.value.negativeCount / sentiment.value.newsCount) * 100
})

const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return timeStr
  }
}

const decodeTitle = (title: string) => {
  if (!title) return ''
  // Handle URL encoded Chinese characters
  try {
    return decodeURIComponent(title.replace(/\\u/g, '%u'))
  } catch {
    return title
  }
}

const openNews = (url: string) => {
  if (url) window.open(url, '_blank')
}

const loadNews = async () => {
  loading.value = true
  error.value = ''
  try {
    // Fetch both news list and sentiment
    const [newsRes, sentimentRes] = await Promise.all([
      newsApi.getNews(props.code, 20),
      newsApi.getSentiment(props.code)
    ])
    newsItems.value = newsRes.newsItems || []
    sentiment.value = sentimentRes
  } catch (e: any) {
    console.error('Failed to load news:', e)
    error.value = e.message || 'Failed to load news data'
  } finally {
    loading.value = false
  }
}

onMounted(loadNews)
</script>

<style scoped>
.stock-news {
  width: 100%;
  padding: 0 4px;
}

.sentiment-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.1));
  border: 1px solid rgba(99, 102, 241, 0.2);
  margin-bottom: 16px;
}

.sentiment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.sentiment-title {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.news-count {
  font-size: 12px;
  color: #94a3b8;
}

.sentiment-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sentiment-score-section {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.sentiment-score {
  font-size: 42px;
  font-weight: bold;
  line-height: 1;
}

.sentiment-score.positive { color: #22c55e; }
.sentiment-score.negative { color: #ef4444; }
.sentiment-score.neutral { color: #94a3b8; }

.sentiment-label {
  font-size: 16px;
  font-weight: 500;
}

.sentiment-label.positive { color: #22c55e; }
.sentiment-label.negative { color: #ef4444; }
.sentiment-label.neutral { color: #94a3b8; }

.sentiment-bar {
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  display: flex;
}

.bar-positive {
  background: linear-gradient(90deg, #22c55e, #4ade80);
  height: 100%;
  transition: width 0.3s ease;
}

.bar-negative {
  background: linear-gradient(90deg, #f97316, #ef4444);
  height: 100%;
  transition: width 0.3s ease;
}

.sentiment-stats {
  display: flex;
  justify-content: space-around;
  gap: 16px;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
}

.stat-label {
  font-size: 11px;
  color: #94a3b8;
}

.stat.positive .stat-value { color: #22c55e; }
.stat.neutral .stat-value { color: #94a3b8; }
.stat.negative .stat-value { color: #ef4444; }

.sentiment-time {
  margin-top: 12px;
  font-size: 11px;
  color: #64748b;
  text-align: right;
}

.news-list {
  margin-top: 8px;
}

.list-title {
  font-size: 13px;
  font-weight: 600;
  color: #94a3b8;
  margin: 0 0 12px 4px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.news-item {
  margin-bottom: 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s ease;
}

.news-item:hover {
  background: rgba(255, 255, 255, 0.06);
  transform: translateX(4px);
}

.news-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.news-source {
  font-size: 11px;
  color: #64748b;
}

.news-time {
  font-size: 11px;
  color: #475569;
  margin-left: auto;
}

.news-title {
  margin: 0;
  font-size: 13px;
  color: #e2e8f0;
  line-height: 1.4;
  font-weight: 400;
}

.no-news {
  text-align: center;
  color: #64748b;
  padding: 32px 16px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 12px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
}
</style>
