<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NIcon,
  NSwitch, NSpin, NEmpty, NButtonGroup
} from 'naive-ui'
import { watchlistApi, stockApi, backupApi, type WatchlistItem, type StockQuote } from '../api'
import { IconRefresh, IconDelete, IconHome, IconPlus, IconBackup, IconCloudDownload } from '../components/icons'
import StockSearch from '../components/StockSearch.vue'
import BackgroundOrbs from '../components/BackgroundOrbs.vue'
import BottomTabs from '../components/BottomTabs.vue'
import { formatVolume } from '../utils/format'

const router = useRouter()

const watchlist = ref<WatchlistItem[]>([])
const quotes = ref<Map<string, StockQuote>>(new Map())
const loading = ref(false)
const refreshing = ref(false)
const backingUp = ref(false)
const backingUpHistory = ref(false)
const backupStatus = ref('')
const historyStatus = ref('')

const autoRefresh = ref(true)
const refreshInterval = ref(30)
const lastRefresh = ref<Date | null>(null)

const showSearch = ref(false)

const bottomTabs = computed(() => [
  { icon: IconHome, label: 'Portfolio', action: () => router.push('/') },
  { icon: IconPlus, label: 'Add', action: () => { showSearch.value = true }, primary: true }
])

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

const handleFetchHistory = async () => {
  if (backingUpHistory.value) return
  backingUpHistory.value = true
  historyStatus.value = 'Fetching historical data...'
  try {
    const result = await watchlistApi.fetchHistory()
    historyStatus.value = `Fetched ${result.newRecords} records`
    setTimeout(() => {
      historyStatus.value = ''
    }, 5000)
  } catch (error) {
    console.error(error)
    historyStatus.value = 'Fetch failed'
    setTimeout(() => {
      historyStatus.value = ''
    }, 3000)
  } finally {
    backingUpHistory.value = false
  }
}

const stopAutoRefresh = () => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }
}

const isMarketOpen = (): boolean => {
  const now = new Date()
  const day = now.getDay()
  if (day === 0 || day === 6) return false

  const shanghaiOffset = 8 * 60 * 60 * 1000
  const shanghaiTime = new Date(now.getTime() + (now.getTimezoneOffset() * 60 * 1000) + shanghaiOffset)
  const hours = shanghaiTime.getHours()
  const minutes = shanghaiTime.getMinutes()
  const totalMinutes = hours * 60 + minutes

  const morningStart = 9 * 60 + 30
  const morningEnd = 11 * 60 + 30
  const afternoonStart = 13 * 60
  const afternoonEnd = 15 * 60

  return (totalMinutes >= morningStart && totalMinutes < morningEnd) ||
         (totalMinutes >= afternoonStart && totalMinutes < afternoonEnd)
}

const refreshTimer = ref<ReturnType<typeof setInterval> | null>(null)

