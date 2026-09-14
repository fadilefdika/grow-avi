<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Top Header -->
    <div class="bg-primary pt-12 pb-24 px-6 rounded-b-[2rem] shadow-md relative">
      <div class="flex justify-between items-center relative z-10">
        <div>
          <button @click="$router.push('/dashboard')" class="text-white mb-2 flex items-center text-sm opacity-80 hover:opacity-100">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali
          </button>
          <h1 class="text-white text-xl font-bold">Katalog Hadiah</h1>
        </div>
        <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm border border-white/30 shadow-inner text-right">
          <p class="text-white/80 text-xs font-medium mb-0.5">Saldo Poin</p>
          <p class="text-white text-lg font-bold">{{ balance }} <span class="text-sm font-normal">pts</span></p>
        </div>
      </div>
    </div>

    <!-- Main Content (Overlapping Header) -->
    <div class="px-5 -mt-16 relative z-20 space-y-5">
      
      <!-- Error / Success Messages -->
      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm relative">
        <span class="block sm:inline">{{ error }}</span>
        <button @click="error = ''" class="absolute top-0 right-0 px-3 py-3">
           <svg class="h-4 w-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
      <div v-if="successMsg" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm relative">
        <span class="block sm:inline">{{ successMsg }}</span>
        <button @click="successMsg = ''" class="absolute top-0 right-0 px-3 py-3">
           <svg class="h-4 w-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Empty State -->
      <div v-else-if="rewards.length === 0" class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 text-center">
        <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
          </svg>
        </div>
        <p class="text-gray-500 font-medium">Belum ada katalog hadiah saat ini.</p>
      </div>

      <!-- Reward Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div v-for="reward in rewards" :key="reward.id" class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 flex flex-col relative overflow-hidden group">
          
          <!-- Sold Out Overlay -->
          <div v-if="reward.stock <= 0" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex items-center justify-center">
            <div class="bg-gray-800 text-white font-bold py-1 px-4 rounded-full shadow-lg transform -rotate-12 border-2 border-white">
              HABIS
            </div>
          </div>

          <div class="flex-1">
            <h3 class="font-bold text-gray-800 text-lg mb-1 leading-tight">{{ reward.title }}</h3>
            <div class="flex items-center gap-2 mb-4">
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-orange-100 text-orange-800">
                Sisa: {{ reward.stock }}
              </span>
            </div>
          </div>
          
          <div class="mt-auto pt-4 border-t border-gray-100 flex items-center justify-between">
            <div class="flex items-baseline gap-1 text-primary">
              <span class="text-2xl font-bold">{{ reward.points_required }}</span>
              <span class="text-xs font-medium">pts</span>
            </div>
            
            <button 
              @click="confirmRedeem(reward)" 
              :disabled="reward.stock <= 0 || balance < reward.points_required || isRedeeming"
              class="px-4 py-2 bg-primary text-white rounded-lg font-medium text-sm transition-colors hover:bg-blue-800 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Tukar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirm Modal -->
    <div v-if="selectedReward" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-sm overflow-hidden">
        <div class="p-5 text-center">
          <div class="w-16 h-16 bg-blue-50 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-1">Konfirmasi Penukaran</h3>
          <p class="text-sm text-gray-500 mb-4">
            Anda akan menukarkan <span class="font-bold text-gray-700">{{ selectedReward.points_required }} pts</span> untuk mendapatkan <span class="font-bold text-gray-700">{{ selectedReward.title }}</span>. Lanjutkan?
          </p>
          
          <div class="grid grid-cols-2 gap-3">
            <button @click="selectedReward = null" class="py-2.5 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 transition-colors text-sm">
              Batal
            </button>
            <button @click="executeRedeem" :disabled="isRedeeming" class="py-2.5 bg-primary text-white font-medium rounded-lg hover:bg-blue-800 transition-colors flex items-center justify-center text-sm disabled:opacity-50">
              <svg v-if="isRedeeming" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              {{ isRedeeming ? 'Memproses' : 'Ya, Tukar' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import apiClient from '../../api/client';

const rewards = ref<any[]>([]);
const balance = ref(0);
const isLoading = ref(true);
const isRedeeming = ref(false);
const error = ref('');
const successMsg = ref('');
const selectedReward = ref<any>(null);

const fetchData = async () => {
  isLoading.value = true;
  error.value = '';
  try {
    const [balRes, rwRes] = await Promise.all([
      apiClient.get('/leaderboard/my-balance'),
      apiClient.get('/rewards')
    ]);
    balance.value = balRes.data.balance || 0;
    rewards.value = rwRes.data.data || [];
  } catch (err) {
    error.value = 'Gagal memuat data katalog hadiah.';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

const confirmRedeem = (reward: any) => {
  error.value = '';
  successMsg.value = '';
  selectedReward.value = reward;
};

const executeRedeem = async () => {
  if (!selectedReward.value) return;
  
  isRedeeming.value = true;
  error.value = '';
  
  try {
    const res = await apiClient.post(`/rewards/redeem/${selectedReward.value.id}`);
    successMsg.value = res.data.message || `Berhasil menukar ${selectedReward.value.title}!`;
    balance.value = res.data.remaining_balance;
    
    // Refresh catalog to update stock
    const rwRes = await apiClient.get('/rewards');
    rewards.value = rwRes.data.data || [];
    
  } catch (err: any) {
    if (err.response && err.response.data && err.response.data.error) {
      error.value = err.response.data.error;
    } else {
      error.value = 'Gagal melakukan penukaran. Coba lagi.';
    }
  } finally {
    isRedeeming.value = false;
    selectedReward.value = null;
  }
};

onMounted(() => {
  fetchData();
});
</script>
