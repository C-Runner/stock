<template>
  <div class="home">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>

    <div class="home-header">
      <div class="header-left">
        <div class="header-text">
          <h1>Stock Portfolio</h1>
          <p class="subtitle">Manage your stock holdings</p>
        </div>
      </div>
      <div class="header-right">
        <button class="logout-btn" @click="handleLogout">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9"/>
          </svg>
          <span>Logout</span>
        </button>
      </div>
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

      <div class="table-content" v-if="stocks.length === 0 && !loading">
        <n-empty description="No stocks yet, click 'Add Stock' to add one" />
      </div>
      <div class="table-container" v-else>
        <div class="table-scroll-wrapper">
          <div class="table-header-row">
            <div class="header-stock-cell">Stock</div>
            <div class="header-scroll-part">
              <div class="header-price-cell">Price</div>
              <div class="header-profit-cell">P/L</div>
              <div class="header-cell-item">Qty</div>
              <div class="header-cell-item">Value</div>
              <div class="header-action-cell"></div>
            </div>
          </div>
          <div
            v-for="stock in stocks"
            :key="stock.code"
            class="table-row"
            @click="goToAnalysis(stock.code)"
          >
            <div class="stock-cell">
              <div class="stock-name">{{ stock.name }}</div>
              <div class="stock-code">{{ stock.code }}</div>
            </div>
            <div class="table-scroll-part">
              <div class="price-cell">
                <div class="price-current">¥{{ stock.currentPrice.toFixed(2) }}</div>
                <div class="price-cost">¥{{ stock.buyPrice.toFixed(2) }}</div>
              </div>
              <div class="profit-cell">
                <div class="profit-amount" :class="(stock.currentPrice - stock.buyPrice) * stock.quantity >= 0 ? 'profit-up' : 'profit-down'">
                  {{ (stock.currentPrice - stock.buyPrice) * stock.quantity >= 0 ? '+' : '' }}¥{{ ((stock.currentPrice - stock.buyPrice) * stock.quantity).toFixed(2) }}
                </div>
                <div class="profit-percent" :class="stock.currentPrice >= stock.buyPrice ? 'profit-up' : 'profit-down'">
                  {{ stock.buyPrice > 0 ? ((stock.currentPrice / stock.buyPrice - 1) * 100).toFixed(2) : '0.00' }}%
                </div>
              </div>
              <div class="cell-item">{{ stock.quantity.toLocaleString() }}</div>
              <div class="cell-item">¥{{ (stock.currentPrice * stock.quantity).toFixed(2) }}</div>
              <div class="action-cell">
                <n-button text @click.stop="showDeletePopup(stock.code, stock.name)" class="delete-btn">
                  <template #icon>
                    <n-icon><IconDelete /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bottom-tabs">
      <div class="tab-item" @click="router.push('/watchlist')">
        <div class="tab-icon">
          <n-icon size="22"><IconSearch /></n-icon>
        </div>
        <span class="tab-label">Watchlist</span>
      </div>
      <div class="tab-item primary" @click="showAddModal = true">
        <div class="tab-icon">
          <n-icon size="24"><IconPlus /></n-icon>
        </div>
        <span class="tab-label">Add Stock</span>
      </div>
    </div>

    <n-modal
      v-model:show="showAddModal"
      preset="card"
      class="add-stock-modal"
      :style="{
        '--n-color': 'rgba(20, 19, 60, 0.7)',
        '--n-color-modal': 'rgba(20, 19, 60, 0.7)',
        background: 'rgba(20, 19, 60, 0.7)',
        backdropFilter: 'blur(24px)'
      }"
      :mask="{ style: 'background: rgba(0, 0, 0, 0.5); backdrop-filter: blur(4px);' }"
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

    <n-modal
      v-model:show="showPopup"
      preset="card"
      class="delete-modal"
      :style="{
        '--n-color': 'rgba(20, 19, 60, 0.5)',
        '--n-color-modal': 'rgba(20, 19, 60, 0.5)',
        background: 'rgba(20, 19, 60, 0.5)',
        backdropFilter: 'blur(32px) saturate(180%)'
      }"
      :mask="{ style: 'background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(8px) saturate(150%);' }"
      :icon="() => null"
      :mask-closable="false"
    >
      <template #header>
        <div class="modal-header">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="modal-icon delete-icon">
            <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
          </svg>
          <span>Delete Stock</span>
        </div>
      </template>
      <div class="modal-body delete-modal-body">
        <p>Are you sure you want to delete <strong>{{ selectedStockName }}</strong>?</p>
      </div>
      <template #action>
        <div class="modal-footer">
          <n-button @click="showPopup = false" class="cancel-btn">Cancel</n-button>
          <n-button type="error" @click="confirmDelete" class="delete-confirm-btn">Delete</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NModal, NForm, NFormItem,
  NInput, NInputNumber, NEmpty, NIcon
} from 'naive-ui'
import { stockApi, authApi, type Stock, type StockRequest } from '../api'
import { IconPlus, IconWallet, IconCoin, IconTrend, IconDataLine, IconRefresh, IconSearch, IconDelete } from '../components/icons'

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

const showPopup = ref(false)
const selectedStockCode = ref('')
const selectedStockName = ref('')

const showDeletePopup = (code: string, name: string) => {
  selectedStockCode.value = code
  selectedStockName.value = name
  showPopup.value = true
}

const confirmDelete = async () => {
  await stockApi.deleteStock(selectedStockCode.value)
  showPopup.value = false
  await fetchStocks()
}

const goToAnalysis = (code: string) => {
  router.push({ path: `/analysis/${code}`, query: { from: '/home' } })
}


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

onMounted(() => fetchStocks())

