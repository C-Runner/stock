<template>
  <div class="settings">
    <div class="background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>

    <div class="settings-header">
      <div class="header-left">
        <button class="back-btn" @click="goBack">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          <span>Back</span>
        </button>
        <div class="header-text">
          <h1>Settings</h1>
        </div>
      </div>
    </div>

    <div class="settings-content">
      <section class="settings-section">
        <div class="section-header">
          <div class="section-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"/>
              <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
            </svg>
          </div>
          <div class="section-title">
            <h2>AI LLM Settings</h2>
            <p class="section-desc">Configure the AI language model for stock analysis</p>
          </div>
        </div>

        <div v-if="loading" class="loading">
          <div class="spinner"></div>
          <p>Loading settings...</p>
        </div>

        <form v-else @submit.prevent="saveSettings" class="settings-form">
          <div class="form-group">
            <label for="apiKey">API Key</label>
            <div class="input-with-hint">
              <input
                type="password"
                id="apiKey"
                v-model="form.apiKey"
                placeholder="Enter API key"
                autocomplete="off"
              />
              <span class="hint">Leave empty to keep current key. Enter new value to update.</span>
            </div>
            <div v-if="settings.hasApiKey" class="current-key">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
              Current key: ********
            </div>
          </div>

          <div class="form-group">
            <label for="apiUrl">API URL</label>
            <input
              type="text"
              id="apiUrl"
              v-model="form.apiUrl"
              placeholder="https://api.anthropic.com/v1/messages"
              autocomplete="off"
            />
            <span class="hint">API endpoint URL for your LLM provider</span>
          </div>

          <div class="form-group">
            <label for="model">Model</label>
            <select id="model" v-model="form.model">
              <option value="MiniMax-M2.7">MiniMax-M2.7</option>
              <option value="DeepSeek-V4-Pro">DeepSeek-V4-Pro</option>
              <option value="deepseek-v4-flash">DeepSeek-V4-Flash</option>
            </select>
            <span class="hint">Select the AI model to use for analysis</span>
          </div>

          <div class="form-group toggle-group">
            <div class="toggle-label">
              <label for="enabled">Enable AI Analysis</label>
              <span class="hint">When disabled, heuristic analysis will be used</span>
            </div>
            <label class="toggle">
              <input type="checkbox" id="enabled" v-model="form.enabled" />
              <span class="slider"></span>
            </label>
          </div>

          <div class="form-actions">
            <button type="submit" class="save-btn" :disabled="saving">
              {{ saving ? 'Saving...' : 'Save Settings' }}
            </button>
            <button type="button" class="test-btn" :disabled="saving || testing" @click="testAI">
              {{ testing ? 'Testing...' : 'Test AI' }}
            </button>
            <span v-if="saveSuccess" class="success-msg">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 6L9 17l-5-5"/>
              </svg>
              Settings saved successfully!
            </span>
            <span v-if="saveError" class="error-msg">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              {{ saveError }}
            </span>
            <span v-if="testSuccess" class="test-success-msg">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 6L9 17l-5-5"/>
              </svg>
              {{ testSuccess }}
            </span>
            <span v-if="testError" class="test-error-msg">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              {{ testError }}
            </span>
          </div>
        </form>
      </section>
    </div>

    <div class="bottom-tabs">
      <div class="tab-item" @click="router.push('/watchlist')">
        <div class="tab-icon">
          <n-icon size="22"><IconSearch /></n-icon>
        </div>
        <span class="tab-label">Watchlist</span>
      </div>
      <div class="tab-item" @click="router.push('/home')">
        <div class="tab-icon">
          <n-icon size="22"><IconWallet /></n-icon>
        </div>
        <span class="tab-label">Holdings</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { settingsApi, type AISettingsResponse, type AISettingsRequest } from '../api'
import { IconSearch, IconWallet } from '../components/icons'

const router = useRouter()

const loading = ref(true)
const saving = ref(false)
const testing = ref(false)
const saveSuccess = ref(false)
const saveError = ref('')
const testSuccess = ref('')
const testError = ref('')

const settings = ref<AISettingsResponse>({
  apiKey: '',
  apiUrl: 'https://api.minimax.chat/v1/text/chatcompletion_v2',
  model: 'MiniMax-Text-01',
  enabled: true,
  hasApiKey: false
})

const form = ref({
  apiKey: '',
  apiUrl: 'https://api.minimax.chat/v1/text/chatcompletion_v2',
  model: 'MiniMax-Text-01',
  enabled: true
})

