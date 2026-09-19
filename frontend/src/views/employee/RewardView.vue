<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Hero Header -->
    <div class="hero-gradient pt-12 pb-32 px-6 relative overflow-hidden">
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-8 -right-8 w-40 h-40 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFA64D, transparent);"></div>
        <div class="absolute bottom-0 left-1/3 w-24 h-24 rounded-full opacity-10" style="background: radial-gradient(circle, #fff, transparent);"></div>
      </div>
      <div class="flex justify-between items-start relative z-10">
        <div>
          <button @click="$router.push('/dashboard')" class="flex items-center gap-1.5 text-white/70 text-sm font-bold mb-3 hover:text-white transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali
          </button>
          <h1 class="text-white text-2xl font-black">🎁 Katalog Hadiah</h1>
          <p class="text-white/60 text-xs font-semibold mt-0.5">Tukarkan poinmu dengan hadiah menarik!</p>
        </div>
        <!-- Balance display -->
        <div class="glass-card rounded-2xl p-3 text-right">
          <p class="text-white/70 text-xs font-bold uppercase tracking-wide mb-0.5">Saldo Poin</p>
          <div class="flex items-baseline gap-1 justify-end">
            <span class="text-white text-xl font-black">{{ balance }}</span>
            <span class="text-white/70 text-xs font-bold">pts</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="px-4 -mt-24 relative z-20 space-y-4">

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Empty State -->
      <div v-else-if="rewards.length === 0" class="bg-white rounded-3xl p-10 shadow-sm border border-gray-100 text-center">
        <p class="text-5xl mb-4">🎪</p>
        <p class="font-black text-gray-700 text-base">Katalog hadiah sedang kosong.</p>
        <p class="text-gray-400 text-sm font-semibold mt-1">Cek lagi nanti ya!</p>
      </div>

      <!-- Reward Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div
          v-for="reward in rewards"
          :key="reward.id"
          :class="[
            'relative rounded-3xl overflow-hidden transition-all duration-300 flex flex-col',
            balance >= reward.points_required && reward.stock > 0
              ? 'card-eligible shadow-lg hover:shadow-xl'
              : 'bg-white shadow-sm border border-gray-100 hover:shadow-md'
          ]"
        >
          <!-- Sold Out Overlay -->
          <div v-if="reward.stock <= 0" class="absolute inset-0 bg-white/80 backdrop-blur-[2px] z-10 flex items-center justify-center rounded-3xl">
            <div class="bg-gray-800 text-white font-black py-2 px-5 rounded-full shadow-xl transform -rotate-6 text-sm tracking-widest">
              HABIS
            </div>
          </div>

          <!-- Eligible shimmer effect -->
          <div v-if="balance >= reward.points_required && reward.stock > 0" class="absolute inset-0 pointer-events-none overflow-hidden rounded-3xl">
            <div class="absolute inset-0 opacity-5" style="background: linear-gradient(90deg, transparent 0%, #58CC02 50%, transparent 100%); background-size: 200% 100%; animation: shimmer 3s linear infinite;"></div>
          </div>

          <div class="p-5 flex-1 flex flex-col">
            <!-- Reward icon area -->
            <div class="flex items-start justify-between mb-4">
              <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-orange-100 to-orange-50 flex items-center justify-center shadow-sm border border-orange-100">
                <span class="text-2xl">
                  {{ getRewardEmoji(reward.title) }}
                </span>
              </div>
              <div v-if="balance >= reward.points_required && reward.stock > 0" class="flex items-center gap-1 bg-xp-green/10 text-xp-green px-2 py-1 rounded-xl">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                <span class="text-[10px] font-black">Bisa ditukar!</span>
              </div>
            </div>

            <h3 class="font-black text-gray-800 text-base mb-1 leading-snug">{{ reward.title }}</h3>
            
            <div class="flex-1"></div>

            <div class="flex items-center justify-between pt-4 border-t border-gray-100 mt-4">
              <div class="flex items-baseline gap-1">
                <span class="text-2xl font-black text-primary">{{ reward.points_required }}</span>
                <span class="text-xs font-bold text-gray-400">pts</span>
              </div>
              
              <button
                @click="confirmRedeem(reward)"
                :disabled="reward.stock <= 0 || balance < reward.points_required || isRedeeming"
                :class="[
                  'px-5 py-2.5 rounded-xl font-black text-sm transition-all duration-200',
                  reward.stock > 0 && balance >= reward.points_required
                    ? 'btn-game-orange text-white'
                    : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                ]"
              >
                Tukar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirm Modal -->
    <div v-if="selectedReward" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-t-3xl sm:rounded-3xl shadow-2xl w-full max-w-sm overflow-hidden animate-slide-up">
        <div class="p-6 text-center">
          <!-- Icon -->
          <div class="w-20 h-20 rounded-3xl bg-orange-50 border-2 border-orange-100 flex items-center justify-center mx-auto mb-4">
            <span class="text-4xl">{{ getRewardEmoji(selectedReward.title) }}</span>
          </div>
          
          <h3 class="text-xl font-black text-gray-900 mb-1">Konfirmasi Penukaran</h3>
          <p class="text-sm text-gray-500 font-semibold mb-1">
            Tukar <span class="font-black text-primary">{{ selectedReward.points_required }} pts</span> untuk
          </p>
          <p class="font-black text-gray-800 text-base mb-5">{{ selectedReward.title }}</p>

          <!-- Point Selection -->
          <div class="text-left mb-5">
            <p class="text-xs font-bold text-gray-500 uppercase tracking-wide mb-2">Pilih Poin ID GROW</p>
            <div class="max-h-48 overflow-y-auto space-y-2 pr-1 mb-3 custom-scrollbar">
              <label v-for="pt in availablePoints" :key="pt.id" class="flex items-center justify-between p-3 rounded-xl border cursor-pointer transition-colors" :class="selectedSubmissionIds.has(pt.id) ? 'bg-blue-50 border-primary' : 'bg-white border-gray-200 hover:border-gray-300'">
                <div class="flex items-center gap-3">
                  <input type="checkbox" :checked="selectedSubmissionIds.has(pt.id)" @change="togglePoint(pt.id)" class="w-4 h-4 text-primary bg-gray-100 border-gray-300 rounded focus:ring-primary" />
                  <div>
                    <p class="text-xs font-bold" :class="selectedSubmissionIds.has(pt.id) ? 'text-primary' : 'text-gray-900'">{{ pt.grow_id }}</p>
                    <p class="text-[10px] font-semibold text-gray-500 truncate max-w-[150px]">{{ pt.activity_name }}</p>
                  </div>
                </div>
                <span class="text-sm font-black" :class="selectedSubmissionIds.has(pt.id) ? 'text-primary' : 'text-gray-900'">{{ pt.points_awarded }} pts</span>
              </label>
              <div v-if="availablePoints.length === 0" class="text-center py-4 bg-gray-50 rounded-xl">
                <p class="text-xs font-semibold text-gray-500">Tidak ada tiket poin yang tersedia.</p>
              </div>
            </div>
            
            <div class="bg-gray-50 rounded-xl p-3 flex justify-between items-center border border-gray-100">
              <span class="text-xs font-bold text-gray-500">Total Terpilih:</span>
              <span class="text-sm font-black" :class="totalSelectedPoints >= selectedReward.points_required ? 'text-green-600' : 'text-red-500'">{{ totalSelectedPoints }} / {{ selectedReward.points_required }} pts</span>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-3">
            <button @click="selectedReward = null" class="py-3.5 bg-gray-100 text-gray-700 font-black rounded-2xl hover:bg-gray-200 transition-colors text-sm">
              Batal
            </button>
            <button @click="executeRedeem" :disabled="isRedeeming || totalSelectedPoints < selectedReward.points_required" class="btn-game-orange py-3.5 text-sm flex items-center justify-center gap-2 disabled:opacity-60">
              <svg v-if="isRedeeming" class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              {{ isRedeeming ? 'Memproses...' : 'Tukar! 🎉' }}
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
import { useToast } from 'vue-toastification';

