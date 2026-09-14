<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Top Header -->
    <div class="bg-primary pt-12 pb-24 px-6 rounded-b-[2rem] shadow-md relative">
      <div class="flex justify-between items-center relative z-10">
        <div>
          <h2 class="text-white text-sm opacity-80">Halo,</h2>
          <h1 class="text-white text-xl font-bold">{{ auth.user?.name || auth.user?.npk || 'Karyawan' }}</h1>
        </div>
        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm border border-white/30 shadow-inner">
          <span class="text-white font-bold text-lg">{{ getInitial(auth.user?.name || auth.user?.npk || 'U') }}</span>
        </div>
      </div>
    </div>

    <!-- Main Content (Overlapping Header) -->
    <div class="px-5 -mt-16 relative z-20 space-y-5">
      
      <!-- Balance Card -->
      <div class="bg-white rounded-2xl p-5 shadow-lg border border-gray-100 flex items-center justify-between">
        <div>
          <p class="text-gray-500 text-xs font-medium mb-1">Total Poin GROW</p>
          <div class="flex items-end gap-1">
            <span class="text-3xl font-extrabold text-gray-900 leading-none">{{ balance }}</span>
            <span class="text-sm font-semibold text-secondary mb-1">pts</span>
          </div>
        </div>
        <div class="h-16 w-[1px] bg-gray-100"></div>
        <div class="text-right">
          <p class="text-gray-500 text-xs font-medium mb-1">Peringkat</p>
          <div class="flex items-center justify-end gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
            </svg>
            <span class="text-2xl font-bold text-gray-900">#{{ rank }}</span>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-2 gap-4">
        <router-link to="/submission" class="bg-white rounded-xl p-4 shadow-sm border border-gray-100 flex flex-col items-center justify-center gap-2 hover:border-primary transition-all group">
          <div class="w-12 h-12 rounded-full bg-blue-50 flex items-center justify-center group-hover:bg-primary transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-primary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <span class="text-xs font-semibold text-gray-800 text-center">Buat<br>Pengajuan</span>
        </router-link>
        
        <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100 flex flex-col items-center justify-center gap-2 hover:border-secondary transition-all group cursor-not-allowed opacity-80">
          <div class="w-12 h-12 rounded-full bg-orange-50 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
          </div>
          <span class="text-xs font-semibold text-gray-800 text-center">Tukar<br>Hadiah</span>
        </div>
      </div>

      <!-- Recent Submissions Tracker -->
      <div class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-gray-800">Riwayat Pengajuan</h3>
          <button class="text-xs text-primary font-medium">Lihat Semua</button>
        </div>

        <div v-if="isLoading" class="flex justify-center py-8">
          <svg class="animate-spin h-6 w-6 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <div v-else-if="submissions.length === 0" class="text-center py-8">
          <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-sm text-gray-500">Belum ada pengajuan poin.</p>
        </div>

        <div v-else class="space-y-3">
          <div v-for="sub in submissions" :key="sub.id" class="flex items-center justify-between p-3 rounded-xl border border-gray-100 bg-gray-50/50">
            <div class="flex items-center gap-3">
              <div :class="getStatusIconBg(sub.status)" class="w-10 h-10 rounded-full flex items-center justify-center shrink-0">
                <svg v-if="sub.status === 'PENDING'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                <svg v-else-if="sub.status === 'APPROVED'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </div>
              <div class="overflow-hidden">
                <p class="text-sm font-semibold text-gray-800 truncate">{{ sub.grow_id || 'Pengajuan Baru' }}</p>
                <p class="text-xs text-gray-500 truncate">{{ formatDate(sub.created_at) }}</p>
              </div>
            </div>
            <div class="text-right shrink-0">
              <span :class="getStatusTextColor(sub.status)" class="text-xs font-bold px-2 py-1 rounded-md" :style="{ backgroundColor: getStatusBgColor(sub.status) }">
                {{ sub.status }}
              </span>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';

const auth = useAuthStore();
const balance = ref(0); 
const rank = ref('-');
const submissions = ref<any[]>([]);
const isLoading = ref(true);

const getInitial = (name: string) => {
  return name.charAt(0).toUpperCase();
};

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
};

const getStatusIconBg = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'bg-green-100';
    case 'REJECTED': return 'bg-red-100';
    default: return 'bg-yellow-100';
  }
};

const getStatusTextColor = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'text-green-700';
    case 'REJECTED': return 'text-red-700';
    default: return 'text-yellow-700';
  }
};

const getStatusBgColor = (status: string) => {
  switch (status) {
    case 'APPROVED': return '#dcfce7'; // green-100
    case 'REJECTED': return '#fee2e2'; // red-100
    default: return '#fef9c3'; // yellow-100
  }
};

const fetchDashboardData = async () => {
  isLoading.value = true;
  try {
    // Try to fetch leaderboard balance
    try {
      const balRes = await apiClient.get('/leaderboard/my-balance');
      balance.value = balRes.data.balance || 0;
      rank.value = balRes.data.rank || '-';
    } catch (err) {
      console.warn("Could not fetch balance", err);
    }

    // Fetch recent submissions
    try {
      const subRes = await apiClient.get('/submissions/me?limit=5');
      submissions.value = subRes.data || [];
    } catch (err) {
      console.warn("Could not fetch submissions", err);
    }
    
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchDashboardData();
});
</script>
