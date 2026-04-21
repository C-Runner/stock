<template>
  <div class="home">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>

    <div class="home-header">
      <div class="header-left">
        <div class="logo">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 3v18h18"/>
            <path d="M18 9l-5 5-4-4-3 3"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>Stock Portfolio</h1>
          <p class="subtitle">Manage your stock holdings</p>
        </div>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true" class="add-btn">
        <template #icon>
          <n-icon><IconPlus /></n-icon>
        </template>
        Add Stock
      </n-button>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-grid">
          <div class="stat-item">
            <div class="stat-icon blue">
              <n-icon size="20"><IconWallet /></n-icon>
            </div>
            <div class="stat-info">
              <span class="stat-label">Total Stocks</span>
              <span class="stat-value">{{ stocks.length }}</span>
            </div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-icon green">
              <n-icon size="20"><IconCoin /></n-icon>
            </div>
            <div class="stat-info">
              <span class="stat-label">Market Value</span>
              <span class="stat-value">¥{{ totalMarketValue.toFixed(2) }}</span>
            </div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-icon orange">
              <n-icon size="20"><IconTrend /></n-icon>
            </div>
            <div class="stat-info">
              <span class="stat-label">Total Cost</span>
              <span class="stat-value">¥{{ totalCost.toFixed(2) }}</span>
            </div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-icon" :class="totalProfit >= 0 ? 'red' : 'teal'">
              <n-icon size="20"><IconDataLine /></n-icon>
            </div>
            <div class="stat-info">
              <span class="stat-label">Total P/L</span>
              <span class="stat-value" :class="totalProfit >= 0 ? 'up' : 'down'">
                {{ totalProfit >= 0 ? '+' : '' }}¥{{ totalProfit.toFixed(2) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="table-card">
      <div class="table-header">
        <span>Holdings</span>
        <n-button text @click="fetchStocks" :loading="loading" class="refresh-icon-btn" style="background: transparent !important;">
          <template #icon>
            <n-icon style="background: transparent !important;"><IconRefresh /></n-icon>
          </template>
        </n-button>
      </div>

      <div class="table-content">
        <n-empty v-if="!loading && stocks.length === 0" description="No stocks yet, click 'Add Stock' to add one" />
        <n-data-table
          v-else
          :columns="columns"
          :data="stocks"
          :loading="loading"
          :pagination="false"
          :bordered="false"
        />
      </div>
    </div>

    <n-modal
      v-model:show="showAddModal"
      preset="card"
      class="add-stock-modal"
      :style="{
        '--n-color': 'rgba(20, 19, 60, 0.95)',
        '--n-color-modal': 'rgba(20, 19, 60, 0.95)',
        background: 'rgba(20, 19, 60, 0.95)',
        backdropFilter: 'blur(24px)'
      }"
      :mask="{ style: 'background: rgba(0, 0, 0, 0.7); backdrop-filter: blur(4px);' }"
      :icon="() => null"
      :mask-closable="false"
    >
      <template #header>
        <div class="modal-header">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="modal-icon">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          <span>Add New Stock</span>
        </div>
      </template>
      <div class="modal-body">
        <n-form :model="stockForm" class="add-stock-form">
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
      </div>
      <template #action>
        <div class="modal-footer">
          <n-button @click="showAddModal = false" class="cancel-btn">Cancel</n-button>
          <n-button type="primary" @click="handleAddStock" :loading="submitting" class="submit-btn">Add Stock</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NDataTable, NModal, NForm, NFormItem,
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
  { title: 'Code', key: 'code', width: 120, render: (row: Stock) => row.name },
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
    // Also fetch real-time quotes to update current prices
    if (stocks.value.length > 0) {
      const updatedStocks = await Promise.all(
        stocks.value.map(async (stock) => {
          try {
            const quote = await stockApi.getQuote(stock.code)
            return { ...stock, currentPrice: quote.current }
          } catch {
            return stock
          }
        })
      )
      stocks.value = updatedStocks
    }
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
  min-height: calc(100vh - 60px);
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
  padding: 40px 24px;
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
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: #8b5cf6;
  bottom: -100px;
  right: -100px;
  animation: float 10s ease-in-out infinite reverse;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: #06b6d4;
  top: 50%;
  left: 50%;
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

.home-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  position: relative;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo {
  width: 48px;
  height: 48px;
  padding: 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 14px;
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
  font-size: 26px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.5px;
}

.subtitle {
  margin: 4px 0 0;
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
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

.stats-row {
  margin-bottom: 32px;
  position: relative;
}

.stat-card {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  padding: 24px;
  box-sizing: border-box;
}

.stat-grid {
  display: flex;
  align-items: stretch;
  justify-content: space-between;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  padding: 8px 0;
}

.stat-divider {
  width: 1px;
  align-self: stretch;
  margin: 8px 0;
  background: rgba(255, 255, 255, 0.1);
}

.table-card {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  position: relative;
  overflow: hidden;
  padding: 20px;
  box-sizing: border-box;
}

.table-content {
  overflow-x: auto;
  background: transparent !important;
}

.table-content :deep(*) {
  background: transparent !important;
}

.stat-icon {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.blue { background: linear-gradient(135deg, #6366f1, #8b5cf6); }
.stat-icon.green { background: linear-gradient(135deg, #11998e, #38ef7d); }
.stat-icon.orange { background: linear-gradient(135deg, #f093fb, #f5576c); }
.stat-icon.red { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.stat-icon.teal { background: linear-gradient(135deg, #38ef7d, #11998e); }

.stat-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.stat-value {
  font-size: 22px;
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
  color: #fff;
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
}

.table-card :deep(.n-data-table-td) {
  background: transparent !important;
  color: #fff;
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

.table-card :deep(.n-data-table tr:hover .n-data-table-td) {
  background: rgba(99, 102, 241, 0.08) !important;
}

.table-card :deep(.profit-up) {
  color: #ff6b6b;
}

.table-card :deep(.profit-down) {
  color: #38ef7d;
}

@media (max-width: 1024px) {
  .stat-grid {
    flex-wrap: wrap;
  }

  .stat-item {
    flex: 1 1 calc(50% - 8px);
    min-width: 140px;
  }

  .stat-divider {
    display: none;
  }
}

@media (max-width: 768px) {
  .home {
    padding: 24px 16px;
  }

  .home-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .header-left {
    flex-direction: column;
  }

  .table-card :deep(.n-data-table) {
    font-size: 12px;
  }

  .table-card :deep(.n-data-table-th),
  .table-card :deep(.n-data-table-td) {
    padding: 8px 4px;
  }
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

@media (max-width: 480px) {
  .home {
    padding: 16px 12px;
    padding-bottom: calc(16px + env(safe-area-inset-bottom));
  }

  .header-text h1 {
    font-size: 20px;
  }

  .subtitle {
    font-size: 12px;
  }

  .stat-card {
    padding: 16px;
  }

  .stat-item {
    flex: 1 1 calc(50% - 8px);
  }

  .stat-icon {
    width: 32px;
    height: 32px;
    border-radius: 8px;
  }

  .stat-icon :deep(.n-icon) {
    font-size: 16px !important;
  }

  .stat-value {
    font-size: 16px;
  }

  .stat-label {
    font-size: 10px;
  }
}

.add-stock-modal {
  --n-color: rgba(20, 19, 60, 0.95) !important;
  --n-border-radius: 24px !important;
  --n-border: 1px solid rgba(99, 102, 241, 0.2) !important;
  background-color: var(--n-color) !important;
}

.add-stock-modal :deep(.n-card) {
  background-color: rgba(20, 19, 60, 0.95) !important;
}

.add-stock-modal :deep(.n-dialog__close) {
  top: 20px;
  right: 20px;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.6) !important;
  transition: all 0.2s ease;
}

.add-stock-modal :deep(.n-dialog__close:hover) {
  background: rgba(255, 107, 107, 0.15);
  border-color: rgba(255, 107, 107, 0.3);
  color: #ff6b6b !important;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  font-size: 18px;
  color: #fff;
  margin-bottom: 0;
}

.modal-icon {
  width: 28px;
  height: 28px;
  padding: 6px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 8px;
  color: #fff;
  flex-shrink: 0;
}

.modal-body,
.add-stock-modal :deep(.n-card),
.add-stock-modal :deep(.n-card__content),
.add-stock-modal :deep(.n-card .n-card__content),
.add-stock-modal :deep(.n-base-card),
.add-stock-modal :deep(.n-base-card .n-card__content) {
  padding-top: 24px;
  background: rgba(20, 19, 60, 0.95) !important;
  background-color: rgba(20, 19, 60, 0.95) !important;
  backdrop-filter: blur(24px) !important;
}

.add-stock-form :deep(.n-form-item) {
  margin-bottom: 20px;
}

.add-stock-form :deep(.n-form-item-label) {
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  font-weight: 500;
  padding-bottom: 8px;
  letter-spacing: 0.3px;
}

.add-stock-form :deep(.n-form-item-blank) {
  min-height: auto !important;
}

.add-stock-form :deep(.n-input),
.add-stock-form :deep(.n-input-number) {
  --n-color: var(--n-color-modal) !important;
  --n-color-modal: rgba(20, 19, 60, 0.95) !important;
  --n-color-popup: rgba(20, 19, 60, 0.95) !important;
  background: rgba(20, 19, 60, 0.95) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px;
  color: #fff;
  transition: all 0.2s ease;
}

.add-stock-form :deep(.n-input .n-input-wrapper),
.add-stock-form :deep(.n-input-number .n-input-number__wrapper),
.add-stock-form :deep(.n-input-wrapper),
.add-stock-form :deep(.n-input-number__wrapper),
.add-stock-form :deep(.n-input-number-input),
.add-stock-form :deep(.n-input__input-el),
.add-stock-form :deep(.n-base-input),
.add-stock-form :deep(.n-base-input .n-input-wrapper),
.add-stock-form :deep(.n-input-number__input),
.add-stock-form :deep(.n-input-number__textarea-el),
.add-stock-form :deep(.n-input__input),
.add-stock-form :deep(.n-input__textarea-el),
.add-stock-form :deep(.n-base-input .n-input__input-el),
.add-stock-form :deep(.n-input-number .n-base-input),
.add-stock-form :deep(.n-form-item),
.add-stock-form :deep(.n-form-item-blank),
.add-stock-form :deep(.n-input-number__suffix),
.add-stock-form :deep(.n-input-number__prefix),
.add-stock-form :deep(.n-input-number__add),
.add-stock-form :deep(.n-input-number__sub) {
  --n-color: var(--n-color-modal) !important;
  --n-color-modal: rgba(20, 19, 60, 0.95) !important;
  --n-color-popup: rgba(20, 19, 60, 0.95) !important;
  background: rgba(20, 19, 60, 0.95) !important;
  box-shadow: none !important;
  border-radius: 10px !important;
}

.add-stock-form :deep(.n-input-number-button),
.add-stock-form :deep(.n-input-number__button),
.add-stock-form :deep(.n-input-number--button-type),
.add-stock-form :deep(.n-input-number .n-button),
.add-stock-form :deep(.n-input-number__wrapper .n-button) {
  background: rgba(20, 19, 60, 0.95) !important;
  border-radius: 10px !important;
  border: none !important;
  --n-color: rgba(20, 19, 60, 0.95) !important;
}

.add-stock-form :deep(.n-input-number__actions),
.add-stock-form :deep(.n-input-number__action) {
  background: rgba(20, 19, 60, 0.95) !important;
  --n-color: rgba(20, 19, 60, 0.95) !important;
}

.add-stock-modal :deep(.n-input-number:not(.n-input-number--disabled) .n-input-number__button) {
  background: rgba(20, 19, 60, 0.95) !important;
  --n-color: rgba(20, 19, 60, 0.95) !important;
}

.add-stock-modal .add-stock-form :deep(.n-input-number__button) {
  background: rgba(20, 19, 60, 0.95) !important;
}

.add-stock-form :deep(.n-input:hover),
.add-stock-form :deep(.n-input-number:hover) {
  border-color: rgba(99, 102, 241, 0.5) !important;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.add-stock-form :deep(.n-input:focus),
.add-stock-form :deep(.n-input-number:focus) {
  border-color: #6366f1 !important;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2), 0 4px 16px rgba(99, 102, 241, 0.15) !important;
}

.add-stock-form :deep(.n-input--disabled),
.add-stock-form :deep(.n-input-number--disabled) {
  background: rgba(255, 255, 255, 0.02) !important;
  color: rgba(255, 255, 255, 0.4) !important;
}

.add-stock-form :deep(.n-input__placeholder) {
  color: rgba(255, 255, 255, 0.3) !important;
  font-size: 14px;
}

.add-stock-form :deep(.n-input-number-suffix) {
  color: rgba(255, 255, 255, 0.5) !important;
}

.add-stock-form :deep(.n-button--tertiary-type) {
  background: rgba(99, 102, 241, 0.15) !important;
  border: 1px solid rgba(99, 102, 241, 0.3) !important;
  color: #fff !important;
  border-radius: 12px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.add-stock-form :deep(.n-button--tertiary-type:hover) {
  background: rgba(99, 102, 241, 0.25) !important;
  border-color: rgba(99, 102, 241, 0.5) !important;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

.cancel-btn {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  color: rgba(255, 255, 255, 0.7) !important;
  border-radius: 12px;
  padding: 0 24px !important;
  height: 44px !important;
  font-weight: 500;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
}

.submit-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border: none !important;
  color: #fff !important;
  font-weight: 600;
  border-radius: 12px;
  padding: 0 28px !important;
  height: 44px !important;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3);
  transition: all 0.2s ease;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
}

.table-header :deep(.n-button--text) {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.table-header :deep(.n-button__icon) {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.table-header :deep(.n-button__icon .n-icon) {
  background: transparent !important;
  color: inherit !important;
}

.refresh-icon-btn {
  background: transparent !important;
}

.refresh-icon-btn:hover {
  background: rgba(99, 102, 241, 0.1) !important;
}

.refresh-icon-btn .n-button__icon {
  background: transparent !important;
}

.refresh-icon-btn .n-icon {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.7) !important;
}

.refresh-icon-btn .n-icon svg {
  background: transparent !important;
  fill: currentColor !important;
}
</style>
