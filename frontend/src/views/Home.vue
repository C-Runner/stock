<template>
  <div class="home">
    <!-- Header -->
    <div class="home-header">
      <div class="header-left">
        <h1>Stock Portfolio</h1>
        <p class="subtitle">Manage your stock holdings</p>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true">
        <template #icon>
          <n-icon><Plus /></n-icon>
        </template>
        Add Stock
      </n-button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon blue">
            <n-icon size="24"><Wallet /></n-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">Total Stocks</span>
            <span class="stat-value">{{ stocks.length }}</span>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon green">
            <n-icon size="24"><Coin /></n-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">Total Market Value</span>
            <span class="stat-value">¥{{ totalMarketValue.toFixed(2) }}</span>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon orange">
            <n-icon size="24"><TrendCharts /></n-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">Total Cost</span>
            <span class="stat-value">¥{{ totalCost.toFixed(2) }}</span>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon" :class="totalProfit >= 0 ? 'red' : 'green'">
            <n-icon size="24"><DataLine /></n-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">Total P/L</span>
            <span class="stat-value" :class="totalProfit >= 0 ? 'up' : 'down'">
              {{ totalProfit >= 0 ? '+' : '' }}¥{{ totalProfit.toFixed(2) }}
            </span>
          </div>
        </div>
      </n-card>
    </div>

    <!-- Stock Table -->
    <n-card class="table-card" :bordered="false">
      <template #header>
        <div class="table-header">
          <span>Holdings</span>
          <n-button text @click="fetchStocks" :loading="loading">
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
          </n-button>
        </div>
      </template>

      <n-empty v-if="!loading && stocks.length === 0" description="No stocks yet, click 'Add Stock' to add one" />
      <n-data-table
        v-else
        :columns="columns"
        :data="stocks"
        :loading="loading"
        :pagination="false"
        :bordered="false"
        striped
      />
    </n-card>

    <!-- Add Stock Modal -->
    <n-modal v-model:show="showAddModal" preset="card" title="Add New Stock" style="width: 500px;">
      <n-form :model="stockForm">
        <n-form-item label="Stock Code" path="code">
          <n-space vertical :size="12" style="width: 100%">
            <n-input v-model:value="stockForm.code" placeholder="e.g. 600519 or sh600519" @blur="lookupStock" />
            <n-button :loading="lookingUp" @click="lookupStock" :disabled="!stockForm.code" block>
              Search Info
            </n-button>
          </n-space>
        </n-form-item>
        <n-form-item label="Stock Name" path="name">
          <n-input v-model:value="stockForm.name" placeholder="Auto-filled after search" disabled />
        </n-form-item>
        <n-form-item label="Current Price" path="currentPrice">
          <n-input-number
            v-model:value="stockForm.currentPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
            placeholder="Current market price"
          />
        </n-form-item>
        <n-form-item label="Quantity" path="quantity">
          <n-input-number
            v-model:value="stockForm.quantity"
            :min="0"
            style="width: 100%"
            placeholder="Number of shares"
          />
        </n-form-item>
        <n-form-item label="Buy Price" path="buyPrice">
          <n-input-number
            v-model:value="stockForm.buyPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
            placeholder="Average purchase price"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showAddModal = false">Cancel</n-button>
          <n-button type="primary" @click="handleAddStock" :loading="submitting">Add Stock</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NCard, NDataTable, NModal, NForm, NFormItem,
  NInput, NInputNumber, NEmpty, NSpace, NIcon
} from 'naive-ui'
import { stockApi, type Stock, type StockRequest } from '../api'

// Icons
const Plus = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z' })
])
const Wallet = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M21 18v1c0 1.1-.9 2-2 2H5c-1.11 0-2-.9-2-2V5c0-1.1.89-2 2-2h14c1.1 0 2 .9 2 2v1h-9c-1.11 0-2 .9-2 2v8c0 1.1.89 2 2 2h9zm-9-2h10V8H12v8zm4-2.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z' })
])
const Coin = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm.31-8.86c-1.77-.45-2.34-.94-2.34-1.67 0-.84.79-1.43 2.1-1.43 1.38 0 1.9.66 1.94 1.64h1.71c-.05-1.34-.87-2.57-2.49-2.97V5H10.9v1.69c-1.51.32-2.72 1.3-2.72 2.81 0 1.79 1.49 2.69 3.66 3.21 1.95.46 2.34 1.15 2.34 1.87 0 .53-.39 1.39-2.1 1.39-1.6 0-2.23-.72-2.32-1.64H8.04c.1 1.7 1.36 2.66 2.86 2.97V19h2.34v-1.76c1.86-.29 2.81-1.4 2.81-2.82 0-2.39-1.79-2.96-3.66-3.42z' })
])
const TrendCharts = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z' })
])
const DataLine = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z' })
])
const Refresh = () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', fill: 'currentColor' }, [
  h('path', { d: 'M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z' })
])

const router = useRouter()

const stocks = ref<Stock[]>([])
const loading = ref(false)
const showAddModal = ref(false)
const submitting = ref(false)
const lookingUp = ref(false)

const stockForm = ref<StockRequest>({
  code: '',
  name: '',
  currentPrice: 0,
  quantity: 0,
  buyPrice: 0
})

// Computed stats
const totalMarketValue = computed(() =>
  stocks.value.reduce((sum, s) => sum + s.currentPrice * s.quantity, 0)
)

const totalCost = computed(() =>
  stocks.value.reduce((sum, s) => sum + s.buyPrice * s.quantity, 0)
)

