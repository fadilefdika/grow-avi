<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Top Header -->
    <div class="bg-primary pt-12 pb-24 px-6 rounded-b-[2rem] shadow-md relative">
      <div class="flex justify-between items-center relative z-10">
        <div>
          <h2 class="text-white text-sm opacity-80">Halo,</h2>
          <h1 class="text-white text-xl font-bold">{{ auth.user?.name || auth.user?.npk || 'Karyawan' }}</h1>
        </div>
        <div class="flex items-center gap-3 relative z-10">
          <div class="w-10 h-10 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm border border-white/30 shadow-inner">
            <span class="text-white font-bold text-lg">{{ getInitial(auth.user?.name || auth.user?.npk || 'U') }}</span>
          </div>
          <button @click="handleLogout" class="w-8 h-8 rounded-full bg-white/10 hover:bg-red-500/80 flex items-center justify-center transition-colors text-white border border-white/20">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
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
        <div class="text-right cursor-pointer group" @click="$router.push('/leaderboard')">
          <p class="text-gray-500 text-xs font-medium mb-1 group-hover:text-primary transition-colors">Peringkat <span class="text-[10px] bg-primary/10 text-primary px-1 rounded ml-1">Lihat</span></p>
          <div class="flex items-center justify-end gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
            </svg>
            <span class="text-2xl font-bold text-gray-900 group-hover:text-primary transition-colors">#{{ rank }}</span>
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
        
        <router-link to="/rewards" class="bg-white rounded-xl p-4 shadow-sm border border-gray-100 flex flex-col items-center justify-center gap-2 hover:border-secondary transition-all group">
          <div class="w-12 h-12 rounded-full bg-orange-50 flex items-center justify-center group-hover:bg-secondary transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-secondary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
          </div>
          <span class="text-xs font-semibold text-gray-800 text-center">Tukar<br>Hadiah</span>
        </router-link>
      </div>

      <!-- Admin Actions (Only if admin) -->
      <div v-if="auth.user?.role === 'admin'" class="bg-gray-800 rounded-xl p-4 shadow-sm flex items-center justify-between mt-2">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-gray-700 rounded-full flex items-center justify-center">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <div>
            <h3 class="text-white font-bold text-sm">Mode Administrator</h3>
            <p class="text-gray-400 text-xs">Akses panel admin penuh</p>
          </div>
        </div>
        <router-link to="/admin/dashboard" class="px-4 py-2 bg-yellow-500 text-gray-900 font-bold rounded-lg text-sm hover:bg-yellow-400 transition-colors">
          Buka Panel
        </router-link>
      </div>

      <!-- Recent Submissions Tracker -->
      <div class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-gray-800">Riwayat Pengajuan</h3>
          <button @click="$router.push('/history')" v-if="allSubmissions.length > 5" class="text-xs text-primary font-medium hover:underline">
            Lihat Semua
          </button>
        </div>

        <div v-if="isLoading" class="flex justify-center py-8">
          <svg class="animate-spin h-6 w-6 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <div v-else-if="allSubmissions.length === 0" class="text-center py-8">
          <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-sm text-gray-500">Belum ada pengajuan poin.</p>
        </div>

        <div v-else class="space-y-3">
          <div v-for="sub in displayedSubmissions" :key="sub.id" @click="openDetail(sub)" class="flex items-center justify-between p-3 rounded-xl border border-gray-100 bg-gray-50/50 cursor-pointer hover:bg-gray-100 transition-colors group">
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

    <!-- Detail Modal -->
    <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl max-h-[90vh] flex flex-col overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50 shrink-0">
          <h3 class="text-lg font-bold text-gray-900">Detail Pengajuan</h3>
          <button @click="closeDetailModal" class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="p-6 overflow-y-auto">
          <div v-if="isLoadingDetail" class="flex justify-center py-12">
            <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          
          <div v-else-if="selectedDetail" class="space-y-6">
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs text-gray-500 font-medium">GROW ID</p>
                <p class="text-lg font-bold text-gray-900">{{ selectedDetail.grow_id }}</p>
              </div>
              <span :class="getStatusTextColor(selectedDetail.status)" :style="{ backgroundColor: getStatusBgColor(selectedDetail.status) }" class="px-3 py-1 text-xs font-bold rounded-full border">
                {{ selectedDetail.status }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="bg-gray-50 p-3 rounded-lg border border-gray-100" v-if="selectedDetail.category_name !== 'INNOVATION'">
                <p class="text-xs text-gray-500 mb-1">Aktivitas</p>
                <p class="text-sm font-semibold text-gray-900">{{ (selectedDetail.activity_name === 'Yang lain (Custom)' && selectedDetail.custom_activity_type?.Valid) ? selectedDetail.custom_activity_type.String : (selectedDetail.activity_name || '-') }}</p>
              </div>
              <div class="bg-gray-50 p-3 rounded-lg border border-gray-100">
                <p class="text-xs text-gray-500 mb-1">Kategori</p>
                <p class="text-sm font-semibold text-gray-900">{{ selectedDetail.category_name || '-' }}</p>
              </div>
              <div class="bg-gray-50 p-3 rounded-lg border border-gray-100">
                <p class="text-xs text-gray-500 mb-1">Poin</p>
                <p class="text-sm font-bold text-primary">{{ selectedDetail.points_awarded || '-' }} pts</p>
              </div>
              <div class="bg-gray-50 p-3 rounded-lg border border-gray-100">
                <p class="text-xs text-gray-500 mb-1">Tanggal Aktivitas</p>
                <p class="text-sm font-semibold text-gray-900">{{ formatDate(selectedDetail.activity_date) }}</p>
              </div>
            </div>

            <div v-if="selectedDetail.custom_reference?.Valid" class="bg-gray-50 p-3 rounded-lg border border-gray-100">
              <p class="text-xs text-gray-500 mb-1">Nama Kegiatan/Aktivitas</p>
              <p class="text-sm font-medium text-gray-900">{{ selectedDetail.custom_reference.String }}</p>
            </div>

            <div v-if="selectedDetail.nomor_ss?.Valid" class="bg-gray-50 p-3 rounded-lg border border-gray-100">
              <p class="text-xs text-gray-500 mb-1">Nomor SS</p>
              <p class="text-sm font-bold text-gray-900">{{ selectedDetail.nomor_ss.String }}</p>
            </div>

            <div v-if="selectedDetail.admin_notes?.Valid" class="bg-yellow-50 p-4 rounded-xl border border-yellow-200">
              <p class="text-xs font-bold text-yellow-800 mb-1">Catatan Admin:</p>
              <p class="text-sm text-yellow-700">{{ selectedDetail.admin_notes.String }}</p>
            </div>

            <!-- Evidence Section -->
            <div v-if="selectedDetail.evidence && selectedDetail.evidence.length > 0">
              <p class="text-sm font-bold text-gray-900 mb-3">Bukti Lampiran</p>
              <div class="grid grid-cols-2 gap-3">
                <a v-for="(ev, idx) in selectedDetail.evidence" :key="idx" :href="ev" target="_blank" class="block group relative rounded-lg overflow-hidden border border-gray-200 bg-gray-50">
                  <div class="aspect-video w-full">
                    <img :src="ev" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                  </div>
                  <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                  </div>
                </a>
              </div>
            </div>
          </div>
        </div>
        
        <div class="p-4 border-t border-gray-100 bg-gray-50 shrink-0 flex gap-3">
          <button v-if="selectedDetail.status === 'REJECTED'" @click="goToRevision" class="flex-1 py-2.5 px-4 bg-primary text-white font-medium rounded-lg hover:bg-blue-700 transition-colors">
            Revisi Pengajuan
          </button>
          <button @click="closeDetailModal" class="flex-1 py-2.5 px-4 bg-white border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors">
            Tutup
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';

const authStore = useAuthStore();
const router = useRouter();

const auth = useAuthStore();
const balance = ref(0); 
const rank = ref('-');
const allSubmissions = ref<any[]>([]);
const showAll = ref(false); // No longer used but kept to avoid undefined errors if referenced elsewhere
const isLoading = ref(true);

const displayedSubmissions = computed(() => {
  return allSubmissions.value.slice(0, 5);
});

const showDetailModal = ref(false);
const selectedDetail = ref<any>(null);
const isLoadingDetail = ref(false);

const openDetail = async (sub: any) => {
  try {
    isLoadingDetail.value = true;
    showDetailModal.value = true;
    selectedDetail.value = null;
    const res = await apiClient.get(`/submissions/${sub.id}`);
    selectedDetail.value = res.data.data;
  } catch (err) {
    console.error("Gagal mengambil detail", err);
    showDetailModal.value = false;
  } finally {
    isLoadingDetail.value = false;
  }
};

const closeDetailModal = () => {
  showDetailModal.value = false;
  selectedDetail.value = null;
};

const goToRevision = () => {
  if (selectedDetail.value) {
    router.push(`/submission?revisi_id=${selectedDetail.value.id}`);
  }
};

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
      const subRes = await apiClient.get('/submissions/me');
      allSubmissions.value = subRes.data.submissions || [];
    } catch (err) {
      console.warn("Could not fetch submissions", err);
    }
    
  } finally {
    isLoading.value = false;
  }
};

const handleLogout = async () => {
  try {
    await apiClient.post('/auth/logout');
  } catch (err) {
    console.error('Logout error:', err);
  } finally {
    auth.logout();
    router.push('/login');
  }
};

onMounted(() => {
  fetchDashboardData();
});
</script>
