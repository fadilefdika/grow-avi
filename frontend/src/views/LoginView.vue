<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Decorative background elements -->
    <div class="absolute top-0 left-0 w-full h-64 bg-primary rounded-b-[40px] shadow-lg transform -translate-y-10 sm:-translate-y-20"></div>
    
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-2xl shadow-xl z-10 border border-gray-100">
      <div>
        <div class="flex justify-center mb-6">
          <img src="/img/header-grow-avi.png" alt="Astra Visteon GROW" class="h-16 object-contain" />
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
          Selamat Datang
        </h2>
        <p class="mt-2 text-center text-sm text-gray-500">
          Silakan masuk menggunakan NPK Anda
        </p>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <!-- Error Alert -->
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg relative text-sm" role="alert">
          <span class="block sm:inline">{{ error }}</span>
        </div>

        <div class="space-y-4">
          <div>
            <label for="npk" class="block text-sm font-medium text-gray-700 mb-1">NPK</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <input id="npk" name="npk" type="text" required v-model="npk" class="appearance-none rounded-lg relative block w-full px-3 py-3 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary focus:z-10 sm:text-sm transition-colors" placeholder="Masukkan NPK" />
            </div>
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password (Awork)</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <input id="password" name="password" type="password" required v-model="password" class="appearance-none rounded-lg relative block w-full px-3 py-3 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary focus:z-10 sm:text-sm transition-colors" placeholder="Password Awork" />
            </div>
          </div>
        </div>

        <div>
          <button type="submit" :disabled="isLoading" class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-semibold rounded-lg text-white bg-primary hover:bg-blue-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="isLoading">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>Masuk Aplikasi</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';

const router = useRouter();
const authStore = useAuthStore();

const npk = ref('');
const password = ref('');
const isLoading = ref(false);
const error = ref('');

const handleLogin = async () => {
  isLoading.value = true;
  error.value = '';
  
  try {
    const response = await apiClient.post('/auth/login', {
      npk: npk.value,
      password: password.value
    });
    
    const { access_token, user } = response.data;
    
    console.log('[DEBUG LOGIN] Data dari backend:', response.data);

    // Save to pinia store
    authStore.setAuth(access_token, user);
    
    // Navigate explicitly to dashboard
    router.push('/dashboard');
    
  } catch (err: any) {
    if (err.response) {
      if (err.response.status === 429) {
        error.value = 'Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi dalam 15 menit.';
      } else if (err.response.status === 401) {
        error.value = 'NPK atau Password salah.';
      } else {
        error.value = err.response.data.error || 'Terjadi kesalahan pada server.';
      }
    } else {
      error.value = 'Tidak dapat terhubung ke server. Periksa koneksi Anda.';
    }
  } finally {
    isLoading.value = false;
  }
};
</script>
