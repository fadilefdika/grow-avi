<template>
  <div class="w-64 bg-white border-r border-gray-200 flex flex-col hidden md:flex shrink-0">
    <div class="p-6">
      <h2 class="text-xl font-bold text-gray-800">GROW Admin</h2>
    </div>
    <nav class="mt-2 flex-1">
      <router-link to="/admin/dashboard" class="block px-6 py-3" :class="isActive('/admin/dashboard')">Dashboard</router-link>
      <router-link to="/admin/submissions" class="block px-6 py-3" :class="isActive('/admin/submissions')">Antrean Approval</router-link>
      <router-link to="/admin/master-data" class="block px-6 py-3" :class="isActive('/admin/master-data')">Master Data</router-link>
      <router-link to="/admin/rewards" class="block px-6 py-3" :class="isActive('/admin/rewards')">Manajemen Hadiah</router-link>
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
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const isActive = (path: string) => {
  return route.path === path 
    ? 'bg-blue-50 text-blue-700 font-medium border-r-4 border-blue-700' 
    : 'text-gray-600 hover:bg-gray-50 hover:text-primary';
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
