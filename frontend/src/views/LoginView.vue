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
          Silakan masuk ke akun Anda
        </p>
      </div>

      <!-- Login Type Toggle -->
      <div class="flex p-1 space-x-1 bg-gray-100/80 rounded-xl">
        <button 
          @click="loginType = 'employee'" 
          :class="['w-1/2 py-2.5 text-sm font-medium rounded-lg transition-all duration-200', loginType === 'employee' ? 'bg-white text-primary shadow-sm' : 'text-gray-500 hover:text-gray-700']"
        >
          Karyawan
        </button>
        <button 
          @click="loginType = 'admin'" 
          :class="['w-1/2 py-2.5 text-sm font-medium rounded-lg transition-all duration-200', loginType === 'admin' ? 'bg-white text-primary shadow-sm' : 'text-gray-500 hover:text-gray-700']"
        >
          Admin
        </button>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <!-- Error Alert -->
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg relative text-sm" role="alert">
          <span class="block sm:inline">{{ error }}</span>
        </div>

        <div class="space-y-4">
          <!-- Username / NPK Field -->
          <div>
            <label :for="loginType === 'employee' ? 'npk' : 'username'" class="block text-sm font-medium text-gray-700 mb-1">
              {{ loginType === 'employee' ? 'NPK' : 'Username' }}
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <input 
                :id="loginType === 'employee' ? 'npk' : 'username'" 
                :name="loginType === 'employee' ? 'npk' : 'username'" 
                :type="loginType === 'employee' ? 'number' : 'text'"
                :inputmode="loginType === 'employee' ? 'numeric' : 'text'"
                :pattern="loginType === 'employee' ? '[0-9]*' : undefined"
                required 
                v-model="usernameInput" 
                class="appearance-none rounded-lg relative block w-full px-3 py-3 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary focus:z-10 sm:text-sm transition-colors" 
                :placeholder="loginType === 'employee' ? 'Masukkan NPK (angka)' : 'Masukkan Username'" 
              />
            </div>
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
              {{ loginType === 'employee' ? 'Password (Awork)' : 'Password' }}
            </label>
            <div class="relative">
              <!-- Ikon gembok di kiri -->
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <input 
                id="password" 
                name="password" 
                :type="showPassword ? 'text' : 'password'" 
                required 
                v-model="passwordInput" 
                class="appearance-none rounded-lg relative block w-full px-3 py-3 pl-10 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary focus:z-10 sm:text-sm transition-colors" 
                :placeholder="loginType === 'employee' ? 'Password Awork' : 'Password Admin'" 
              />
              <!-- Tombol show/hide password (ikon mata) di kanan -->
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition-colors focus:outline-none"
                :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
              >
                <!-- Eye icon (password tersembunyi) -->
                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <!-- Eye-off icon (password terlihat) -->
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.064 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              </button>
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
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';

const router = useRouter();
const authStore = useAuthStore();

const loginType = ref<'employee' | 'admin'>('employee');
const usernameInput = ref('');
const passwordInput = ref('');

// Kosongkan field input setiap kali tab login berganti
watch(loginType, () => {
  usernameInput.value = '';
  passwordInput.value = '';
  error.value = '';
});
const showPassword = ref(false);
const isLoading = ref(false);
const error = ref('');
const countdown = ref(0);
let timerInterval: ReturnType<typeof setInterval> | null = null;

const formatCountdown = (seconds: number) => {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m} menit ${s} detik`;
};

const handleLogin = async () => {
  isLoading.value = true;
  error.value = '';
  
  try {
    const response = await apiClient.post('/auth/login', {
      npk: String(usernameInput.value), // Pastikan selalu string meski input type=number
      password: passwordInput.value,
      type: loginType.value
    });
    
    const { access_token, user } = response.data;
    
    console.log('[DEBUG LOGIN] Data dari backend:', response.data);

    // Save to pinia store
    authStore.setAuth(access_token, user);
    
    // Navigate explicitly based on role
    if (user.role === 'admin') {
      router.push('/admin/dashboard');
    } else {
      router.push('/dashboard');
    }
    
  } catch (err: any) {
    if (err.response) {
      if (err.response.status === 429) {
        const retryAfter = err.response.data.retry_after;
        if (retryAfter && retryAfter > 0) {
          countdown.value = retryAfter;
          error.value = `Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi dalam ${formatCountdown(countdown.value)}.`;
          if (timerInterval) clearInterval(timerInterval);
          timerInterval = setInterval(() => {
            countdown.value--;
            if (countdown.value <= 0) {
              if (timerInterval) clearInterval(timerInterval);
              error.value = ''; // clear error so user knows they can try again
            } else {
              error.value = `Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi dalam ${formatCountdown(countdown.value)}.`;
            }
          }, 1000);
        } else {
          error.value = 'Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi nanti.';
        }
      } else if (err.response.status === 401) {
        error.value = loginType.value === 'employee' ? 'NPK atau Password salah.' : 'Username atau Password salah.';
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
