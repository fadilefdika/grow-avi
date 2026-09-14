<template>
  <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 flex justify-around items-center h-16 z-50 pb-safe shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)]">
    <router-link to="/" class="flex flex-col items-center justify-center w-full h-full text-gray-500 hover:text-primary transition-colors" active-class="text-primary font-semibold">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      <span class="text-[10px]">Dashboard</span>
    </router-link>
    
    <router-link to="/submission" class="flex flex-col items-center justify-center w-full h-full text-gray-500 hover:text-primary transition-colors" active-class="text-primary font-semibold">
      <div class="bg-primary text-white p-2 rounded-full -mt-5 border-4 border-gray-50 shadow-md">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </div>
      <span class="text-[10px] mt-1">Pengajuan</span>
    </router-link>

    <div @click="handleLogout" class="flex flex-col items-center justify-center w-full h-full text-gray-500 hover:text-red-500 transition-colors cursor-pointer">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
      </svg>
      <span class="text-[10px]">Keluar</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';

const router = useRouter();
const authStore = useAuthStore();

const handleLogout = async () => {
  try {
    await apiClient.post('/auth/logout');
  } catch (error) {
    console.error("Logout API failed, but clearing local state anyway");
  } finally {
    authStore.logout();
    router.push('/login');
  }
};
</script>
