<script setup lang="ts">
import { NIcon } from 'naive-ui'

interface TabItem {
  icon: any
  label: string
  action: () => void
  primary?: boolean
}

defineProps<{
  tabs: TabItem[]
}>()
</script>

<template>
  <div class="bottom-tabs">
    <div
      v-for="(tab, index) in tabs"
      :key="index"
      class="tab-item"
      :class="{ primary: tab.primary }"
      @click="tab.action"
    >
      <div class="tab-icon">
        <n-icon :size="tab.primary ? 24 : 22"><component :is="tab.icon" /></n-icon>
      </div>
      <span class="tab-label">{{ tab.label }}</span>
    </div>
  </div>
</template>
