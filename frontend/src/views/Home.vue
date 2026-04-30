<template>
  <div class="home">
    <BackgroundOrbs />

    <div class="home-header">
      <div class="header-left">
        <div class="header-text">
          <h1>Stock Portfolio</h1>
          <p class="subtitle">Manage your stock holdings</p>
        </div>
      </div>
      <div class="header-right">
        <button class="settings-btn" @click="goToSettings">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
          </svg>
          <span>Settings</span>
        </button>
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
                <n-button text @click.stop="showEditModal(stock)" class="edit-btn">
                  <template #icon>
                    <n-icon><IconEdit /></n-icon>
                  </template>
                </n-button>
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

    <BottomTabs :tabs="bottomTabs" />

    <AddStockModal v-model:show="showAddModal" @added="fetchStocks" />
    <EditStockModal v-model:show="showEditStockModal" :stock="editStockForm" @saved="fetchStocks" />
    <DeleteModal v-model:show="showPopup" :stock-name="selectedStockName" @confirm="confirmDelete" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NEmpty, NIcon } from 'naive-ui'
import { stockApi, authApi, type Stock, type StockRequest } from '../api'
import { IconPlus, IconWallet, IconCoin, IconTrend, IconDataLine, IconRefresh, IconSearch, IconDelete, IconEdit } from '../components/icons'
import BackgroundOrbs from '../components/BackgroundOrbs.vue'
import BottomTabs from '../components/BottomTabs.vue'
import AddStockModal from '../components/AddStockModal.vue'
import EditStockModal from '../components/EditStockModal.vue'
import DeleteModal from '../components/DeleteModal.vue'

const router = useRouter()

const stocks = ref<Stock[]>([])
const loading = ref(false)
const showAddModal = ref(false)

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

const showEditStockModal = ref(false)
const editStockForm = ref<StockRequest | null>(null)

const bottomTabs = computed(() => [
  { icon: IconSearch, label: 'Watchlist', action: () => router.push('/watchlist') },
  { icon: IconPlus, label: 'Add Stock', action: () => { showAddModal.value = true }, primary: true }
])

const showDeletePopup = (code: string, name: string) => {
  selectedStockCode.value = code
  selectedStockName.value = name
  showPopup.value = true
}

const showEditModal = (stock: Stock) => {
  editStockForm.value = {
    code: stock.code,
    name: stock.name,
    currentPrice: stock.currentPrice,
    quantity: stock.quantity,
    buyPrice: stock.buyPrice
  }
  showEditStockModal.value = true
}

const confirmDelete = async () => {
  await stockApi.deleteStock(selectedStockCode.value)
  showPopup.value = false
  await fetchStocks()
}

const goToAnalysis = (code: string) => {
  router.push({ path: `/analysis/${code}`, query: { from: '/home' } })
}

const goToSettings = () => {
  router.push('/settings')
}

const fetchStocks = async () => {
  loading.value = true
  try {
    stocks.value = await stockApi.getStocks()
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

onMounted(() => fetchStocks())

const handleLogout = async () => {
  try {
    await authApi.logout()
  } catch {
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

.settings-btn {
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

.settings-btn:hover {
  background: rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.3);
  color: #667eea;
}

.settings-btn svg {
  width: 16px;
  height: 16px;
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

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
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

.table-card :deep(.n-data-table-tr),
.table-card :deep(.n-base-table),
.table-card :deep(.n-base-table-tbody),
.table-card :deep(.n-base-table-tr),
.table-card :deep(.n-base-td) {
  background: transparent !important;
}

.table-card :deep(.n-data-table tr:hover .n-data-table-td) {
  background: rgba(99, 102, 241, 0.08) !important;
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

.refresh-icon-btn .n-button__icon,
.refresh-icon-btn .n-icon,
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
  background: rgba(10, 10, 15, 1);
}

.header-action-cell {
  width: 80px;
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
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
  width: 80px;
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.delete-btn {
  color: rgba(255, 107, 107, 0.6) !important;
}

.delete-btn:hover {
  color: #ff6b6b !important;
}

.edit-btn {
  color: rgba(99, 102, 241, 0.6) !important;
}

.edit-btn:hover {
  color: #6366f1 !important;
}

.profit-up { color: #ff6b6b; }
.profit-down { color: #38ef7d; }
</style>
