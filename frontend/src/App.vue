<script setup lang="ts">
import { NLayout, NLayoutHeader, NLayoutContent, NMenu, NAvatar, NDropdown, NMessageProvider, NConfigProvider, darkTheme } from 'naive-ui'
import { useRouter, useRoute } from 'vue-router'

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
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <n-layout style="min-height: 100vh" class="dark-layout">
      <n-layout-header bordered class="dark-header">
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
            <n-avatar round style="cursor: pointer; background: #333">User</n-avatar>
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
  background: #000 !important;
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
