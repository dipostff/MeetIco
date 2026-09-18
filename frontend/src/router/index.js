import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/Login.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/Register.vue')
    },
    {
      path: '/book/:username/:slug',
      name: 'booking',
      component: () => import('@/views/public/BookingPage.vue')
    },
    {
      path: '/cancel/:token',
      name: 'cancel',
      component: () => import('@/views/public/CancelPage.vue')
    },
    {
      path: '/dashboard',
      component: () => import('@/views/dashboard/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard/bookings'
        },
        {
          path: 'bookings',
          name: 'dashboard-bookings',
          component: () => import('@/views/dashboard/Bookings.vue')
        },
        {
          path: 'event-types',
          name: 'dashboard-event-types',
          component: () => import('@/views/dashboard/EventTypes.vue')
        },
        {
          path: 'settings',
          name: 'dashboard-settings',
          component: () => import('@/views/dashboard/Settings.vue')
        }
      ]
    },
    {
      path: '/',
      redirect: '/dashboard/bookings'
    }
  ]
})

// Navigation guard for auth
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    next('/dashboard/bookings')
  } else {
    next()
  }
})

export default router
