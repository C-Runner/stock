<template>
  <div class="login-container">
    <n-card class="login-card" :bordered="false">
      <template #header>
        <h2>Stock App</h2>
      </template>
      <n-form ref="formRef" :model="form" :rules="rules">
        <n-form-item path="username" label="Username">
          <n-input
            v-model:value="form.username"
            placeholder="Enter username"
            @keydown.enter="handleLogin"
          />
        </n-form-item>
        <n-form-item path="password" label="Password">
          <n-input
            v-model:value="form.password"
            type="password"
            placeholder="Enter password"
            @keydown.enter="handleLogin"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button type="primary" block :loading="loading" @click="handleLogin">Login</n-button>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { api } from '../api'

const router = useRouter()
const message = useMessage()

const formRef = ref()
const loading = ref(false)
const form = ref({
  username: '',
  password: ''
})

const rules = {
  username: { required: true, message: 'Username is required', trigger: 'blur' },
  password: { required: true, message: 'Password is required', trigger: 'blur' }
}

const handleLogin = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const response = await api.post<{ token: string; expiresAt: number }>('/api/login', form.value)
    localStorage.setItem('token', response.token)
    localStorage.setItem('tokenExpiry', String(response.expiresAt))
    message.success('Login successful')
    router.push('/')
  } catch (error) {
    message.error('Invalid credentials')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #1a1a1a !important;
  border: 1px solid #333;
  border-radius: 12px;
}

.login-card h2 {
  text-align: center;
  color: #fff;
  margin: 0;
}

:deep(.n-card-header) {
  padding-bottom: 0;
}

:deep(.n-form-item) {
  margin-bottom: 20px;
}

:deep(.n-input) {
  --n-color: #1a1a1a !important;
  --n-color-focus: #1a1a1a !important;
  --n-border: 1px solid #333 !important;
  --n-border-focus: 1px solid #667eea !important;
  background: #1a1a1a !important;
}
</style>
