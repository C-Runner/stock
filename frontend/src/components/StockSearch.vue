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
              <span :class="['change', selectedQuote.current >= selectedQuote.open ? 'up' : 'down']">
                {{ selectedQuote.current >= selectedQuote.open ? '+' : '' }}
                {{ (((selectedQuote.current - selectedQuote.open) / selectedQuote.open) * 100).toFixed(2) }}%
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
  font-size: 18px;
  color: #fff;
}

.modal-body {
  padding-top: 20px;
}

.selected-stock {
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
}

.stock-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stock-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stock-name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.stock-code {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
}

.stock-price {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.price {
  font-size: 24px;
  font-weight: bold;
  color: #fff;
}

.change {
  font-size: 14px;
  padding: 4px 10px;
  border-radius: 8px;
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
    rgba(99, 102, 241, 0.15) 0%,
    rgba(139, 92, 246, 0.1) 50%,
    rgba(0, 0, 0, 0.5) 100%
  ) !important;
  backdrop-filter: blur(8px) saturate(150%);
}

.stock-search-modal :deep(.n-dialog__content) {
  background: rgba(20, 19, 60, 0.95) !important;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-bottom: none;
  backdrop-filter: blur(24px);
  padding: 20px 24px !important;
}

.stock-search-modal :deep(.n-dialog__content--last) {
  background: rgba(20, 19, 60, 0.95) !important;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-top: none;
  padding: 0 24px 20px !important;
}

.stock-search-modal :deep(.n-dialog__close) {
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.6) !important;
  transition: all 0.2s ease;
}

.stock-search-modal :deep(.n-dialog__close:hover) {
  background: rgba(255, 107, 107, 0.15);
  border-color: rgba(255, 107, 107, 0.3);
  color: #ff6b6b !important;
}

.stock-search-modal :deep(.n-input) {
  background: rgba(255, 255, 255, 0.06) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 12px !important;
  transition: all 0.2s ease !important;
}

.stock-search-modal :deep(.n-input:hover) {
  border-color: rgba(99, 102, 241, 0.4) !important;
}

.stock-search-modal :deep(.n-input:focus) {
  border-color: #6366f1 !important;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2) !important;
}

.stock-search-modal :deep(.n-input__input-el) {
  color: #fff !important;
  font-size: 14px !important;
}

.stock-search-modal :deep(.n-input__placeholder) {
  color: rgba(255, 255, 255, 0.3) !important;
}

.stock-search-modal :deep(.n-list-item) {
  background: rgba(255, 255, 255, 0.03) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 12px !important;
  margin-bottom: 8px !important;
  transition: all 0.2s ease !important;
}

.stock-search-modal :deep(.n-list-item:hover) {
  background: rgba(99, 102, 241, 0.12) !important;
  border-color: rgba(99, 102, 241, 0.3) !important;
}

.stock-search-modal :deep(.n-thing-title) {
  color: #fff !important;
  font-weight: 500 !important;
}

.stock-search-modal :deep(.n-thing-description) {
  color: rgba(255, 255, 255, 0.4) !important;
  font-size: 12px !important;
}

.stock-search-modal :deep(.n-divider) {
  border-color: rgba(255, 255, 255, 0.08) !important;
  margin: 12px 0 !important;
}

.stock-search-modal :deep(.n-button--primary-type) {
  background: linear-gradient(135deg, #6366f1, #8b5cf6) !important;
  border: none !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3) !important;
  transition: all 0.2s ease !important;
}

.stock-search-modal :deep(.n-button--primary-type:hover) {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4) !important;
}
</style>
