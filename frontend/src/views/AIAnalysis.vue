<template>
  <div class="ai-analysis">
    <BackgroundOrbs />

    <div class="page-header">
      <button class="back-btn" @click="goBack">
        <n-icon size="20"><IconArrowLeft /></n-icon>
      </button>
      <div class="header-info">
        <h1 class="page-title">AI Analysis</h1>
        <p class="subtitle">{{ report?.name || code }} ({{ code }})</p>
      </div>
    </div>

    <n-space vertical :size="12" class="content" :style="{ width: '100%' }">
      <n-spin v-if="loading" show description="Generating AI analysis..." />
      <n-alert v-else-if="error" type="error">{{ error }}</n-alert>

      <template v-else-if="report">
        <!-- Heuristic Mode Notice -->
        <div v-if="report.analysisMethod === 'heuristic'" class="heuristic-notice">
          <n-icon size="16"><IconLight /></n-icon>
          <span>Algorithm analysis (AI not configured)</span>
        </div>

        <!-- Cache Notice -->
        <div v-else-if="report.fromCache" class="cache-notice">
          <n-icon size="16"><IconClock /></n-icon>
          <span>Cached result (refreshes every 5 minutes)</span>
        </div>

        <!-- AI Analysis Card -->
        <n-card v-if="report.rawAnalysis" class="analysis-card raw-analysis-card" :bordered="false" :content-style="{ padding: '20px' }" :header-style="{ padding: '16px 20px', borderBottom: 'none' }" :shadow="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20"><IconRobot /></n-icon>
              <span>AI Analysis</span>
            </div>
          </template>
          <div class="raw-analysis-content" v-html="renderMarkdown(report.rawAnalysis)"></div>
        </n-card>

        <div class="generated-at">
          Generated at: {{ report.generatedAt }}
        </div>
      </template>
    </n-space>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard, NSpace, NSpin, NAlert, NIcon
} from 'naive-ui'
import { aiApi, type AIAnalysisReport } from '../api'
import { IconArrowLeft, IconClock, IconLight, IconRobot } from '../components/icons'
import BackgroundOrbs from '../components/BackgroundOrbs.vue'

const route = useRoute()
const router = useRouter()

const code = route.params.code as string
const report = ref<AIAnalysisReport | null>(null)
const loading = ref(true)
const error = ref('')

const markedInstance = shallowRef<any>(null)

const loadMarkdown = async () => {
  if (markedInstance.value) return
  const markedModule = await import('marked')
  const { markedHighlight } = await import('marked-highlight')
  const hljsModule = await import('highlight.js')
  const hljs = hljsModule.default
  await import('highlight.js/styles/github-dark.css')
  markedModule.marked.use(markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code: string, lang: string) {
      if (lang && hljs.getLanguage(lang)) {
        return hljs.highlight(code, { language: lang }).value
      }
      return hljs.highlightAuto(code).value
    }
  }))
  markedInstance.value = markedModule.marked
}

