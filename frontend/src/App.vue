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
.n-button__icon,
.n-icon,
.n-icon-slot,
.n-base-icon,
.n-button .n-icon-slot,
.n-button--text .n-icon-slot,
[class*="icon-slot"] {
  background-color: transparent !important;
  background-image: none !important;
}

.n-button__icon::before,
.n-button__icon::after,
.n-icon::before,
.n-icon::after {
  background-color: transparent !important;
}
</style>
