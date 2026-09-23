<template>
  <div
    class="h-screen w-full flex flex-col items-center justify-between overflow-hidden relative"
    style="background: linear-gradient(180deg, #D9F1FF 0%, #3CA5E1 100%);"
  >
    <!-- Background Clouds -->
    <img
      src="/img/awan-login-pojok-kiri.png"
      alt=""
      class="absolute top-0 left-0 w-44 md:w-72 opacity-90 pointer-events-none"
    />
    <img
      src="/img/awan-login-pojok-kanan.png"
      alt=""
      class="absolute top-44 right-0 w-40 md:w-72 opacity-90 pointer-events-none"
    />

    <!-- Header + Card (grup atas) -->
    <div class="relative z-20 flex flex-col items-center w-full px-4 pt-20 md:pt-10 lg:pt-2">
      <img
        src="/img/Judul-GROW.png"
        alt="GROW Logo"
        class="h-40 md:h-28 lg:h-44 mb-4 drop-shadow-md"
      />

      <div class="w-full max-w-sm md:max-w-md bg-white rounded-3xl shadow-2xl p-5 md:p-6">
        <!-- Tab Switcher -->
        <div class="flex p-1 gap-1 bg-gray-100 rounded-full mb-5 border border-gray-200">
          <button
            @click="loginType = 'employee'"
            :class="[
              'flex-1 py-2 text-sm font-bold rounded-full transition-all duration-300',
              loginType === 'employee'
                ? 'bg-[#1E6FA8] text-white shadow-md'
                : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            Karyawan
          </button>
          <button
            @click="loginType = 'admin'"
            :class="[
              'flex-1 py-2 text-sm font-bold rounded-full transition-all duration-300 flex items-center justify-center gap-1.5',
              loginType === 'admin'
                ? 'bg-[#1E6FA8] text-white shadow-md'
                : 'text-[#1E6FA8] hover:bg-blue-50'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            Admin
          </button>
        </div>

        <div class="text-center mb-5">
          <h2 class="text-xl font-black text-[#1E6FA8]">Selamat Datang!</h2>
          <p class="text-gray-500 text-sm font-medium mt-1">Masuk untuk memulai perjalananmu</p>
        </div>

        <!-- Form -->
        <form class="space-y-4" @submit.prevent="handleLogin">
          <!-- NPK Field -->
          <div class="space-y-1.5">
            <label class="block text-sm font-bold text-[#1E6FA8] ml-2">
              {{ loginType === 'employee' ? 'NPK' : 'Username' }}
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <input
                type="text"
                required
                v-model="usernameInput"
                class="w-full pl-11 pr-4 py-3 border border-[#1E6FA8]/40 rounded-full text-gray-900 font-semibold text-sm placeholder-gray-400 focus:outline-none focus:border-[#1E6FA8] focus:ring-1 focus:ring-[#1E6FA8] transition-all bg-white shadow-sm"
                :placeholder="loginType === 'employee' ? 'Masukkan NPK Anda' : 'Masukkan Username'"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div class="space-y-1.5">
            <label class="block text-sm font-bold text-[#1E6FA8] ml-2">
              Password (A-Work)
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <input
                :type="showPassword ? 'text' : 'password'"
                required
                v-model="passwordInput"
                class="w-full pl-11 pr-12 py-3 border border-[#1E6FA8]/40 rounded-full text-gray-900 font-semibold text-sm placeholder-gray-400 focus:outline-none focus:border-[#1E6FA8] focus:ring-1 focus:ring-[#1E6FA8] transition-all bg-white shadow-sm"
                placeholder="Password"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-[#1E6FA8] transition-colors focus:outline-none"
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
            type="submit"
            :disabled="isLoading"
            class="w-full py-3.5 mt-2 font-bold text-white rounded-full bg-[#1E6FA8] hover:bg-[#0D47A1] transition-all duration-300 disabled:opacity-70 flex justify-center items-center shadow-lg hover:shadow-xl"
          >
            <span v-if="isLoading" class="flex items-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>Masuk</span>
          </button>
        </form>
      </div>
    </div>

        <!-- ===================== DEKORASI MOBILE (overlap penuh, solid) ===================== -->
    <div
      class="flex md:hidden relative z-10 w-full h-56 pointer-events-none overflow-visible"
    >
      <!-- kiri: list icon -->
      <img
        src="/img/list-login-pojok-kiri.png"
        alt=""
        class="absolute left-0 -bottom-8 w-32 z-10"
      />

      <!-- raket + kamera (tengah-kiri) -->
      <img
        src="/img/raket-samping-kamera.png"
        alt=""
        class="absolute left-16 -bottom-8 w-52 z-20"
      />

      <!-- helm -->
      <img
        src="/img/helm-login-diatas-raket.png"
        alt=""
        class="absolute left-9 -bottom-2 w-48 z-30"
      />

      <!-- kamera + bola -->
      <img
        src="/img/kamera-login-samping-bola.png"
        alt=""
        class="absolute right-10 -bottom-8 w-56 z-50"
      />

      <!-- buku, menindih kamera -->
      <img
        src="/img/buku-login-diatas-kamera.png"
        alt=""
        class="absolute right-24 -bottom-5 w-52 z-30"
      />

      <!-- bola pojok kanan bawah -->
      <img
        src="/img/bola-login-pojok-kanan-bawah.png"
        alt=""
        class="absolute right-0 -bottom-10 w-28 z-40"
      />
    </div>

    <!-- Gradient overlay biru → transparan, relatif ke SELURUH LAYAR, mentok bawah -->
    <div
      class="absolute bottom-0 left-0 w-full h-32 md:h-40 lg:h-48 xl:h-56 bg-gradient-to-t from-[#0D47A1]/85 via-[#1E6FA8]/40 to-transparent z-[15] pointer-events-none"
    ></div>

    <!-- ===================== DEKORASI DESKTOP (dua cluster pojok) ===================== -->
    <div class="hidden md:block absolute inset-0 pointer-events-none z-10 overflow-hidden">
      <!-- Cluster kiri-bawah -->
      <div class="absolute bottom-0 left-0 flex items-end">
        <img
          src="/img/list-login-pojok-kiri.png"
          alt=""
          class="w-44 lg:w-52 xl:w-60 2xl:w-72 -ml-6 lg:-ml-8 translate-y-4"
        />
        <img
          src="/img/kamera-login-samping-bola.png"
          alt=""
          class="w-32 lg:w-40 xl:w-48 2xl:w-56 -ml-8 lg:-ml-10"
        />
        <img
          src="/img/raket-samping-kamera.png"
          alt=""
          class="w-28 lg:w-36 xl:w-40 2xl:w-48 -ml-6 lg:-ml-8"
        />
      </div>

      <!-- Cluster kanan-bawah -->
      <div class="absolute bottom-0 right-0 flex items-end">
        <img
          src="/img/buku-login-diatas-kamera.png"
          alt=""
          class="w-32 lg:w-40 xl:w-48 2xl:w-56 -mr-2 translate-y-2"
        />
        <img
          src="/img/helm-login-diatas-raket.png"
          alt=""
          class="w-28 lg:w-36 xl:w-44 2xl:w-52 -mr-2"
        />
        <img
          src="/img/bola-login-pojok-kanan-bawah.png"
          alt=""
          class="w-52 lg:w-64 xl:w-72 2xl:w-80 -mr-8 lg:-mr-10 translate-y-4"
        />
      </div>
    </div>

    <!-- Footer -->
    <div class="relative z-20 w-full text-center pb-3">
      <p class="text-white/90 text-[10px] md:text-xs font-semibold">
        © 2026 PT ASTRA VISTEON INDONESIA
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

<style scoped>
/* Keyframes untuk floating animation */
@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
}

@keyframes float-slow {
  0% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-15px) rotate(2deg); }
  100% { transform: translateY(0px) rotate(0deg); }
}

@keyframes float-delay {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-8px); }
  100% { transform: translateY(0px); }
}

@keyframes float-slow-delay {
  0% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-12px) rotate(-2deg); }
  100% { transform: translateY(0px) rotate(0deg); }
}

.animate-float {
  animation: float 4s ease-in-out infinite;
}

.animate-float-slow {
  animation: float-slow 6s ease-in-out infinite;
}

.animate-float-delay {
  animation: float-delay 5s ease-in-out infinite;
  animation-delay: 1.5s;
}

.animate-float-slow-delay {
  animation: float-slow-delay 7s ease-in-out infinite;
  animation-delay: 2s;
}
</style>