const startAutoRefresh = () => {
  if (refreshTimer.value) clearInterval(refreshTimer.value)
  refreshTimer.value = setInterval(() => {
    if (isMarketOpen()) fetchQuotes()
  }, refreshInterval.value * 1000)
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
    <BackgroundOrbs />

    <div class="watchlist-header">
      <div class="header-left">
        <div class="header-text">
          <h1>Watchlist</h1>
          <p class="subtitle">Track your favorite stocks</p>
        </div>
      </div>
      <div class="header-right">
        <n-button @click="handleFetchHistory" class="history-btn" :class="{ 'is-backing-up': backingUpHistory }" title="Fetch 180 days historical data">
          <template #icon>
            <n-icon>
              <IconCloudDownload />
            </n-icon>
          </template>
        </n-button>
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
      <div class="history-status" v-if="historyStatus">
        {{ historyStatus }}
      </div>
    </div>

    <div class="table-card">
      <div class="table-header">
        <span>Watchlist</span>
        <span class="last-refresh" v-if="lastRefresh">
          {{ lastRefresh.toLocaleTimeString() }}
        </span>
      </div>
      <div class="table-container">
        <div class="table-scroll-wrapper">
          <div class="table-header-row" v-if="tableData.length > 0">
            <div class="header-stock-cell">Stock</div>
            <div class="header-scroll-part">
              <div class="header-price-cell">Price</div>
              <div class="header-change-cell">Change</div>
              <div class="header-hilow-cell">High / Low</div>
              <div class="header-volume-cell">Volume</div>
              <div class="header-action-cell"></div>
            </div>
          </div>
          <div
            v-for="row in tableData"
            :key="row.code"
            class="table-row"
            @click="router.push({ path: `/analysis/${row.code}`, query: { from: '/watchlist' } })"
          >
            <div class="stock-cell">
              <div class="stock-name">{{ row.name }}</div>
              <div class="stock-code">{{ row.code }}</div>
            </div>
            <div class="table-scroll-part">
              <div class="price-cell">
                <span v-if="row.quote" class="price-value">¥{{ row.quote.current.toFixed(2) }}</span>
                <span v-else class="price-value">-</span>
              </div>
              <div class="change-cell">
                <span v-if="row.quote" class="change-value" :class="row.quote.current >= row.quote.prevClose ? 'up' : 'down'">
                  {{ row.quote.current >= row.quote.prevClose ? '+' : '' }}{{ (((row.quote.current - row.quote.prevClose) / row.quote.prevClose) * 100).toFixed(2) }}%
                </span>
                <span v-else class="change-value">-</span>
              </div>
              <div class="hilow-cell">
                <div v-if="row.quote" class="hilow">
                  <span class="hi">H: ¥{{ row.quote.high.toFixed(2) }}</span>
                  <span class="lo">L: ¥{{ row.quote.low.toFixed(2) }}</span>
                </div>
                <span v-else class="price-value">-</span>
              </div>
              <div class="volume-cell">
                <span v-if="row.quote" class="vol-value">{{ formatVolume(row.quote.volume) }}</span>
                <span v-else class="vol-value">-</span>
              </div>
              <div class="action-cell">
                <n-button text @click.stop="handleRemove(row.code)" class="delete-btn">
                  <template #icon>
                    <n-icon><IconDelete /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </div>
          <n-empty v-if="!loading && tableData.length === 0" description="No stocks in watchlist" />
          <div v-if="loading" class="loading-overlay">
            <n-spin :show="loading" description="Loading..." />
          </div>
        </div>
      </div>
    </div>

    <StockSearch
      v-model:show="showSearch"
      @select="handleAddToWatchlist"
    />

    <BottomTabs :tabs="bottomTabs" />
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
  display: flex;
  flex-direction: column;
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

.history-btn {
  background: rgba(99, 102, 241, 0.15) !important;
  border: 1px solid rgba(99, 102, 241, 0.3) !important;
  transition: all 0.2s ease;
}

.history-btn:hover:not(.is-backing-up) {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border-color: transparent !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
}

.history-btn.is-backing-up {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border-color: transparent !important;
  transform: none;
  box-shadow: none;
}

.history-btn.is-backing-up .n-icon {
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

.history-status {
  color: #6366f1;
  font-size: 11px;
  padding: 4px 8px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
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
  backdrop-filter: blur(20px);
  position: relative;
  padding: 0;
  box-sizing: border-box;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.table-scroll-wrapper {
  flex: 1;
  min-height: 0;
  overflow-x: auto;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-x pan-y;
  overscroll-behavior: contain;
}

.table-scroll-wrapper::-webkit-scrollbar {
  display: none;
}

.table-container {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  width: 100%;
  overflow: hidden;
}

.table-header-row {
  display: flex;
  height: 40px;
  background: rgba(10, 10, 15, 0.95);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-stock-cell {
  width: 100px;
  min-width: 100px;
  display: flex;
  align-items: center;
  font-weight: 600;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  padding: 0 12px;
  box-sizing: border-box;
  background: rgba(10, 10, 15, 1);
  position: sticky;
  left: 0;
  z-index: 11;
}

.header-scroll-part {
  display: flex;
  flex: 1;
}

.header-price-cell,
.header-change-cell,
.header-hilow-cell,
.header-volume-cell {
  width: 110px;
  min-width: 110px;
  display: flex;
  align-items: center;
  font-weight: 600;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(10, 10, 15, 1);
  padding: 0 8px;
  box-sizing: border-box;
}

.header-action-cell {
  width: 60px;
  min-width: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(10, 10, 15, 1);
}

.table-row {
  display: flex;
  height: 60px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: background 0.2s ease;
  width: max-content;
  min-width: 100%;
}

.table-row:hover {
  background: rgba(99, 102, 241, 0.08);
}

.stock-cell {
  width: 100px;
  min-width: 100px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 12px;
  gap: 4px;
  line-height: 1.3;
  box-sizing: border-box;
  background: rgba(10, 10, 15, 1);
  position: sticky;
  left: 0;
  z-index: 1;
}

.table-scroll-part {
  display: flex;
  flex: 1;
}

.price-cell,
.change-cell,
.hilow-cell,
.volume-cell {
  width: 110px;
  min-width: 110px;
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #fff;
  background: rgba(10, 10, 15, 1);
  padding: 0 8px;
  box-sizing: border-box;
}

.action-cell {
  width: 60px;
  min-width: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(10, 10, 15, 1);
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(30, 27, 75, 0.8);
  z-index: 20;
}

.stock-name {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stock-name:hover {
  color: #6366f1;
}

.stock-code {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.price-value {
  font-weight: 500;
}

.change-value {
  font-weight: 600;
}

.hilow {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 11px;
}

.hilow .hi {
  color: rgba(255, 255, 255, 0.6);
}

.hilow .lo {
  color: rgba(255, 255, 255, 0.6);
}

.vol-value {
  color: rgba(255, 255, 255, 0.6);
}

.up { color: #ff6b6b; }
.down { color: #38ef7d; }

.delete-btn {
  color: rgba(255, 107, 107, 0.6) !important;
}

.delete-btn:hover {
  color: #ff6b6b !important;
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

  .watchlist-table th,
  .watchlist-table td {
    padding: 10px 8px;
    font-size: 12px;
  }

  .stock-name {
    font-size: 13px;
  }

  .stock-code {
    font-size: 10px;
  }

  .hilow {
    font-size: 10px;
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
