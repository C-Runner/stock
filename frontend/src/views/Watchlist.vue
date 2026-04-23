<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NIcon,
  NSwitch, NSpin, NEmpty, NButtonGroup
} from 'naive-ui'
import { watchlistApi, stockApi, backupApi, type WatchlistItem, type StockQuote } from '../api'
import { IconRefresh, IconDelete, IconHome, IconPlus, IconBackup } from '../components/icons'
import StockSearch from '../components/StockSearch.vue'

const router = useRouter()

const watchlist = ref<WatchlistItem[]>([])
const quotes = ref<Map<string, StockQuote>>(new Map())
const loading = ref(false)
const refreshing = ref(false)
const backingUp = ref(false)
const backupStatus = ref('')

const autoRefresh = ref(true)
const refreshInterval = ref(30)
const lastRefresh = ref<Date | null>(null)
const refreshTimer = ref<ReturnType<typeof setInterval> | null>(null)

const showSearch = ref(false)

interface WatchlistRow {
  code: string
  name: string
  addedAt: string
  quote: StockQuote | null
}

const tableData = computed<WatchlistRow[]>(() =>
  watchlist.value.map(item => ({
    ...item,
    quote: quotes.value.get(item.code) || null
  }))
)

const fetchWatchlist = async () => {
  loading.value = true
  try {
    watchlist.value = await watchlistApi.getWatchlist()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchQuotes = async () => {
  if (watchlist.value.length === 0) return

  refreshing.value = true
  try {
    // Fetch quotes for all watchlist items in parallel
    const quotesData = await Promise.all(
      watchlist.value.map(item =>
        stockApi.getQuote(item.code).then(q => ({ code: item.code, quote: q })).catch(() => null)
      )
    )

    // Clear old quotes and set new ones, matching by code to avoid index misalignment
    quotes.value.clear()
    quotesData.forEach((result) => {
      if (result && result.quote) {
        quotes.value.set(result.code, result.quote)
      }
    })
    lastRefresh.value = new Date()
  } catch (error) {
    console.error(error)
  } finally {
    refreshing.value = false
  }
}

const handleRemove = async (code: string) => {
  try {
    await watchlistApi.removeFromWatchlist(code)
    watchlist.value = watchlist.value.filter(item => item.code !== code)
    quotes.value.delete(code)
  } catch (error) {
    console.error(error)
  }
}

const handleAddToWatchlist = async (code: string, name: string) => {
  try {
    await watchlistApi.addToWatchlist(code, name)
    await fetchWatchlist()
    const quote = await stockApi.getQuote(code)
    quotes.value.set(code, quote)
  } catch (error) {
    console.error(error)
  }
}

const handleBackup = async () => {
  backingUp.value = true
  backupStatus.value = 'Starting backup...'
  try {
    const result = await backupApi.triggerBackup()
    backupStatus.value = result.message
    setTimeout(() => {
      backupStatus.value = ''
    }, 5000)
  } catch (error) {
    console.error(error)
    backupStatus.value = 'Backup failed'
    setTimeout(() => {
      backupStatus.value = ''
    }, 3000)
  } finally {
    backingUp.value = false
  }
}

const startAutoRefresh = () => {
  if (refreshTimer.value) clearInterval(refreshTimer.value)
  refreshTimer.value = setInterval(fetchQuotes, refreshInterval.value * 1000)
}

const stopAutoRefresh = () => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }
}

watch(autoRefresh, (enabled) => {
  if (enabled) startAutoRefresh()
  else stopAutoRefresh()
})

watch(refreshInterval, () => {
  if (autoRefresh.value) startAutoRefresh()
})

onMounted(async () => {
  await fetchWatchlist()
  await fetchQuotes()
  if (autoRefresh.value) startAutoRefresh()
})

onUnmounted(() => stopAutoRefresh())
</script>

