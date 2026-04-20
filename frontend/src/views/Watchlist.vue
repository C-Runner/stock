<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NCard, NDataTable, NSpace, NIcon,
  NSwitch, NSpin, NEmpty, NButtonGroup
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { watchlistApi, stockApi, type WatchlistItem, type StockQuote } from '../api'
import { IconSearch, IconRefresh } from '../components/icons'
import { formatVolume } from '../utils/format'
import StockSearch from '../components/StockSearch.vue'

const router = useRouter()

const watchlist = ref<WatchlistItem[]>([])
const quotes = ref<Map<string, StockQuote>>(new Map())
const loading = ref(false)
const refreshing = ref(false)

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
    <div class="watchlist-header">
      <div class="header-left">
        <h1>Watchlist</h1>
        <p class="subtitle">Track your favorite stocks</p>
      </div>
      <n-space>
        <n-button type="primary" @click="showSearch = true">
          <template #icon>
            <n-icon><IconSearch /></n-icon>
          </template>
          Add Stock
        </n-button>
        <n-button @click="fetchQuotes" :loading="refreshing">
          <template #icon>
            <n-icon><IconRefresh /></n-icon>
          </template>
        </n-button>
      </n-space>
    </div>

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
