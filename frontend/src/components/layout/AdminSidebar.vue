<template>
  <!-- Mobile Menu Button -->
  <button 
    @click="isOpen = true" 
    class="md:hidden fixed top-6 right-6 z-40 bg-white border border-gray-200 text-gray-800 p-2.5 rounded-xl shadow-sm hover:bg-gray-50 transition-colors focus:outline-none"
  >
    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
  </button>

  <!-- Mobile Overlay -->
  <Transition name="fade">
    <div v-if="isOpen" @click="isOpen = false" class="md:hidden fixed inset-0 z-40 bg-gray-900/50 backdrop-blur-sm"></div>
  </Transition>

  <!-- Sidebar -->
  <div :class="[
    'fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200 flex flex-col transform transition-transform duration-300 ease-in-out md:relative md:translate-x-0 shrink-0',
    isOpen ? 'translate-x-0' : '-translate-x-full'
  ]">
    <div class="p-6 flex justify-between items-center">
      <h2 class="text-xl font-bold text-gray-800">GROW Admin</h2>
      <button @click="isOpen = false" class="md:hidden text-gray-400 hover:text-gray-600 focus:outline-none">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>
    <nav class="mt-2 flex-1 overflow-y-auto">
      <router-link @click="isOpen = false" to="/admin/dashboard" class="block px-6 py-3" :class="isActive('/admin/dashboard', ['/admin/leaderboard'])">Dashboard</router-link>
      <router-link @click="isOpen = false" to="/admin/submissions" class="flex items-center justify-between px-6 py-3" :class="isActive('/admin/submissions')">
        <span>Antrean Approval</span>
        <span v-if="pendingCount > 0" class="bg-red-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full shadow-sm">{{ pendingCount > 99 ? '99+' : pendingCount }}</span>
      </router-link>
      <router-link @click="isOpen = false" to="/admin/activity-log" class="block px-6 py-3" :class="isActive('/admin/activity-log')">Activity Log</router-link>
      <router-link @click="isOpen = false" to="/admin/master-data" class="block px-6 py-3" :class="isActive('/admin/master-data')">Master Data</router-link>
      <router-link @click="isOpen = false" to="/admin/rewards" class="block px-6 py-3" :class="isActive('/admin/rewards')">Manajemen Hadiah</router-link>
    </nav>
    <div class="p-4 border-t border-gray-200">
       <button @click="handleLogout" class="w-full text-left px-2 py-2 text-sm text-red-500 hover:text-red-700 font-medium flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Keluar
       </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const isOpen = ref(false);
const pendingCount = ref(0);

const fetchPendingCount = async () => {
  try {
    const [resGrow, resClaim] = await Promise.all([
      apiClient.get('/admin/submissions/pending', { params: { limit: 1 } }),
      apiClient.get('/admin/redemptions/pending', { params: { limit: 1 } })
    ]);
    pendingCount.value = (resGrow.data.total || 0) + (resClaim.data.total || 0);
  } catch (err) {
    console.error("Gagal memuat badge pending", err);
  }
};

onMounted(() => {
  if (authStore.user?.role === 'admin') {
    fetchPendingCount();
    // Refresh count periodically every 1 menit
    setInterval(fetchPendingCount, 60000);
  }
});

const isActive = (path: string, relatedPaths: string[] = []) => {
  const isMatch = route.path === path || route.path.startsWith(path + '/') || relatedPaths.includes(route.path);
  return isMatch
    ? 'bg-blue-50 text-blue-700 font-medium border-r-4 border-blue-700' 
    : 'text-gray-600 hover:bg-gray-50 hover:text-primary transition-colors';
};

const handleLogout = async () => {
  try {
    await apiClient.post('/auth/logout');
  } catch (err) {
    console.error('Logout error:', err);
  } finally {
    authStore.logout();
    router.push('/login');
  }
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