const rewards = ref<any[]>([]);
const balance = ref(0);
const isLoading = ref(true);
const isRedeeming = ref(false);
const selectedReward = ref<any>(null);
const availablePoints = ref<any[]>([]);
const selectedSubmissionIds = ref<Set<number>>(new Set());
const toast = useToast();

import { computed } from 'vue';

const totalSelectedPoints = computed(() => {
  return availablePoints.value
    .filter(p => selectedSubmissionIds.value.has(p.id))
    .reduce((sum, p) => sum + p.points_awarded, 0);
});

const togglePoint = (id: number) => {
  if (selectedSubmissionIds.value.has(id)) {
    selectedSubmissionIds.value.delete(id);
  } else {
    selectedSubmissionIds.value.add(id);
  }
};

const getRewardEmoji = (title: string): string => {
  const t = title?.toLowerCase() || '';
  if (t.includes('grab food') || t.includes('grabfood')) return '🍔';
  if (t.includes('grab')) return '🚗';
  if (t.includes('gojek') || t.includes('go-jek')) return '🛵';
  if (t.includes('shopee')) return '🛍️';
  if (t.includes('tokopedia')) return '🟢';
  if (t.includes('voucher') || t.includes('gift')) return '🎫';
  if (t.includes('makanan') || t.includes('makan') || t.includes('food')) return '🍱';
  if (t.includes('belanja') || t.includes('shopping')) return '🛒';
  if (t.includes('pulsa') || t.includes('telpon')) return '📱';
  return '🎁';
};

const fetchData = async () => {
  isLoading.value = true;
  try {
    const [balRes, rwRes, apRes] = await Promise.all([
      apiClient.get('/leaderboard/my-balance'),
      apiClient.get('/rewards'),
      apiClient.get('/available-points')
    ]);
    balance.value = balRes.data.balance || 0;
    rewards.value = rwRes.data.data || [];
    availablePoints.value = apRes.data.data || [];
  } catch (err) {
    toast.error('Gagal memuat data katalog hadiah.');
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

const confirmRedeem = (reward: any) => {
  selectedReward.value = reward;
  selectedSubmissionIds.value.clear();
};

const executeRedeem = async () => {
  if (!selectedReward.value) return;
  
  isRedeeming.value = true;
  
  try {
    const res = await apiClient.post(`/rewards/redeem/${selectedReward.value.id}`, {
      submission_ids: Array.from(selectedSubmissionIds.value)
    });
    toast.success(res.data.message || `Berhasil menukar ${selectedReward.value.title}! 🎉`);
    
    // Refresh catalog and points to update stock and balance
    fetchData();
    
  } catch (err: any) {
    if (err.response && err.response.data && err.response.data.error) {
      toast.error(err.response.data.error);
    } else {
      toast.error('Gagal melakukan penukaran. Coba lagi.');
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