<template>
  <div class="watchlist">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>

    <div class="watchlist-header">
      <div class="header-left">
        <div class="logo">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>Watchlist</h1>
          <p class="subtitle">Track your favorite stocks</p>
        </div>
      </div>
      <div class="header-right">
        <n-button @click="handleBackup" :loading="backingUp" class="backup-btn" style="background: rgba(255, 255, 255, 0.05) !important;">
          <template #icon>
            <n-icon style="background: transparent !important;"><IconBackup /></n-icon>
          </template>
        </n-button>
        <n-button @click="fetchQuotes" :loading="refreshing" class="refresh-btn" style="background: rgba(255, 255, 255, 0.05) !important;">
          <template #icon>
            <n-icon style="background: transparent !important;"><IconRefresh /></n-icon>
          </template>
        </n-button>
      </div>
    </div>

    <div class="refresh-card">
      <div class="refresh-left">
        <div class="refresh-toggle">
          <span>Auto Refresh</span>
          <n-switch v-model:value="autoRefresh" />
        </div>
        <div class="refresh-interval" v-if="autoRefresh">
          <n-button-group size="small">
            <n-button
              :type="refreshInterval === 10 ? 'primary' : 'default'"
              @click="refreshInterval = 10"
            >
              10s
            </n-button>
            <n-button
              :type="refreshInterval === 30 ? 'primary' : 'default'"
              @click="refreshInterval = 30"
            >
              30s
            </n-button>
            <n-button
              :type="refreshInterval === 60 ? 'primary' : 'default'"
              @click="refreshInterval = 60"
            >
              60s
            </n-button>
          </n-button-group>
        </div>
      </div>
      <div class="backup-status" v-if="backupStatus">
        {{ backupStatus }}
      </div>
    </div>

    <div class="table-card">
      <div class="table-header">
        <span>Watchlist</span>
        <span class="last-refresh" v-if="lastRefresh">
          {{ lastRefresh.toLocaleTimeString() }}
        </span>
      </div>
      <div class="table-content">
        <n-spin :show="loading">
          <n-empty v-if="!loading && tableData.length === 0" description="No stocks in watchlist" />
          <div v-else class="watchlist-cards">
            <div
              v-for="row in tableData"
              :key="row.code"
              class="watchlist-card"
              @click="router.push({ path: `/analysis/${row.code}`, query: { from: '/watchlist' } })"
            >
              <div class="card-left">
                <div class="card-name">{{ row.name }}</div>
                <div class="card-code">{{ row.code }}</div>
              </div>
              <div class="card-right" v-if="row.quote">
                <div class="card-price">¥{{ row.quote.current.toFixed(2) }}</div>
                <div class="card-change" :class="row.quote.current >= row.quote.prevClose ? 'up' : 'down'">
                  {{ row.quote.current >= row.quote.prevClose ? '+' : '' }}{{ (((row.quote.current - row.quote.prevClose) / row.quote.prevClose) * 100).toFixed(2) }}%
                </div>
              </div>
              <div class="card-right" v-else>
                <div class="card-price">-</div>
                <div class="card-change">-</div>
              </div>
              <n-button
                size="small"
                type="error"
                quaternary
                @click.stop="handleRemove(row.code)"
                class="card-delete-btn"
              >
                <template #icon>
                  <n-icon><IconDelete /></n-icon>
                </template>
              </n-button>
            </div>
          </div>
        </n-spin>
      </div>
    </div>

    <StockSearch
      v-model:show="showSearch"
      @select="handleAddToWatchlist"
    />

    <div class="bottom-tabs">
      <div class="tab-item" @click="router.push('/')">
        <div class="tab-icon">
          <n-icon size="22"><IconHome /></n-icon>
        </div>
        <span class="tab-label">Portfolio</span>
      </div>
      <div class="tab-item primary" @click="showSearch = true">
        <div class="tab-icon">
          <n-icon size="24"><IconPlus /></n-icon>
        </div>
        <span class="tab-label">Add</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.watchlist {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  min-height: calc(100vh - 60px);
  padding: 16px;
  padding-bottom: calc(80px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  position: relative;
  overflow: visible;
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
  right: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: #8b5cf6;
  bottom: -100px;
  left: -100px;
  animation: float 10s ease-in-out infinite reverse;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: #06b6d4;
  top: 50%;
  left: 30%;
  transform: translate(-50%, -50%);
  animation: pulse 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}

@keyframes pulse {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.3; }
  50% { transform: translate(-50%, -50%) scale(1.2); opacity: 0.5; }
}

.watchlist-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  position: relative;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 36px;
  height: 36px;
  padding: 8px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo svg {
  width: 100%;
  height: 100%;
  color: #fff;
}

.header-text h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.5px;
}

.subtitle {
  margin: 2px 0 0;
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
}

.backup-btn {
  background: linear-gradient(135deg, #10b981, #059669) !important;
  border: none !important;
  font-weight: 600;
}

.backup-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(16, 185, 129, 0.4);
}

