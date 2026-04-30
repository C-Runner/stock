<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NInputNumber, NButton } from 'naive-ui'
import { stockApi, type StockRequest } from '../api'
import { ref } from 'vue'

const props = defineProps<{
  show: boolean
  stock: StockRequest | null
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  'saved': []
}>()

const submitting = ref(false)
const editStockForm = ref<StockRequest>({
  code: '',
  name: '',
  currentPrice: 0,
  quantity: 0,
  buyPrice: 0
})

import { watch } from 'vue'
watch(() => props.show, (val) => {
  if (val && props.stock) {
    editStockForm.value = { ...props.stock }
  }
})

const handleEditStock = async () => {
  if (!editStockForm.value.code) return
  if (!editStockForm.value.currentPrice || editStockForm.value.currentPrice <= 0) return
  if (!editStockForm.value.quantity || editStockForm.value.quantity <= 0) return
  if (!editStockForm.value.buyPrice || editStockForm.value.buyPrice <= 0) return

  submitting.value = true
  try {
    await stockApi.updateStock(editStockForm.value.code, editStockForm.value)
    emit('update:show', false)
    emit('saved')
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <n-modal
    :show="show"
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
    @update:show="emit('update:show', $event)"
  >
    <template #header>
      <div class="modal-header">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="modal-icon">
          <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
          <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
        <span>Edit Stock</span>
      </div>
    </template>
    <div class="modal-body">
      <n-form :model="editStockForm" class="add-stock-form">
        <n-form-item label="Stock Code" path="code">
          <n-input v-model:value="editStockForm.code" disabled />
        </n-form-item>
        <n-form-item label="Stock Name" path="name">
          <n-input v-model:value="editStockForm.name" disabled />
        </n-form-item>
        <n-form-item label="Current Price" path="currentPrice">
          <n-input-number
            v-model:value="editStockForm.currentPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="Quantity" path="quantity">
          <n-input-number
            v-model:value="editStockForm.quantity"
            :min="0"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="Buy Price" path="buyPrice">
          <n-input-number
            v-model:value="editStockForm.buyPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </n-form-item>
      </n-form>
    </div>
    <template #action>
      <div class="modal-footer">
        <n-button @click="emit('update:show', false)" class="cancel-btn">Cancel</n-button>
        <n-button type="primary" @click="handleEditStock" :loading="submitting" class="submit-btn">Save Changes</n-button>
      </div>
    </template>
  </n-modal>
</template>

<style scoped>
.add-stock-modal {
  --n-color: rgba(20, 19, 60, 0.7) !important;
  --n-border-radius: 16px !important;
  --n-border: 1px solid rgba(99, 102, 241, 0.2) !important;
  background-color: var(--n-color) !important;
}

.add-stock-modal :deep(.n-card) {
  background-color: rgba(20, 19, 60, 0.95) !important;
}

.modal-body,
.add-stock-modal :deep(.n-card__content) {
  padding-top: 16px;
  background: rgba(20, 19, 60, 0.95) !important;
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
}

.add-stock-form :deep(.n-form-item-blank) {
  min-height: auto !important;
}

.add-stock-form :deep(.n-input),
.add-stock-form :deep(.n-input-number) {
  background: rgba(20, 19, 60, 0.95) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px;
  color: #fff;
}

.add-stock-form :deep(.n-input .n-input-wrapper),
.add-stock-form :deep(.n-input-number .n-input-number__wrapper),
.add-stock-form :deep(.n-input-wrapper),
.add-stock-form :deep(.n-input-number__wrapper) {
  background: rgba(20, 19, 60, 0.95) !important;
  box-shadow: none !important;
  border-radius: 10px !important;
}

.add-stock-form :deep(.n-input:hover),
.add-stock-form :deep(.n-input-number:hover) {
  border-color: rgba(99, 102, 241, 0.5) !important;
}

.add-stock-form :deep(.n-input--disabled),
.add-stock-form :deep(.n-input-number--disabled) {
  background: rgba(255, 255, 255, 0.02) !important;
  color: rgba(255, 255, 255, 0.4) !important;
}
</style>
