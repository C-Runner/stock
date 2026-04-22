<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, h, watch } from 'vue'
import {
  NButton, NDataTable, NSpace, NIcon,
  NSwitch, NSpin, NEmpty, NButtonGroup
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { watchlistApi, stockApi, backupApi, type WatchlistItem, type StockQuote } from '../api'
import { IconSearch, IconRefresh, IconBackup } from '../components/icons'
import { formatVolume } from '../utils/format'
import StockSearch from '../components/StockSearch.vue'


const watchlist = ref<WatchlistItem[]>([])
const quotes = ref<Map<string, StockQuote>>(new Map())
const loading = ref(false)
const refreshing = ref(false)
const backingUp = ref(false)
const backupStatus = ref<string>('')

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

const columns: DataTableColumns<WatchlistRow> = [
  {
    title: 'Stock',
    key: 'stock',
    width: 90,
    render: (row) => h('div', { class: 'stock-cell' }, [
      h('div', {
        class: 'stock-name',
        onClick: (e: MouseEvent) => {
          e.stopPropagation()
          window.location.href = `/analysis/${row.code}`
        }
      }, row.name),
      h('div', { class: 'stock-code' }, row.code)
    ])
  },
  {
    title: 'Price',
    key: 'price',
    width: 90,
    render: (row) => row.quote ? `¥${row.quote.current.toFixed(2)}` : '-'
  },
  {
    title: 'Change',
    key: 'change',
    width: 90,
    render: (row) => {
      if (!row.quote) return '-'
      const rate = row.quote.prevClose > 0
        ? ((row.quote.current - row.quote.prevClose) / row.quote.prevClose) * 100
        : 0
      return h('span', { class: rate >= 0 ? 'change-up' : 'change-down' },
        `${rate >= 0 ? '+' : ''}${rate.toFixed(2)}%`
      )
    }
  },
  {
    title: 'Volume',
    key: 'volume',
    width: 90,
    render: (row) => row.quote ? formatVolume(row.quote.volume) : '-'
  },
  {
    title: 'Action',
    key: 'actions',
    width: 80,
    render: (row) => h(NButton, {
      size: 'small',
      type: 'error',
      quaternary: true,
      onClick: (e: Event) => { e.stopPropagation(); handleRemove(row.code) }
    }, () => 'Remove')
  }
]

const tableData = computed<WatchlistRow[]>(() =>
  watchlist.value.map(item => ({
    ...item,
    quote: quotes.value.get(item.code) || null
  }))
)

const rowProps = (row: WatchlistRow) => {
  const goToAnalysis = () => {
    window.location.href = `/analysis/${row.code}`
  }
  return {
    style: 'cursor: pointer',
    onClick: goToAnalysis
  }
}

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
    // Poll for completion or just show started status
    setTimeout(() => {
      backupStatus.value = ''
    }, 3000)
  } catch (error) {
    console.error(error)
    backupStatus.value = 'Backup failed'
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
      <n-space>
        <n-button type="primary" @click="handleBackup" :loading="backingUp" class="backup-btn">
          <template #icon>
            <n-icon><IconBackup /></n-icon>
          </template>
          Backup
        </n-button>
        <n-button type="primary" @click="showSearch = true" class="add-btn">
          <template #icon>
            <n-icon><IconSearch /></n-icon>
          </template>
          Add Stock
        </n-button>
        <n-button @click="fetchQuotes" :loading="refreshing" class="refresh-btn" style="background: rgba(255, 255, 255, 0.05) !important;">
          <template #icon>
            <n-icon style="background: transparent !important;"><IconRefresh /></n-icon>
          </template>
        </n-button>
      </n-space>
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
      <div class="refresh-right">
        <div class="backup-status" v-if="backupStatus">
          {{ backupStatus }}
        </div>
      </div>
    </div>

    <div class="table-card">
      <div class="table-header">
        <span>Watchlist</span>
        <span class="last-refresh" v-if="lastRefresh">
          Last: {{ lastRefresh.toLocaleTimeString() }}
        </span>
      </div>
      <div class="table-content">
        <n-spin :show="loading">
          <n-empty v-if="!loading && tableData.length === 0" description="No stocks in watchlist" />
          <n-data-table
            v-else
            :columns="columns"
            :data="tableData"
            :pagination="false"
            :bordered="false"
            striped
            :row-props="rowProps"
          />
        </n-spin>
      </div>
    </div>

    <StockSearch
      v-model:show="showSearch"
      @select="handleAddToWatchlist"
    />
  </div>
</template>

<style scoped>
.watchlist {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  min-height: calc(100vh - 60px);
  padding: 16px 24px;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
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

.add-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border: none !important;
  font-weight: 600;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
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

.refresh-right {
  display: flex;
  align-items: center;
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
  background: transparent !important;
}

.table-card :deep(.n-data-table) {
  font-size: 14px;
  background: transparent !important;
}

.table-card :deep(.n-data-table-wrapper) {
  background: transparent !important;
}

.table-card :deep(.n-data-table-th) {
  background: transparent !important;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6) !important;
  padding: 6px 6px !important;
}

.table-card :deep(.n-data-table-td) {
  background: transparent !important;
  color: #fff;
  padding: 6px 6px !important;
}

.table-card :deep(.n-data-table-tr) {
  background: transparent !important;
}

.table-card :deep(.n-base-table) {
  background: transparent !important;
}

.table-card :deep(.n-base-table-tbody) {
  background: transparent !important;
}

.table-card :deep(.n-base-table-tr) {
  background: transparent !important;
}

.table-card :deep(.n-base-td) {
  background: transparent !important;
}

.table-card :deep(.n-data-table-td__cell) {
  white-space: normal !important;
}

.table-card :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: rgba(99, 102, 241, 0.08) !important;
}

.table-card :deep(.n-data-table-tr) {
  cursor: pointer;
}

.change-up {
  color: #ff6b6b;
}

.change-down {
  color: #38ef7d;
}

.code-link {
  color: #6366f1;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.2s;
}

.code-link:hover {
  color: #818cf8;
  text-decoration: underline;
}

.stock-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  line-height: 1.3;
}

.stock-name {
  font-size: 14px;
  color: #fff;
  cursor: pointer;
}

.stock-name:hover {
  color: #6366f1;
}

.stock-code {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

@media (max-width: 768px) {
  .table-card :deep(.n-data-table) {
    font-size: 12px;
  }

  .table-card :deep(.n-data-table-th),
  .table-card :deep(.n-data-table-td) {
    padding: 6px 4px;
  }

  .table-card :deep(.n-data-table-wrapper) {
    background: transparent !important;
  }

  .table-card :deep(.n-data-table-tr) {
    background: transparent !important;
  }

  .table-card :deep(.n-data-table-td) {
    background: transparent !important;
  }

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
</style>
