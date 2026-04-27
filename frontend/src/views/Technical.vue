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
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard, NSpace, NSpin, NAlert, NIcon
} from 'naive-ui'
import { stockApi, type TechnicalAnalysis, type StockQuote } from '../api'
import { IconChart, IconArrowBack } from '../components/icons'
import KLineChart from '../components/KLineChart.vue'
import { RecommendationCard, PeriodAnalysis, LevelsSection, PatternSection } from '../components/technical'

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

@media (max-width: 768px) {
  .technical {
    padding: 12px;
  }
}
</style>
