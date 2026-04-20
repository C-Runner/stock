<template>
  <div class="home">
    <div class="home-header">
      <div class="header-left">
        <h1>Stock Portfolio</h1>
        <p class="subtitle">Manage your stock holdings</p>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true">
        <template #icon>
          <n-icon><IconPlus /></n-icon>
        </template>
        Add Stock
      </n-button>
    </div>

    <div class="stats-row">
      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon blue">
            <n-icon size="24"><IconWallet /></n-icon>
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
            <n-icon size="24"><IconCoin /></n-icon>
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
            <n-icon size="24"><IconTrend /></n-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">Total Cost</span>
            <span class="stat-value">¥{{ totalCost.toFixed(2) }}</span>
          </div>
        </div>
      </n-card>

      <n-card class="stat-card" :bordered="false">
        <div class="stat-content">
          <div class="stat-icon" :class="totalProfit >= 0 ? 'red' : 'teal'">
            <n-icon size="24"><IconDataLine /></n-icon>
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

    <n-card class="table-card" :bordered="false">
      <template #header>
        <div class="table-header">
          <span>Holdings</span>
          <n-button text @click="fetchStocks" :loading="loading">
            <template #icon>
              <n-icon><IconRefresh /></n-icon>
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
import { IconPlus, IconWallet, IconCoin, IconTrend, IconDataLine, IconRefresh } from '../components/icons'

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
.stat-icon.teal { background: linear-gradient(135deg, #38ef7d 0%, #11998e 100%); }

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
