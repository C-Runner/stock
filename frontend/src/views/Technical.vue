<template>
  <div class="technical">
    <BackgroundOrbs />

    <n-space vertical :size="12" class="content" :style="{ width: '100%' }">
      <div class="page-header">
        <div class="header-info">
          <h1 class="page-title">Technical Analysis</h1>
          <p class="subtitle" v-if="quote">{{ quote.name }} ({{ stockCode }})</p>
        </div>
      </div>

      <n-spin v-if="loading" show description="Loading..." />

      <template v-else-if="technical">
        <RecommendationCard :technical="technical" :quote="quote" />

        <PeriodAnalysis :technical="technical" />

        <LevelsSection :technical="technical" />

        <PatternSection :technical="technical" />

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
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard, NSpace, NSpin, NIcon
} from 'naive-ui'
import { stockApi, type TechnicalAnalysis, type StockQuote } from '../api'
import { IconChart } from '../components/icons'
import KLineChart from '../components/KLineChart.vue'
import { RecommendationCard, PeriodAnalysis, LevelsSection, PatternSection } from '../components/technical'
import BackgroundOrbs from '../components/BackgroundOrbs.vue'

const route = useRoute()
const router = useRouter()

const stockCode = ref(route.params.code as string)
const loading = ref(false)
const error = ref('')
const technical = ref<TechnicalAnalysis | null>(null)
const quote = ref<StockQuote | null>(null)

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

let touchStartX = 0

const handleTouchStart = (e: TouchEvent) => {
  if (e.touches.length > 0) {
    touchStartX = e.touches[0]?.clientX ?? 0
  }
}

const handleTouchEnd = (e: TouchEvent) => {
  const touchEndX = e.changedTouches[0]?.clientX ?? 0
  const swipeDistance = touchEndX - touchStartX
  if (touchStartX < 50 && swipeDistance > 80) {
    router.back()
  }
}

onMounted(() => {
  fetchData()
  document.addEventListener('touchstart', handleTouchStart, { passive: true })
  document.addEventListener('touchend', handleTouchEnd, { passive: true })
})

onUnmounted(() => {
  document.removeEventListener('touchstart', handleTouchStart)
  document.removeEventListener('touchend', handleTouchEnd)
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

.content {
  position: relative;
  align-items: stretch;
  width: 100%;
}

.analysis-card:hover {
  border-color: rgba(255, 255, 255, 0.15) !important;
  box-shadow: 0 8px 32px rgba(16, 185, 129, 0.1) !important;
}

.card-header :deep(.n-icon) {
  color: #10b981 !important;
}

@media (max-width: 768px) {
  .technical {
    padding: 12px;
  }
}
</style>
