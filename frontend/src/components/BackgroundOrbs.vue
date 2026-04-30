<script setup lang="ts">
interface Props {
  variant?: 'default' | 'green'
}

withDefaults(defineProps<Props>(), {
  variant: 'default'
})
</script>

<template>
  <div class="page-background">
    <div
      class="gradient-orb orb-1"
      :style="variant === 'green' ? { background: '#10b981' } : { background: '#6366f1' }"
    ></div>
    <div
      class="gradient-orb orb-2"
      :style="variant === 'green' ? { background: '#059669' } : { background: '#8b5cf6' }"
    ></div>
    <div
      v-if="variant === 'default'"
      class="gradient-orb orb-3"
    ></div>
  </div>
</template>

<style scoped>
.orb-1 {
  width: 500px;
  height: 500px;
  top: -150px;
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  bottom: -100px;
  right: -100px;
  animation: float 10s ease-in-out infinite reverse;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: #06b6d4;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: pulse 6s ease-in-out infinite;
}
</style>
