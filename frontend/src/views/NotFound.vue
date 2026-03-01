<template>
  <div class="not-found">
    <n-result status="404" title="404 Not Found" description="The page you are looking for does not exist.">
      <template #footer>
        <n-button @click="goHome">Go Home</n-button>
      </template>
    </n-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NResult, NButton } from 'naive-ui'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.not-found {
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}
</style>
