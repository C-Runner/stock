import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/watchlist',
    name: 'Watchlist',
    component: () => import('../views/Watchlist.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/analysis/:code',
    name: 'Analysis',
    component: () => import('../views/Analysis.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const expiresAt = localStorage.getItem('tokenExpiry')

  if (to.meta.requiresAuth !== false) {
    if (!token || !expiresAt || Date.now() > Number(expiresAt) * 1000) {
      localStorage.removeItem('token')
      localStorage.removeItem('tokenExpiry')
      next('/login')
      return
    }
  }

  if (to.path === '/login' && token && expiresAt && Date.now() <= Number(expiresAt) * 1000) {
    next('/')
    return
  }

  next()
})

export default router
