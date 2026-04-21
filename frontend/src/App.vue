<script setup lang="ts">
import { NLayout, NLayoutHeader, NLayoutContent, NMenu, NAvatar, NDropdown, NMessageProvider, NConfigProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import { useRouter, useRoute } from 'vue-router'

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

const router = useRouter()
const route = useRoute()

const menuOptions = [
  {
    label: 'Home',
    key: '/'
  },
  {
    label: 'Watchlist',
    key: '/watchlist'
  }
]

const handleMenuUpdate = (key: string) => {
  router.push(key)
}

const handleUserAction = (key: string) => {
  if (key === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('tokenExpiry')
    router.push('/login')
  }
}

const userOptions = [
  {
    label: 'Logout',
    key: 'logout'
  }
]
</script>

<template>
  <n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-layout style="min-height: 100vh" class="dark-layout">
      <n-layout-header v-if="false" bordered class="dark-header">
        <div class="header-content">
          <div style="display: flex; align-items: center; gap: 20px">
            <h2 style="margin: 0; color: #fff">Stock App</h2>
            <n-menu
              mode="horizontal"
              :options="menuOptions"
              :value="route.path"
              @update:value="handleMenuUpdate"
            />
          </div>
          <n-dropdown :options="userOptions" trigger="click" @select="handleUserAction">
            <n-avatar round style="cursor: pointer; background: rgba(99, 102, 241, 0.2)">User</n-avatar>
          </n-dropdown>
        </div>
      </n-layout-header>
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

.dark-header {
  background: #111 !important;
  border-color: #333 !important;
  padding: 12px 16px;
  padding-top: calc(12px + env(safe-area-inset-top));
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 480px) {
  .dark-header {
    padding: 8px 12px;
    padding-top: calc(8px + env(safe-area-inset-top));
  }

  h2 {
    font-size: 16px !important;
  }
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