.backup-status {
  color: #10b981;
  font-size: 11px;
  padding: 4px 8px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 6px;
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.refresh-btn .n-button__icon {
  background: transparent !important;
}

.refresh-btn .n-icon {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.refresh-btn .n-icon svg {
  background: transparent !important;
  fill: currentColor !important;
}

.refresh-card {
  margin-bottom: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  backdrop-filter: blur(20px);
  position: relative;
  padding: 10px 16px;
  box-sizing: border-box;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.refresh-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.refresh-toggle {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  font-weight: 500;
}

.refresh-toggle :deep(.n-switch) {
  --n-rail-color: rgba(255, 255, 255, 0.1) !important;
  --n-rail-width: 40px !important;
  --n-rail-height: 20px !important;
  --n-button-width: 16px !important;
  --n-button-height: 16px !important;
}

.refresh-toggle :deep(.n-switch.n-switch--checked) {
  --n-rail-color: rgba(99, 102, 241, 0.6) !important;
}

.refresh-toggle :deep(.n-switch--checked .n-switch__button) {
  transform: translateX(22px) !important;
}

.refresh-interval {
  display: flex;
  align-items: center;
  gap: 10px;
}

.refresh-interval :deep(.n-button-group) {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px !important;
  overflow: hidden;
}

.refresh-interval :deep(.n-button) {
  background: transparent !important;
  border: none !important;
  color: rgba(255, 255, 255, 0.6) !important;
  font-size: 12px !important;
  padding: 6px 14px !important;
  transition: all 0.2s ease !important;
}

.refresh-interval :deep(.n-button:hover) {
  color: #fff !important;
  background: rgba(99, 102, 241, 0.2) !important;
}

.refresh-interval :deep(.n-button--type-primary) {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  color: #fff !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

.last-refresh {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
  font-weight: 400;
}

.table-card {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-bottom: 2px solid rgba(99, 102, 241, 0.3);
  backdrop-filter: blur(20px);
  position: relative;
  overflow: hidden;
  padding: 16px;
  box-sizing: border-box;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-weight: 600;
  font-size: 16px;
  color: #fff;
  gap: 12px;
}

.table-content {
  overflow-x: auto;
  overflow-y: auto;
  max-height: calc(100vh - 350px - 70px - env(safe-area-inset-bottom));
  background: transparent !important;
}

.table-content::-webkit-scrollbar {
  display: none;
}

.change-up {
  color: #ff6b6b;
}

.change-down {
  color: #38ef7d;
}

@media (max-width: 768px) {
  .watchlist {
    padding: 12px 12px;
  }

  .watchlist-header {
    flex-direction: column;
    gap: 12px;
    align-items: center;
    text-align: center;
    margin-bottom: 12px;
  }

  .header-left {
    flex-direction: column;
  }

  .refresh-controls {
    flex-wrap: wrap;
    justify-content: center;
  }
}

.refresh-btn :deep(.n-button__icon) {
  background: transparent !important;
}

.refresh-btn :deep(.n-icon) {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.table-card :deep(.n-button__icon) {
  background: transparent !important;
}

.table-card :deep(.n-icon) {
  background: transparent !important;
}

.watchlist-cards {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.watchlist-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.watchlist-card:active {
  background: rgba(99, 102, 241, 0.1);
  transform: scale(0.98);
}

.card-left {
  flex: 1;
  min-width: 0;
}

.card-name {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-code {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 2px;
}

.card-right {
  text-align: right;
  min-width: 80px;
}

.card-price {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.card-change {
  font-size: 12px;
  font-weight: 500;
  margin-top: 2px;
}

.card-change.up {
  color: #ff6b6b;
}

.card-change.down {
  color: #38ef7d;
}

.card-delete-btn {
  flex-shrink: 0;
  opacity: 0.6;
}

.card-delete-btn:hover {
  opacity: 1;
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
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
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

@media (max-width: 480px) {
  .watchlist {
    padding: 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }

  .watchlist-header {
    flex-direction: column;
    gap: 12px;
    align-items: center;
    text-align: center;
    margin-bottom: 12px;
  }

  .header-left {
    flex-direction: column;
  }

  .refresh-card {
    flex-direction: column;
    gap: 10px;
    padding: 10px 12px;
  }

  .refresh-left {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
}
</style>