const totalProfit = computed(() => totalMarketValue.value - totalCost.value)

const columns = [
  { title: 'Code', key: 'code', width: 120 },
  { title: 'Name', key: 'name', ellipsis: { tooltip: true } },
  {
    title: 'Current',
    key: 'currentPrice',
    width: 120,
    render: (row: Stock) => `¥${row.currentPrice.toFixed(2)}`
  },
  {
    title: 'Qty',
    key: 'quantity',
    width: 80,
    render: (row: Stock) => row.quantity.toLocaleString()
  },
  {
    title: 'Cost',
    key: 'buyPrice',
    width: 100,
    render: (row: Stock) => `¥${row.buyPrice.toFixed(2)}`
  },
  {
    title: 'Value',
    key: 'value',
    width: 130,
    render: (row: Stock) => `¥${(row.currentPrice * row.quantity).toFixed(2)}`
  },
  {
    title: 'P/L',
    key: 'profit',
    width: 130,
    render: (row: Stock) => {
      const profit = (row.currentPrice - row.buyPrice) * row.quantity
      const rate = ((row.currentPrice - row.buyPrice) / row.buyPrice) * 100
      return h('span', { class: profit >= 0 ? 'profit-up' : 'profit-down' },
        `${profit >= 0 ? '+' : ''}¥${profit.toFixed(2)} (${rate.toFixed(2)}%)`
      )
    }
  },
  {
    title: 'Action',
    key: 'actions',
    width: 160,
    render: (row: Stock) => h(NSpace, { size: 'small' }, () => [
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
        onClick: () => handleDeleteStock(row.code)
      }, () => 'Delete')
    ])
  }
]

const lookupStock = async () => {
  if (!stockForm.value.code) return
  lookingUp.value = true
  try {
    const quote = await stockApi.getQuote(stockForm.value.code)
    stockForm.value.name = quote.name
    stockForm.value.currentPrice = quote.current
  } catch {
    stockForm.value.name = ''
  } finally {
    lookingUp.value = false
  }
}

const fetchStocks = async () => {
  loading.value = true
  try {
    stocks.value = await stockApi.getStocks()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
  {
    loading.value }
}

const handleAddStock = async () => {
  if (!stockForm.value.code) return
  if (!stockForm.value.currentPrice || stockForm.value.currentPrice <= 0) return
  if (!stockForm.value.quantity || stockForm.value.quantity <= 0) return
  if (!stockForm.value.buyPrice || stockForm.value.buyPrice <= 0) return

  const formData = { ...stockForm.value }
  if (!formData.name) formData.name = formData.code

  submitting.value = true
  try {
    await stockApi.createStock(formData)
    showAddModal.value = false
    stockForm.value = { code: '', name: '', currentPrice: 0, quantity: 0, buyPrice: 0 }
    await fetchStocks()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDeleteStock = async (code: string) => {
  await stockApi.deleteStock(code)
  await fetchStocks()
}

onMounted(() => fetchStocks())
</script>

<style scoped>
.home {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  background: #000;
  min-height: calc(100vh - 60px);
  box-sizing: border-box;
}

/* Header */
.home-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
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

/* Stats Row */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 12px;
  background: #1a1a1a !important;
}

.stat-card :deep(.n-card) {
  background: #1a1a1a !important;
  border: 1px solid #333 !important;
}

.table-card {
  border-radius: 12px;
  background: #1a1a1a !important;
}

.table-card :deep(.n-card) {
  background: #1a1a1a !important;
  border: 1px solid #333 !important;
}

.stat-card :deep(.n-card__content) {
  padding: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.blue { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.green { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.stat-icon.orange { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.red { background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%); }
.stat-icon.red.up { background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%); }
.stat-icon.green.down { background: linear-gradient(135deg, #38ef7d 0%, #11998e 100%); }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.stat-value.up { color: #ff6b6b; }
.stat-value.down { color: #38ef7d; }

/* Table Card */
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

:deep(.n-data-table) {
  font-size: 14px;
}

:deep(.n-data-table th) {
  background: #1a1a1a !important;
  font-weight: 600;
}

:deep(.profit-up) {
  color: #ff6b6b;
}

:deep(.profit-down) {
  color: #38ef7d;
}

@media (max-width: 1024px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: 1fr;
  }

  .home-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .table-card :deep(.n-data-table) {
    font-size: 12px;
  }

  .table-card :deep(.n-data-table-th),
  .table-card :deep(.n-data-table-td) {
    padding: 8px 4px;
  }
}

@media (max-width: 480px) {
  .home {
    padding: 12px;
    padding-bottom: calc(12px + env(safe-area-inset-bottom));
  }

  .home-header {
    padding: 0 4px;
  }

  .header-left h1 {
    font-size: 20px;
  }

  .subtitle {
    font-size: 12px;
  }

  .stat-card :deep(.n-card__content) {
    padding: 12px;
  }

  .stat-content {
    gap: 10px;
  }

  .stat-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
  }

  .stat-icon :deep(.n-icon) {
    font-size: 18px !important;
  }

  .stat-value {
    font-size: 14px;
  }

  .stat-label {
    font-size: 10px;
  }

  /* Table optimizations */
  .table-card :deep(.n-data-table) {
    font-size: 11px;
  }

  .table-card :deep(.n-data-table-th),
  .table-card :deep(.n-data-table-td) {
    padding: 6px 2px;
  }

  .table-card :deep(.n-data-table-th) {
    font-size: 10px;
  }
}
</style>
