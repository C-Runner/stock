<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NModal, NInput, NButton, NSpace, NSpin,
  NEmpty, NList, NListItem, NThing, NDivider
} from 'naive-ui'
import { stockApi, type Stock, type StockQuote } from '../api'

interface Props {
  show: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'select', code: string, name: string): void
}>()

const searchQuery = ref('')
const searchResults = ref<Stock[]>([])
const searching = ref(false)
const selectedQuote = ref<StockQuote | null>(null)
const lookingUp = ref(false)

const close = () => {
  emit('update:show', false)
  searchQuery.value = ''
  searchResults.value = []
  selectedQuote.value = null
}

const handleSearch = async () => {
  if (!searchQuery.value.trim()) return

  searching.value = true
  try {
    searchResults.value = await stockApi.searchStocks(searchQuery.value)
  } catch (error) {
    console.error(error)
  } finally {
    searching.value = false
  }
}

const lookupStock = async (code: string) => {
  lookingUp.value = true
  try {
    selectedQuote.value = await stockApi.getQuote(code)
  } catch (error) {
    console.error(error)
  } finally {
    lookingUp.value = false
  }
}

const selectStock = (stock: Stock) => {
  lookupStock(stock.code)
}

const handleAddToWatchlist = () => {
  if (selectedQuote.value) {
    emit('select', selectedQuote.value.code, selectedQuote.value.name)
    close()
  }
}

watch(searchQuery, (newVal) => {
  if (!newVal) {
    searchResults.value = []
    selectedQuote.value = null
  }
})
</script>

<template>
  <n-modal
    :show="props.show"
    preset="dialog"
    class="stock-search-modal"
    :icon="() => null"
    @update:show="(v) => emit('update:show', v)"
  >
    <template #header>
      <div class="modal-header">Search Stocks</div>
    </template>
    <div class="modal-body">
      <n-space vertical :size="16">
        <n-input
          v-model:value="searchQuery"
          placeholder="Enter stock code or name..."
          clearable
          @keyup.enter="handleSearch"
        >
          <template #suffix>
            <n-button
              :loading="searching"
              :disabled="!searchQuery"
              size="small"
              type="primary"
              @click="handleSearch"
            >
              Search
            </n-button>
          </template>
        </n-input>

        <n-spin :show="lookingUp">
          <n-list v-if="searchResults.length > 0" hoverable clickable>
            <n-list-item
              v-for="stock in searchResults"
              :key="stock.code"
              @click="selectStock(stock)"
            >
              <n-thing :title="stock.name" :description="stock.code" />
            </n-list-item>
          </n-list>
          <n-empty v-else-if="!searching && searchQuery" description="No stocks found" />
        </n-spin>

        <div v-if="selectedQuote" class="selected-stock">
          <n-divider />
          <div class="stock-preview">
            <div class="stock-header">
              <span class="stock-name">{{ selectedQuote.name }}</span>
              <span class="stock-code">{{ selectedQuote.code }}</span>
            </div>
            <div class="stock-price">
              <span class="price">¥{{ selectedQuote.current.toFixed(2) }}</span>
              <span :class="['change', selectedQuote.current >= selectedQuote.prevClose ? 'up' : 'down']">
                {{ selectedQuote.current >= selectedQuote.prevClose ? '+' : '' }}
                {{ (((selectedQuote.current - selectedQuote.prevClose) / selectedQuote.prevClose) * 100).toFixed(2) }}%
              </span>
            </div>
            <n-button type="primary" block @click="handleAddToWatchlist">
              Add to Watchlist
            </n-button>
          </div>
        </div>
      </n-space>
    </div>
  </n-modal>
</template>

<style scoped>
.modal-header {
  font-weight: 600;
  font-size: 20px;
  color: #fff;
  letter-spacing: -0.3px;
  padding: 4px 0;
}

.modal-body {
  padding-top: 24px;
}

.selected-stock {
  padding: 20px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.08), rgba(139, 92, 246, 0.05));
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 18px;
  transition: all 0.25s ease;
}

.stock-preview {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.stock-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stock-name {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.3px;
}

.stock-code {
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 10px;
  border-radius: 6px;
}

.stock-price {
  display: flex;
  align-items: baseline;
  gap: 14px;
}

.price {
  font-size: 28px;
  font-weight: bold;
  color: #fff;
  letter-spacing: -0.5px;
}

.change {
  font-size: 14px;
  padding: 6px 12px;
  border-radius: 10px;
  font-weight: 600;
}

.up {
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.15);
}