const handleLogout = async () => {
  try {
    await authApi.logout()
  } catch {
    // ignore error
  }
  localStorage.removeItem('token')
  localStorage.removeItem('tokenExpiry')
  router.push('/login')
}
</script>

<style scoped>
.home {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  height: 100dvh;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
  padding: 16px;
  padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  display: flex;
  flex-direction: column;
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
  margin-bottom: 12px;
  position: relative;
  flex: 0 0 auto;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  background: rgba(255, 107, 107, 0.15);
  border-color: rgba(255, 107, 107, 0.3);
  color: #ff6b6b;
}

.logout-btn svg {
  width: 16px;
  height: 16px;
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

.stats-row {
  margin-bottom: 16px;
  position: relative;
  flex: 0 0 auto;
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
  gap: 8px;
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
  padding: 0;
  box-sizing: border-box;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
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

.table-card :deep(.n-data-table tr:hover .n-data-table-td) {
  background: rgba(99, 102, 241, 0.08) !important;
}

.table-card :deep(.profit-up) {
  color: #ff6b6b;
}

.table-card :deep(.profit-down) {
  color: #38ef7d;
}

.price-current {
  font-size: 14px;
  color: #fff;
}

.price-cost {
  font-size: 14px;
  color: #fff;
}

.profit-amount {
  font-size: 14px;
  color: #fff;
}

.profit-percent {
  font-size: 14px;
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
    padding: 16px 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }

  .home-header {
    flex-direction: column;
    gap: 12px;
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
    padding: 12px 8px;
    padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  }

  .home-header {
    margin-bottom: 12px;
  }

  .stats-row {
    margin-bottom: 12px;
  }

  .header-text h1 {
    font-size: 20px;
  }

  .subtitle {
    font-size: 12px;
  }

  .stat-card {
    padding: 12px;
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
  --n-color: rgba(20, 19, 60, 0.7) !important;
  --n-border-radius: 16px !important;
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
  gap: 10px;
  font-weight: 600;
  font-size: 15px;
  color: #fff;
  margin-bottom: 0;
}

.modal-icon {
  width: 22px;
  height: 22px;
  padding: 5px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 6px;
  color: #fff;
  flex-shrink: 0;
}

.modal-body,
.add-stock-modal :deep(.n-card),
.add-stock-modal :deep(.n-card__content),
.add-stock-modal :deep(.n-card .n-card__content),
.add-stock-modal :deep(.n-base-card),
.add-stock-modal :deep(.n-base-card .n-card__content) {
  padding-top: 16px;
  background: rgba(20, 19, 60, 0.95) !important;
  background-color: rgba(20, 19, 60, 0.95) !important;
  backdrop-filter: blur(24px) !important;
}

.add-stock-form :deep(.n-form-item) {
  margin-bottom: 14px;
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
  gap: 10px;
}

.cancel-btn {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  color: rgba(255, 255, 255, 0.7) !important;
  border-radius: 10px;
  padding: 0 20px !important;
  height: 36px !important;
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
  border-radius: 10px;
  padding: 0 20px !important;
  height: 36px !important;
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

.table-container {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  width: 100%;
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
.header-profit-cell,
.header-cell-item {
  width: 90px;
  min-width: 90px;
  display: flex;
  align-items: center;
  font-weight: 600;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.header-action-cell {
  width: 50px;
  min-width: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
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

.price-cell {
  width: 90px;
  min-width: 90px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  line-height: 1.3;
}

.profit-cell {
  width: 90px;
  min-width: 90px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  line-height: 1.3;
}

.cell-item {
  width: 90px;
  min-width: 90px;
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #fff;
}

.action-cell {
  width: 50px;
  min-width: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.table-row:hover {
  background: rgba(99, 102, 241, 0.08);
}

.delete-btn {
  color: rgba(255, 107, 107, 0.6) !important;
}

.delete-btn:hover {
  color: #ff6b6b !important;
}

.popup-content {
  padding: 8px;
}

.popup-title {
  font-size: 14px;
  color: #fff;
  margin-bottom: 12px;
}

.popup-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.delete-modal {
  --n-color: rgba(20, 19, 60, 0.5) !important;
  --n-border-radius: 20px !important;
  --n-border: 1px solid rgba(255, 107, 107, 0.15) !important;
  background-color: var(--n-color) !important;
  backdrop-filter: blur(32px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(32px) saturate(180%) !important;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.1) !important;
}

.delete-modal :deep(.n-card) {
  background-color: rgba(20, 19, 60, 0.4) !important;
  backdrop-filter: blur(32px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(32px) saturate(180%) !important;
  border: 1px solid rgba(255, 107, 107, 0.1) !important;
}

.delete-modal :deep(.n-dialog__close) {
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.6) !important;
  transition: all 0.25s ease;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.delete-modal :deep(.n-dialog__close:hover) {
  background: rgba(255, 107, 107, 0.2);
  border-color: rgba(255, 107, 107, 0.3);
  color: #ff6b6b !important;
  transform: rotate(90deg);
}

.delete-icon {
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.3), rgba(238, 90, 36, 0.3));
  padding: 5px;
  border-radius: 6px;
  width: 22px;
  height: 22px;
  opacity: 0.7;
}

.delete-modal-body {
  text-align: center;
  padding: 16px 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 15px;
}

.delete-modal-body strong {
  color: #fff;
}

.delete-confirm-btn {
  background: linear-gradient(135deg, #ff6b6b, #ee5a24) !important;
  border: none !important;
  color: #fff !important;
  font-weight: 600;
  border-radius: 10px;
  padding: 0 20px !important;
  height: 36px !important;
  box-shadow: 0 4px 16px rgba(255, 107, 107, 0.3);
  transition: all 0.25s ease;
}

.delete-confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.4);
}
</style>
