<template>
  <div class="min-h-screen relative overflow-hidden flex items-center justify-center p-4" style="background: linear-gradient(135deg, #002f4e 0%, #004f82 35%, #0069AA 65%, #0090d4 100%);">
    
    <!-- Animated Background Shapes -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-20 -left-20 w-80 h-80 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #1CB0F6, transparent); animation-delay: 0s;"></div>
      <div class="absolute top-1/4 -right-16 w-64 h-64 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFA64D, transparent); animation-delay: 1.5s;"></div>
      <div class="absolute -bottom-20 left-1/4 w-96 h-96 rounded-full opacity-8 animate-float" style="background: radial-gradient(circle, #0069AA, transparent); animation-delay: 0.8s;"></div>
      <!-- Decorative dots grid -->
      <div class="absolute top-8 right-8 grid grid-cols-4 gap-3 opacity-20">
        <div v-for="i in 16" :key="i" class="w-1.5 h-1.5 rounded-full bg-white"></div>
      </div>
      <div class="absolute bottom-12 left-8 grid grid-cols-4 gap-3 opacity-20">
        <div v-for="i in 12" :key="i" class="w-1.5 h-1.5 rounded-full bg-white"></div>
      </div>
    </div>

    <!-- Login Card -->
    <div class="relative w-full max-w-sm animate-pop">
      <!-- Logo & Brand Header -->
      <div class="text-center mb-6">
        
        <h1 class="text-2xl font-black text-white tracking-tight">
          GROW <span class="text-secondary">AVI</span>
        </h1>
        <p class="text-white/70 text-sm mt-1 font-medium">
          {{ loginType === 'employee' ? 'Kumpulkan poin, raih hadiahmu! 🏆' : 'Panel Administrator' }}
        </p>
      </div>

      <!-- Main Card -->
      <div class="bg-white rounded-3xl shadow-2xl overflow-hidden">
        <!-- Tab Toggle -->
        <div class="flex p-2 gap-1 bg-gray-50 border-b border-gray-100">
          <button
            id="tab-employee"
            @click="loginType = 'employee'"
            :class="[
              'flex-1 py-2.5 text-sm font-bold rounded-xl transition-all duration-300',
              loginType === 'employee'
                ? 'bg-primary text-white shadow-md'
                : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'
            ]"
          >
            👤 Karyawan
          </button>
          <button
            id="tab-admin"
            @click="loginType = 'admin'"
            :class="[
              'flex-1 py-2.5 text-sm font-bold rounded-xl transition-all duration-300',
              loginType === 'admin'
                ? 'bg-gray-800 text-white shadow-md'
                : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'
            ]"
          >
            🔐 Admin
          </button>
        </div>

        <!-- Form -->
        <form class="p-6 space-y-5" @submit.prevent="handleLogin">
          <!-- Welcome text -->
          <div class="text-center py-1">
            <h2 class="text-xl font-black text-gray-900">Selamat Datang!</h2>
            <p class="text-gray-500 text-sm font-medium mt-0.5">Masuk untuk mulai perjalananmu</p>
          </div>

          <!-- NPK / Username Field -->
          <div class="space-y-1.5">
            <label :for="loginType === 'employee' ? 'npk' : 'username'" class="block text-sm font-bold text-gray-700">
              {{ loginType === 'employee' ? 'NPK' : 'Username' }}
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <input
                :id="loginType === 'employee' ? 'npk' : 'username'"
                :name="loginType === 'employee' ? 'npk' : 'username'"
                type="text"
                :inputmode="loginType === 'employee' ? 'numeric' : 'text'"
                :pattern="loginType === 'employee' ? '[0-9]*' : undefined"
                required
                v-model="usernameInput"
                class="w-full pl-11 pr-4 py-3.5 border-2 border-gray-200 rounded-2xl text-gray-900 font-semibold text-sm placeholder-gray-400 focus:outline-none focus:border-primary transition-all duration-200 bg-gray-50 focus:bg-white [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                :placeholder="loginType === 'employee' ? 'Masukkan NPK Anda' : 'Masukkan Username'"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div class="space-y-1.5">
            <label for="password" class="block text-sm font-bold text-gray-700">
              {{ loginType === 'employee' ? 'Password (Awork)' : 'Password' }}
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <input
                id="password"
                name="password"
                :type="showPassword ? 'text' : 'password'"
                required
                v-model="passwordInput"
                class="w-full pl-11 pr-12 py-3.5 border-2 border-gray-200 rounded-2xl text-gray-900 font-semibold text-sm placeholder-gray-400 focus:outline-none focus:border-primary transition-all duration-200 bg-gray-50 focus:bg-white"
                :placeholder="loginType === 'employee' ? 'Password Awork Anda' : 'Password Admin'"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-primary transition-colors focus:outline-none"
                :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
              >
                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.064 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            id="btn-login"
            type="submit"
            :disabled="isLoading"
            :class="[
              'w-full py-4 font-black text-base rounded-2xl transition-all duration-200 text-white disabled:opacity-60 disabled:cursor-not-allowed',
              loginType === 'employee' ? 'btn-game-primary' : 'bg-gray-800 btn-game'
            ]"
            :style="loginType !== 'employee' ? 'box-shadow: 0 6px 0 #111827;' : ''"
          >
            <span v-if="isLoading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>
              {{ loginType === 'employee' ? 'Masuk & Mulai Kumpulkan Poin 🚀' : 'Masuk sebagai Admin' }}
            </span>
          </button>
        </form>
      </div>

      <!-- Footer branding -->
      <p class="text-center text-white/50 text-xs font-medium mt-5">
        © 2026 PT Astra Visteon Indonesia
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';
import { useToast } from 'vue-toastification';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const loginType = ref<'employee' | 'admin'>('employee');
const usernameInput = ref('');
const passwordInput = ref('');

// Kosongkan field input setiap kali tab login berganti
watch(loginType, () => {
  usernameInput.value = '';
  passwordInput.value = '';
});
const showPassword = ref(false);
const isLoading = ref(false);
const countdown = ref(0);

const formatCountdown = (seconds: number) => {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m} menit ${s} detik`;
};

const handleLogin = async () => {
  isLoading.value = true;
  
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
          toast.error(`Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi dalam ${formatCountdown(countdown.value)}.`);
        } else {
          toast.error('Akun terkunci karena terlalu banyak percobaan gagal. Coba lagi nanti.');
        }
      } else if (err.response.status === 401) {
        toast.error(loginType.value === 'employee' ? 'NPK atau Password salah.' : 'Username atau Password salah.');
      } else {
        toast.error(err.response.data.error || 'Terjadi kesalahan pada server.');
      }
    } else {
      toast.error('Tidak dapat terhubung ke server. Periksa koneksi Anda.');
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped></style>

