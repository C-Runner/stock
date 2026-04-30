<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NInputNumber, NButton, NSpace } from 'naive-ui'
import { stockApi, type StockRequest } from '../api'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  'added': []
}>()

const submitting = ref(false)
const lookingUp = ref(false)

const stockForm = ref<StockRequest>({
  code: '',
  name: '',
  currentPrice: 0,
  quantity: 0,
  buyPrice: 0
})

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
    emit('update:show', false)
    stockForm.value = { code: '', name: '', currentPrice: 0, quantity: 0, buyPrice: 0 }
    emit('added')
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

import { ref } from 'vue'
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
        <n-button @click="emit('update:show', false)" class="cancel-btn">Cancel</n-button>
        <n-button type="primary" @click="handleAddStock" :loading="submitting" class="submit-btn">Add Stock</n-button>
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
</style>
