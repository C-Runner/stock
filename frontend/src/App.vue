<script setup lang="ts">
import { NMessageProvider, NConfigProvider, darkTheme } from 'naive-ui'

const themeOverrides = {
  common: {
    primaryColor: '#6366f1',
    primaryColorHover: '#8b5cf6',
    primaryColorPressed: '#4f46e5',
  }
}
</script>

<template>
  <n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <div class="app-container">
        <router-view />
      </div>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
.app-container {
  width: 100%;
  height: 100dvh;
  background: #0a0a0f;
  overflow: hidden;
  position: relative;
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