const loadSettings = async () => {
  loading.value = true
  try {
    const data = await settingsApi.getAISettings()
    settings.value = data
    form.value = {
      apiKey: '',
      apiUrl: data.apiUrl,
      model: data.model,
      enabled: data.enabled
    }
  } catch (e: any) {
    console.error('Failed to load settings:', e)
    saveError.value = 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  saving.value = true
  saveSuccess.value = false
  saveError.value = ''
  testSuccess.value = ''
  testError.value = ''

  try {
    const request: AISettingsRequest = {
      apiUrl: form.value.apiUrl,
      model: form.value.model,
      enabled: form.value.enabled
    }

    if (form.value.apiKey) {
      request.apiKey = form.value.apiKey
    }

    const data = await settingsApi.updateAISettings(request)
    settings.value = data
    form.value.apiKey = ''
    saveSuccess.value = true
    setTimeout(() => { saveSuccess.value = false }, 3000)
  } catch (e: any) {
    console.error('Failed to save settings:', e)
    saveError.value = e.message || 'Failed to save settings'
  } finally {
    saving.value = false
  }
}

const testAI = async () => {
  testing.value = true
  testSuccess.value = ''
  testError.value = ''
  saveSuccess.value = false
  saveError.value = ''

  try {
    const result = await settingsApi.testAISettings()
    if (result.success) {
      testSuccess.value = result.message || 'AI connection successful'
      setTimeout(() => { testSuccess.value = '' }, 5000)
    } else {
      testError.value = result.error || 'AI test failed'
    }
  } catch (e: any) {
    console.error('AI test failed:', e)
    testError.value = e.message || 'Failed to test AI settings'
  } finally {
    testing.value = false
  }
}

const goBack = () => {
  router.back()
}

onMounted(loadSettings)
</script>

<style scoped>
.settings {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  height: 100dvh;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
  padding: 16px;
  padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  display: flex;
  flex-direction: column;
}

.background {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  z-index: -1;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: #6366f1;
  top: -150px;
  left: -100px;
  animation: float 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: #8b5cf6;
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

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}

@keyframes pulse {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.3; }
  50% { transform: translate(-50%, -50%) scale(1.2); opacity: 0.5; }
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  position: relative;
  flex: 0 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.3);
  color: #667eea;
}

.back-btn svg {
  width: 16px;
  height: 16px;
}

.header-text h1 {
  margin: 0;
  font-size: 26px;
  font-weight: 600;
  color: #fff;
  letter-spacing: -0.5px;
}

.settings-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  position: relative;
}

.settings-section {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  padding: 24px;
  box-sizing: border-box;
}

.section-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.section-icon {
  width: 44px;
  height: 44px;
  min-width: 44px;
  padding: 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.section-icon svg {
  width: 100%;
  height: 100%;
  color: #fff;
}

.section-title h2 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.section-desc {
  margin: 0;
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px;
  gap: 16px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
}

.input-with-hint {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group select {
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  font-size: 14px;
  color: #fff;
  transition: all 0.2s ease;
}

.form-group input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15);
  background: rgba(255, 255, 255, 0.08);
}

.form-group select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='rgba(255,255,255,0.5)' stroke-width='2'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  padding-right: 40px;
}

.form-group select option {
  background: rgba(20, 19, 60, 0.95);
  color: #fff;
}

.hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.current-key {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #38ef7d;
  background: rgba(56, 239, 125, 0.1);
  padding: 6px 10px;
  border-radius: 8px;
  width: fit-content;
}

.current-key svg {
  width: 14px;
  height: 14px;
}

.toggle-group {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.toggle-label {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.toggle-label .hint {
  margin-top: 0;
}

.toggle {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 26px;
  flex-shrink: 0;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.1);
  transition: 0.3s;
  border-radius: 26px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 2px;
  bottom: 2px;
  background-color: rgba(255, 255, 255, 0.6);
  transition: 0.3s;
  border-radius: 50%;
}

.toggle input:checked + .slider {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-color: rgba(99, 102, 241, 0.3);
}

.toggle input:checked + .slider:before {
  transform: translateX(22px);
  background-color: #fff;
}

.form-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 8px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.save-btn {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3);
}

.save-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.test-btn {
  padding: 12px 24px;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.test-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
}

.test-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.success-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #38ef7d;
  font-size: 13px;
}

.success-msg svg {
  width: 16px;
  height: 16px;
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ff6b6b;
  font-size: 13px;
}

.error-msg svg {
  width: 16px;
  height: 16px;
}

.test-success-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #38ef7d;
  font-size: 13px;
}

.test-success-msg svg {
  width: 16px;
  height: 16px;
}

.test-error-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ff6b6b;
  font-size: 13px;
}

.test-error-msg svg {
  width: 16px;
  height: 16px;
}

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
  gap: 24px;
  padding: 0 24px;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
  touch-action: none;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 10px 20px;
  border-radius: 16px;
  transition: all 0.2s ease;
  color: rgba(255, 255, 255, 0.6);
}

.tab-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-label {
  font-size: 12px;
  font-weight: 600;
}

@media (max-width: 768px) {
  .settings {
    padding: 16px 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }

  .settings-section {
    padding: 20px 16px;
  }

  .section-header {
    flex-direction: column;
    gap: 12px;
  }

  .section-icon {
    width: 40px;
    height: 40px;
    padding: 8px;
  }

  .form-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .save-btn {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .settings {
    padding: 12px 8px;
    padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  }

  .header-text h1 {
    font-size: 20px;
  }

  .settings-section {
    padding: 16px 12px;
  }
}
</style>
