<script setup lang="ts">
import { NModal, NButton } from 'naive-ui'

defineProps<{
  show: boolean
  stockName: string
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  'confirm': []
}>()
</script>

<template>
  <n-modal
    :show="show"
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
    @update:show="emit('update:show', $event)"
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
      <p>Are you sure you want to delete <strong>{{ stockName }}</strong>?</p>
    </div>
    <template #action>
      <div class="modal-footer">
        <n-button @click="emit('update:show', false)" class="cancel-btn">Cancel</n-button>
        <n-button type="error" @click="emit('confirm')" class="delete-confirm-btn">Delete</n-button>
      </div>
    </template>
  </n-modal>
</template>

<style scoped>
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
