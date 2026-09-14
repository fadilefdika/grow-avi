import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/submission',
      name: 'submission',
      component: () => import('../views/employee/SubmissionView.vue'),
      meta: { requiresAuth: false } // Set to false for now so user can preview without logging in
    }
  ]
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  let isAuthenticated = !!authStore.token;

  // Jika tidak terotentikasi di memori, coba lakukan silent refresh (menggunakan cookie HTTPOnly)
  // Ini menangani kasus user menekan F5/Refresh browser.
  if (!isAuthenticated && to.name !== 'login') {
    try {
      const response = await axios.post('/api/auth/refresh', {}, { withCredentials: true });
      if (response.data && response.data.access_token) {
        authStore.setAuth(response.data.access_token, response.data.user);
        isAuthenticated = true;
      }
    } catch (e) {
      // Refresh gagal (cookie tidak valid / kadaluarsa), abaikan dan biarkan ter-redirect ke login
    }
  }

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login');
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router;
