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
      redirect: () => {
        const authStore = useAuthStore();
        return authStore.user?.role === 'admin' ? '/admin/dashboard' : '/dashboard';
      }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: () => import('../views/admin/AdminDashboardView.vue'),
      meta: { requiresAuth: true, role: 'admin' }
    },
    {
      path: '/submission',
      name: 'submission',
      component: () => import('../views/employee/SubmissionView.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/rewards',
      name: 'rewards',
      component: () => import('../views/employee/RewardView.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: () => import('../views/employee/LeaderboardView.vue'),
      meta: { requiresAuth: true, role: 'employee' }
    },
    {
      path: '/admin/submissions',
      name: 'admin-submissions',
      component: () => import('../views/admin/SubmissionQueueView.vue'),
      meta: { requiresAuth: true, role: 'admin' }
    },
    {
      path: '/admin/master-data',
      name: 'admin-master-data',
      component: () => import('../views/admin/MasterDataView.vue'),
      meta: { requiresAuth: true, role: 'admin' }
    },
    {
      path: '/admin/rewards',
      name: 'admin-rewards',
      component: () => import('../views/admin/RewardManagementView.vue'),
      meta: { requiresAuth: true, role: 'admin' }
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
    next(authStore.user?.role === 'admin' ? '/admin/dashboard' : '/dashboard');
  } else if (to.meta.requiresAuth && to.meta.role) {
    // Check role explicitly
    if (authStore.user?.role !== to.meta.role) {
       // redirect to appropriate dashboard instead of just blocking
       next(authStore.user?.role === 'admin' ? '/admin/dashboard' : '/dashboard');
    } else {
       next();
    }
  } else {
    next();
  }
});

export default router;
