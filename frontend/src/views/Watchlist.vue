<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
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
  if (backingUp.value) return
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

const formatVolume = (vol: number): string => {
  if (vol >= 100000000) return (vol / 100000000).toFixed(2) + 'E'
  if (vol >= 10000) return (vol / 10000).toFixed(2) + 'W'
  if (vol >= 1000) return (vol / 1000).toFixed(2) + 'K'
  return vol.toFixed(0)
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
        <div class="header-text">
          <h1>Watchlist</h1>
          <p class="subtitle">Track your favorite stocks</p>
        </div>
      </div>
      <div class="header-right">
        <n-button @click="handleBackup" class="backup-btn" :class="{ 'is-backing-up': backingUp }">
          <template #icon>
            <n-icon>
              <IconBackup />
            </n-icon>
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
              <div class="card-price-cell" v-if="row.quote">
                <div class="price-current">¥{{ row.quote.current.toFixed(2) }}</div>
              </div>
              <div class="card-price-cell" v-else>
                <div class="price-current">-</div>
              </div>
              <div class="card-change-cell" v-if="row.quote">
                <div class="change-text" :class="row.quote.current >= row.quote.prevClose ? 'up' : 'down'">
                  {{ row.quote.current >= row.quote.prevClose ? '+' : '' }}{{ (((row.quote.current - row.quote.prevClose) / row.quote.prevClose) * 100).toFixed(2) }}%
                </div>
              </div>
              <div class="card-change-cell" v-else>
                <div class="change-text">-</div>
              </div>
              <div class="card-hilow-cell" v-if="row.quote">
                <div class="price-high">
                  <span class="label">H</span>
                  <span class="value up">¥{{ row.quote.high.toFixed(2) }}</span>
                </div>
                <div class="price-low">
                  <span class="label">L</span>
                  <span class="value down">¥{{ row.quote.low.toFixed(2) }}</span>
                </div>
              </div>
              <div class="card-hilow-cell" v-else>
                <div class="price-high">-</div>
                <div class="price-low">-</div>
              </div>
              <div class="card-vol-cell" v-if="row.quote">
                <div class="vol-value">{{ formatVolume(row.quote.volume) }}</div>
              </div>
              <div class="card-vol-cell" v-else>
                <div class="vol-value">-</div>
              </div>
              <div class="card-action">
                <n-button text @click.stop="handleRemove(row.code)" class="delete-btn">
                  <template #icon>
                    <n-icon><IconDelete /></n-icon>
                  </template>
                </n-button>
              </div>
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
  height: 100dvh;
  padding: 16px;
  padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}

.background {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  z-index: -1;
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
  background: rgba(16, 185, 129, 0.15) !important;
  border: 1px solid rgba(16, 185, 129, 0.3) !important;
  transition: all 0.2s ease;
}

.backup-btn:hover:not(.is-backing-up) {
  background: linear-gradient(135deg, #10b981, #059669) !important;
  border-color: transparent !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(16, 185, 129, 0.4);
}

.backup-btn.is-backing-up {
  background: linear-gradient(135deg, #10b981, #059669) !important;
  border-color: transparent !important;
  transform: none;
  box-shadow: none;
}

.backup-btn.is-backing-up .n-icon {
  opacity: 1;
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
  padding: 0;
  box-sizing: border-box;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 12px 16px 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-weight: 600;
  font-size: 16px;
  color: #fff;
  flex: 0 0 auto;
}

.table-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: transparent !important;
}

.watchlist-cards {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-y;
}

.watchlist-cards::-webkit-scrollbar {
  display: none;
}

@media (max-width: 768px) {
  .watchlist {
    padding: 12px 12px;
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

  .refresh-controls {
    flex-wrap: wrap;
    justify-content: center;
  }

  .watchlist-card {
    grid-template-columns: minmax(80px, 1fr) 70px 70px 80px 60px 32px;
    height: 52px;
    padding: 0 8px;
  }

  .card-name {
    font-size: 13px;
  }

  .card-code {
    font-size: 10px;
  }

  .price-current, .change-text {
    font-size: 13px;
  }

  .price-high, .price-low {
    font-size: 10px;
  }

  .vol-value {
    font-size: 11px;
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
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;
}

.watchlist-cards::-webkit-scrollbar {
  display: none;
}

.watchlist-card {
  display: grid;
  grid-template-columns: minmax(100px, 1fr) 80px 80px 90px 70px 32px;
  align-items: center;
  height: 60px;
  padding: 0 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: background 0.2s ease;
  flex-shrink: 0;
  overflow: visible;
}

.watchlist-card:hover {
  background: rgba(99, 102, 241, 0.08);
}

.card-left {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 3px;
  min-width: 0;
  overflow: hidden;
}

.card-name {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-name:hover {
  color: #6366f1;
}

.card-code {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.card-price-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.price-current {
  font-size: 14px;
  color: #fff;
  font-weight: 500;
}

.card-change-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.change-text {
  font-size: 14px;
  font-weight: 600;
}

.card-hilow-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1px;
}

.price-high, .price-low {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
}

.price-high .label, .price-low .label {
  color: rgba(255, 255, 255, 0.3);
  width: 10px;
}

.price-high .value, .price-low .value {
  font-weight: 500;
}

.card-vol-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.vol-value {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.card-action {
  display: flex;
  align-items: center;
  justify-content: center;
}

.up {
  color: #ff6b6b;
}

.down {
  color: #38ef7d;
}

.delete-btn {
  color: rgba(255, 107, 107, 0.6) !important;
}

.delete-btn:hover {
  color: #ff6b6b !important;
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
    padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
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
