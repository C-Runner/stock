<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NCard, NDataTable, NSpace, NIcon,
  NSwitch, NSpin, NEmpty, NButtonGroup
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { watchlistApi, stockApi, type WatchlistItem, type StockQuote } from '../api'
import StockSearch from '../components/StockSearch.vue'

const router = useRouter()

// Icons
const SearchIcon = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z' })
])
const RefreshIcon = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z' })
])

// Watchlist state
const watchlist = ref<WatchlistItem[]>([])
const quotes = ref<Map<string, StockQuote>>(new Map())
const loading = ref(false)
const refreshing = ref(false)

// Auto-refresh state
const autoRefresh = ref(true)
const refreshInterval = ref(30)
const lastRefresh = ref<Date | null>(null)
const refreshTimer = ref<ReturnType<typeof setInterval> | null>(null)

// Search modal
const showSearch = ref(false)

// Table columns
interface WatchlistRow {
  code: string
  name: string
  addedAt: string
  quote: StockQuote | null
}

const columns: DataTableColumns<WatchlistRow> = [
  { title: 'Code', key: 'code', width: 100 },
  { title: 'Name', key: 'name', ellipsis: { tooltip: true } },
  {
    title: 'Price',
    key: 'price',
    width: 120,
    render: (row) => row.quote ? `¥${row.quote.current.toFixed(2)}` : '-'
  },
  {
    title: 'Change',
    key: 'change',
    width: 100,
    render: (row) => {
      if (!row.quote) return '-'
      const rate = ((row.quote.current - row.quote.open) / row.quote.open) * 100
      return h('span', { class: rate >= 0 ? 'change-up' : 'change-down' },
        `${rate >= 0 ? '+' : ''}${rate.toFixed(2)}%`
      )
    }
  },
  {
    title: 'Volume',
    key: 'volume',
    width: 120,
    render: (row) => row.quote ? formatVolume(row.quote.volume) : '-'
  },
  {
    title: 'Action',
    key: 'actions',
    width: 200,
    render: (row) => h(NSpace, { size: 'small' }, () => [
      h(NButton, {
        size: 'small',
        type: 'primary',
        quaternary: true,
        onClick: () => router.push(`/analysis/${row.code}`)
      }, () => 'Analysis'),
      h(NButton, {
        size: 'small',
        type: 'error',
        quaternary: true,
        onClick: () => handleRemove(row.code)
      }, () => 'Remove')
    ])
  }
]

// Table data with quotes
const tableData = computed<WatchlistRow[]>(() =>
  watchlist.value.map(item => ({
    ...item,
    quote: quotes.value.get(item.code) || null
  }))
)

const formatVolume = (vol: number): string => {
  if (vol >= 100000000) return (vol / 100000000).toFixed(2) + ' 亿'
  if (vol >= 10000) return (vol / 10000).toFixed(2) + ' 万'
  return vol.toLocaleString()
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
    const quotePromises = watchlist.value.map(item =>
      stockApi.getQuote(item.code).catch(() => null)
    )
    const results = await Promise.all(quotePromises)

    quotes.value.clear()
    results.forEach((quote, index) => {
      if (quote && watchlist.value[index]) {
        quotes.value.set(watchlist.value[index].code, quote)
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
  if (enabled) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
})

watch(refreshInterval, () => {
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

onMounted(async () => {
  await fetchWatchlist()
  await fetchQuotes()
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<template>
  <div class="watchlist">
    <!-- Header -->
    <div class="watchlist-header">
      <div class="header-left">
        <h1>Watchlist</h1>
        <p class="subtitle">Track your favorite stocks</p>
      </div>
      <n-space>
        <n-button type="primary" @click="showSearch = true">
          <template #icon>
            <n-icon><SearchIcon /></n-icon>
          </template>
          Add Stock
        </n-button>
        <n-button @click="fetchQuotes" :loading="refreshing">
          <template #icon>
            <n-icon><RefreshIcon /></n-icon>
          </template>
        </n-button>
      </n-space>
    </div>

    <!-- Refresh Controls -->
    <n-card class="refresh-card" :bordered="false">
      <div class="refresh-controls">
        <div class="refresh-toggle">
          <span>Auto Refresh</span>
          <n-switch v-model:value="autoRefresh" />
        </div>
        <div class="refresh-interval" v-if="autoRefresh">
          <span>Interval:</span>
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
        <div class="last-refresh" v-if="lastRefresh">
          Last refresh: {{ lastRefresh.toLocaleTimeString() }}
        </div>
      </div>
    </n-card>

    <!-- Watchlist Table -->
    <n-card class="table-card" :bordered="false">
      <n-spin :show="loading">
        <n-empty v-if="!loading && tableData.length === 0" description="No stocks in watchlist" />
        <n-data-table
          v-else
          :columns="columns"
          :data="tableData"
          :pagination="false"
          :bordered="false"
          striped
        />
      </n-spin>
    </n-card>

    <!-- Search Modal -->
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
  background: #000;
  min-height: calc(100vh - 60px);
  padding: 20px;
  box-sizing: border-box;
}

.watchlist-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #fff;
}

.subtitle {
  margin: 4px 0 0;
  color: #999;
  font-size: 14px;
}

.refresh-card {
  margin-bottom: 16px;
  background: #1a1a1a !important;
  border-radius: 12px;
}

.refresh-card :deep(.n-card) {
  background: #1a1a1a !important;
  border: 1px solid #333 !important;
}

.refresh-controls {
  display: flex;
  align-items: center;
  gap: 24px;
}

.refresh-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #999;
}

.refresh-interval {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #999;
}

.last-refresh {
  color: #666;
  font-size: 12px;
  margin-left: auto;
}

.table-card {
  background: #1a1a1a !important;
  border-radius: 12px;
}

.table-card :deep(.n-card) {
  background: #1a1a1a !important;
  border: 1px solid #333 !important;
}

.change-up {
  color: #ef5350;
}

.change-down {
  color: #26a69a;
}

@media (max-width: 768px) {
  .watchlist {
    padding: 12px;
  }

  .watchlist-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .refresh-controls {
    flex-wrap: wrap;
  }
}
</style>
