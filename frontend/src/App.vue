<script setup lang="ts">
import { NLayout, NLayoutContent, NMessageProvider, NConfigProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#6366f1',
    primaryColorHover: '#8b5cf6',
    primaryColorPressed: '#4f46e5',
  },
  Card: {
    color: 'rgba(20, 19, 60, 0.7)',
    borderColor: 'rgba(99, 102, 241, 0.2)',
    borderRadius: '24px',
    contentColor: 'rgba(20, 19, 60, 0.7)',
  },
  Dialog: {
    color: 'rgba(20, 19, 60, 0.7)',
    borderRadius: '24px',
    border: '1px solid rgba(99, 102, 241, 0.2)',
  },
  DataTable: {
    color: 'transparent',
    thColor: 'transparent',
    tdColor: 'transparent',
    thColorHover: 'transparent',
    tdColorHover: 'transparent',
    borderColor: 'transparent',
  }
}
</script>

<template>
  <n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-layout style="min-height: 100vh" class="dark-layout">
      <n-layout-content content-style="padding: 0">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
.dark-layout {
  background: #0a0a0f !important;
}
</style>

<style>
.n-button__icon {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
  --n-color: rgba(0,0,0,0) !important;
  --n-icon-color: rgba(0,0,0,0) !important;
}

.n-button__icon::before,
.n-button__icon::after {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
}

.n-button--text .n-button__icon {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
  --n-color: rgba(0,0,0,0) !important;
}

.n-icon,
.n-icon-slot,
.n-base-icon,
.n-button .n-icon-slot,
.n-button--text .n-icon-slot {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
  --n-color: rgba(0,0,0,0) !important;
}

.n-icon::before,
.n-icon::after {
  background-color: transparent !important;
}

.n-button__icon .n-icon,
.n-button__icon .n-icon-slot,
.n-button__icon .n-base-icon {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
  --n-color: rgba(0,0,0,0) !important;
}

.n-button__icon .n-icon svg,
.n-button__icon svg {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
}

.n-button__icon > * {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  background-image: none !important;
}

[class*="icon-slot"] {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
  --n-color: rgba(0,0,0,0) !important;
}

.n-base-icon {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
}

.n-icon[class] {
  background-color: transparent !important;
  background: rgba(0,0,0,0) !important;
}
</style>