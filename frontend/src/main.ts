import { createApp } from 'vue'
import router from './router'
import './style.css'
import App from './App.vue'

const app = createApp(App)

app.use(router)

// Apply dark theme globally
const meta = document.createElement('meta')
meta.name = 'theme-color'
meta.content = '#0a0a0f'
document.head.appendChild(meta)

app.mount('#app')