.down {
  color: #38ef7d;
  background: rgba(56, 239, 125, 0.15);
}

.stock-search-modal :deep(.n-dialog) {
  background: transparent !important;
  box-shadow: none !important;
  padding: 0 !important;
  border-radius: 24px;
  overflow: hidden;
}

.stock-search-modal :deep(.n-base-overlay) {
  background: linear-gradient(
    135deg,
    rgba(99, 102, 241, 0.12) 0%,
    rgba(139, 92, 246, 0.08) 40%,
    rgba(6, 182, 212, 0.05) 70%,
    rgba(0, 0, 0, 0.6) 100%
  ) !important;
  backdrop-filter: blur(12px) saturate(180%);
}

.stock-search-modal :deep(.n-dialog__content) {
  background: linear-gradient(
    145deg,
    rgba(30, 27, 75, 0.95) 0%,
    rgba(20, 19, 60, 0.98) 100%
  ) !important;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-bottom: none;
  backdrop-filter: blur(24px);
  padding: 24px !important;
  border-radius: 24px 24px 0 0 !important;
}

.stock-search-modal :deep(.n-dialog__content--last) {
  background: linear-gradient(
    145deg,
    rgba(25, 23, 65, 0.95) 0%,
    rgba(15, 14, 50, 0.98) 100%
  ) !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-top: 1px solid rgba(255, 255, 255, 0.06) !important;
  backdrop-filter: blur(24px);
  padding: 20px 24px 24px !important;
  border-radius: 0 0 24px 24px !important;
}

.stock-search-modal :deep(.n-dialog__close) {
  top: 20px;
  right: 20px;
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.5) !important;
  transition: all 0.25s ease;
}

.stock-search-modal :deep(.n-dialog__close:hover) {
  background: rgba(255, 107, 107, 0.15) !important;
  border-color: rgba(255, 107, 107, 0.3) !important;
  color: #ff6b6b !important;
  transform: rotate(90deg);
}

.stock-search-modal :deep(.n-input) {
  background: rgba(255, 255, 255, 0.06) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 14px !important;
  transition: all 0.25s ease !important;
}

.stock-search-modal :deep(.n-input:hover) {
  border-color: rgba(99, 102, 241, 0.4) !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

.stock-search-modal :deep(.n-input:focus) {
  border-color: #6366f1 !important;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2), 0 8px 24px rgba(99, 102, 241, 0.15) !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

.stock-search-modal :deep(.n-input__input-el) {
  color: #fff !important;
  font-size: 14px !important;
}

.stock-search-modal :deep(.n-input__placeholder) {
  color: rgba(255, 255, 255, 0.3) !important;
}

.stock-search-modal :deep(.n-input__suffix) {
  margin-right: 4px;
}

.stock-search-modal :deep(.n-list) {
  background: transparent !important;
}

.stock-search-modal :deep(.n-list-item) {
  background: rgba(255, 255, 255, 0.03) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 14px !important;
  margin-bottom: 10px !important;
  transition: all 0.25s ease !important;
  padding: 14px 16px !important;
}

.stock-search-modal :deep(.n-list-item:hover) {
  background: rgba(99, 102, 241, 0.12) !important;
  border-color: rgba(99, 102, 241, 0.3) !important;
  transform: translateX(4px);
}

.stock-search-modal :deep(.n-thing) {
  padding: 0 !important;
}

.stock-search-modal :deep(.n-thing-main) {
  gap: 4px;
}

.stock-search-modal :deep(.n-thing-title) {
  color: #fff !important;
  font-weight: 600 !important;
  font-size: 15px !important;
}

.stock-search-modal :deep(.n-thing-description) {
  color: rgba(255, 255, 255, 0.4) !important;
  font-size: 13px !important;
  font-family: 'JetBrains Mono', monospace;
}

.stock-search-modal :deep(.n-divider) {
  border-color: rgba(255, 255, 255, 0.08) !important;
  margin: 16px 0 !important;
}

.stock-search-modal :deep(.n-button--primary-type) {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border: none !important;
  border-radius: 14px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
  height: 46px !important;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.35) !important;
  transition: all 0.25s ease !important;
}

.stock-search-modal :deep(.n-button--primary-type:hover) {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.45) !important;
}

.stock-search-modal :deep(.n-spin) {
  min-height: 100px;
}

.stock-search-modal :deep(.n-empty) {
  padding: 32px 0;
}
</style>
