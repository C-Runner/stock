<script setup lang="ts">
import { NIcon } from 'naive-ui'
import { IconHome, IconPlus, IconStar } from './icons'

defineProps<{
  activeTab: 'home' | 'add' | 'watchlist'
  onAdd: () => void
}>()

const emit = defineEmits<{
  navigate: [tab: 'home' | 'watchlist']
}>()
</script>

<template>
  <div class="bottom-tabs">
    <div
      class="tab-item"
      :class="{ active: activeTab === 'home' }"
      @click="emit('navigate', 'home')"
    >
      <div class="tab-icon">
        <n-icon :size="22"><IconHome /></n-icon>
      </div>
      <span class="tab-label">Home</span>
    </div>

    <div
      class="tab-item primary"
      :class="{ active: activeTab === 'add' }"
      @click="onAdd"
    >
      <div class="tab-icon">
        <n-icon :size="22"><IconPlus /></n-icon>
      </div>
      <span class="tab-label">Add</span>
    </div>

    <div
      class="tab-item"
      :class="{ active: activeTab === 'watchlist' }"
      @click="emit('navigate', 'watchlist')"
    >
      <div class="tab-icon">
        <n-icon :size="22"><IconStar /></n-icon>
      </div>
      <span class="tab-label">Watchlist</span>
    </div>
  </div>
</template>

<style scoped>
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
  gap: 12px;
  padding: 0 16px;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
  touch-action: none;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  padding: 8px 24px;
  border-radius: 16px;
  transition: all 0.2s ease;
  color: rgba(255, 255, 255, 0.6);
  min-width: 80px;
  flex: 1;
  max-width: 120px;
}

.tab-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
}

.tab-item.active {
  color: rgba(255, 255, 255, 0.9);
}

.tab-item.active .tab-icon {
  transform: scale(1.1);
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

.tab-item.primary.active {
  background: linear-gradient(135deg, #818cf8, #a78bfa);
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.6);
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  transition: transform 0.2s ease;
}

.tab-label {
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}
</style>
