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
    preset="card"
    title="Search Stocks"
    style="width: 500px; max-width: 90vw"
    @update:show="(v) => emit('update:show', v)"
  >
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
  </n-modal>
</template>

<style scoped>
.selected-stock {
  padding: 12px;
  background: #1a1a1a;
  border-radius: 8px;
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
  color: #999;
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
  padding: 2px 8px;
  border-radius: 4px;
}

.up {
  color: #ef5350;
  background: rgba(239, 83, 80, 0.1);
}

.down {
  color: #26a69a;
  background: rgba(38, 166, 154, 0.1);
}
</style>
