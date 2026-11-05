import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const LandingView = () => import('../views/LandingView.vue')
const LoginView = () => import('../views/LoginView.vue')
const DashboardView = () => import('../views/DashboardView.vue')
const AdminView = () => import('../views/AdminView.vue')
const PasienView = () => import('../views/PasienView.vue')
const RekamMedisView = () => import('../views/RekamMedisView.vue')
const KlaimView = () => import('../views/KlaimView.vue')
const MasterICDView = () => import('../views/MasterICDView.vue')

const routes = [
  {
    path: '/',
    name: 'landing',
    component: LandingView
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminView,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/pasien',
    name: 'pasien',
    component: PasienView,
    meta: { requiresAuth: true }
  },
  {
    path: '/rekam-medis',
    name: 'rekam-medis',
    component: RekamMedisView,
    meta: { requiresAuth: true }
  },
  {
    path: '/klaim',
    name: 'klaim',
    component: KlaimView,
    meta: { requiresAuth: true }
  },
  {
    path: '/master-icd',
    name: 'master-icd',
    component: MasterICDView,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated && to.meta.requiresAuth) {
    await authStore.fetchMe()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/dashboard')
  } else if (to.meta.requiresAdmin && (!authStore.isAuthenticated || authStore.user?.role !== 'admin_ti')) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router