const fetchAnalysis = async () => {
  loading.value = true
  error.value = ''
  try {
    await loadMarkdown()
    report.value = await aiApi.getAnalysis(code)
  } catch (e: any) {
    error.value = e.message || 'Failed to load AI analysis'
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const renderMarkdown = (text: string): string => {
  if (!text || !markedInstance.value) return ''
  text = text.replace(/<\/?p[^>]*>/gi, '')
  text = text.replace(/(#{1,6})([^\s])/g, '$1 $2')
  text = text.replace(/^\*\*\*\*\*分析报告\s*$/gm, '')
  text = text.replace(/^\*\*\*\*\*\s*$/gm, '')
  text = text.replace(/^\*\*\*技术分析报告\s*$/gm, '')
  text = text.replace(/^\*\*\*\s*技术分析报告\s*$/gm, '')
  const result = markedInstance.value.parse(text, { breaks: true }) as string
  return result.replace(/<h([1-6])>(#+\s*)/g, '<h$1>')
}

onMounted(fetchAnalysis)
</script>

<style scoped>
.ai-analysis {
  width: 100%;
  height: 100%;
  margin: 0 auto;
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px;
  padding-bottom: calc(70px + 6px + env(safe-area-inset-bottom));
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-y;
}

.content {
  position: relative;
  align-items: stretch;
  width: 100%;
}

.content > * {
  flex-shrink: 0;
  margin-bottom: 0 !important;
}

.cache-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(251, 191, 36, 0.1);
  border: 1px solid rgba(251, 191, 36, 0.2);
  border-radius: 12px;
  font-size: 13px;
  color: #fbbf24;
}

.cache-notice :deep(.n-icon) {
  color: #fbbf24;
}

.heuristic-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 12px;
  font-size: 13px;
  color: #818cf8;
}

.heuristic-notice :deep(.n-icon) {
  color: #818cf8;
}

.generated-at {
  text-align: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
  margin-top: 8px;
  padding: 16px;
}

.raw-analysis-card {
  margin-bottom: 16px !important;
}

.raw-analysis-content {
  font-size: 14px;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.85);
}

.raw-analysis-content h2 {
  font-size: 18px;
  margin: 12px 0 8px;
  color: #fff;
}

.raw-analysis-content h3 {
  font-size: 16px;
  margin: 10px 0 6px;
  color: #fff;
}

.raw-analysis-content h4 {
  font-size: 14px;
  margin: 8px 0 4px;
  color: rgba(255, 255, 255, 0.9);
}

.raw-analysis-content p {
  margin: 4px 0;
}

.raw-analysis-content ul {
  margin: 8px 0;
  padding-left: 20px;
}

.raw-analysis-content li {
  margin: 4px 0;
}

.raw-analysis-content strong {
  color: #818cf8;
}

.raw-analysis-content em {
  color: rgba(255, 255, 255, 0.7);
  font-style: italic;
}

.raw-analysis-content .section-hr {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(99, 102, 241, 0.3), transparent);
  margin: 16px 0;
}

.raw-analysis-content .para {
  margin: 0;
  line-height: 1.7;
}

.raw-analysis-content .formatted-content {
  line-height: 1.8;
}

.raw-analysis-content .section-item {
  display: flex;
  padding: 8px 12px;
  margin: 4px 0;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  border-left: 3px solid #6366f1;
}

.raw-analysis-content .section-label {
  color: #818cf8;
  font-weight: 500;
  min-width: 120px;
  flex-shrink: 0;
}

.raw-analysis-content .section-value {
  color: rgba(255, 255, 255, 0.85);
}

.raw-analysis-content .bullet-item {
  display: block;
  padding: 4px 12px;
  margin: 4px 0;
  color: rgba(255, 255, 255, 0.85);
}

.raw-analysis-content .bullet-item::before {
  content: '●';
  color: #6366f1;
  margin-right: 10px;
}

.raw-analysis-content .numbered-item {
  display: block;
  padding: 4px 12px;
  margin: 4px 0;
  color: rgba(255, 255, 255, 0.85);
}

/* Markdown rendering styles */
.raw-analysis-content .md-header {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 16px 0 8px 0;
  padding-left: 0;
}

.raw-analysis-content .md-header::before {
  content: '▎';
  color: #6366f1;
  margin-right: 8px;
}

.raw-analysis-content .md-list {
  margin: 8px 0;
  padding-left: 20px;
  list-style: none;
}

.raw-analysis-content .md-list-item {
  padding: 6px 12px;
  margin: 4px 0;
  color: rgba(255, 255, 255, 0.85);
  position: relative;
}

.raw-analysis-content .md-list-item::before {
  content: '●';
  color: #6366f1;
  position: absolute;
  left: -4px;
}

.raw-analysis-content .md-number-item {
  padding: 6px 12px;
  margin: 4px 0;
  color: rgba(255, 255, 255, 0.85);
}

.raw-analysis-content .md-paragraph {
  margin: 8px 0;
  padding: 8px 12px;
  color: rgba(255, 255, 255, 0.85);
  line-height: 1.7;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
}

.raw-analysis-content .md-bold {
  color: #818cf8;
  font-weight: 600;
}

@media (max-width: 768px) {
  .ai-analysis {
    padding: 12px;
    padding-bottom: calc(80px + env(safe-area-inset-bottom));
  }
}
</style>