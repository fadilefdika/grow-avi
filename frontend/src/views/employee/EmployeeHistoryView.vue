<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Hero Header -->
    <div class="hero-gradient pt-12 pb-28 px-6 relative overflow-hidden">
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-8 -right-8 w-40 h-40 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #1CB0F6, transparent);"></div>
        <div class="absolute top-20 -left-6 w-28 h-28 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFA64D, transparent); animation-delay: 1s;"></div>
      </div>
      <div class="flex items-start gap-4 relative z-10">
        <button @click="$router.push('/dashboard')" class="w-10 h-10 rounded-xl glass-card flex items-center justify-center text-white hover:bg-white/25 transition-all mt-0.5 flex-shrink-0">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-white text-2xl font-black">📋 Riwayat Pengajuan</h1>
          <p class="text-white/60 text-xs font-semibold mt-0.5">Pantau semua perkembangan pengajuanmu</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="px-4 -mt-20 relative z-20 space-y-4">

      <!-- Filter Card -->
      <div class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100">
        <div class="flex gap-2">
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari GROW ID..."
              class="w-full pl-9 pr-3 py-2.5 bg-gray-50 border-2 border-gray-200 rounded-xl text-sm font-semibold focus:outline-none focus:border-primary focus:bg-white transition-colors placeholder-gray-300"
              @input="onSearch"
            />
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400 absolute left-3 top-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <select
            v-model="selectedCategory"
            @change="fetchHistory(1)"
            class="w-32 px-2 py-2.5 bg-gray-50 border-2 border-gray-200 rounded-xl text-sm font-semibold focus:outline-none focus:border-primary focus:bg-white transition-colors"
          >
            <option value="">Semua</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
      </div>

      <!-- List Content -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div v-if="isLoading" class="flex justify-center py-12">
          <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        
        <div v-else-if="submissions.length === 0" class="text-center py-12 px-6">
          <div class="w-16 h-16 bg-gray-50 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-sm font-black text-gray-400">Tidak ada data riwayat ditemukan.</p>
          <p class="text-xs text-gray-300 font-semibold mt-1">Coba ubah filter pencarianmu.</p>
        </div>
        
        <div v-else class="divide-y divide-gray-100">
          <div
            v-for="sub in submissions"
            :key="sub.id"
            @click="openDetail(sub)"
            class="p-4 cursor-pointer hover:bg-blue-50/50 transition-colors group flex flex-col gap-3"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3 overflow-hidden">
                <div :class="getStatusIconBg(sub.status)" class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-sm">
                  <svg v-if="sub.status === 'PENDING'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                  <svg v-else-if="sub.status === 'APPROVED'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" /></svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" /></svg>
                </div>
                <div class="overflow-hidden">
                  <p class="text-xs font-bold text-gray-400 truncate">{{ sub.grow_id || 'Pengajuan Baru' }}</p>
                  <p class="text-sm font-black text-gray-800 truncate group-hover:text-primary transition-colors mt-0.5">{{ sub.category_name }}</p>
                </div>
              </div>
              <div class="flex flex-col items-end shrink-0 pl-2">
                <span :class="getStatusBadgeClass(sub.status)" class="mb-1 text-[10px]">
                  {{ sub.status }}
                </span>
                <p class="text-[10px] font-bold text-gray-400">{{ formatDate(sub.created_at) }}</p>
              </div>
            </div>
            
            <!-- Alasan Reject -->
            <div v-if="sub.status === 'REJECTED' && sub.admin_notes && (sub.admin_notes.Valid ? sub.admin_notes.String : sub.admin_notes)" class="bg-red-50 border border-red-100 rounded-xl p-3 flex gap-2 items-start animate-fade-in">
              <span class="text-sm">❌</span>
              <div>
                <p class="text-[10px] font-black text-red-500 uppercase tracking-wide">Alasan Penolakan</p>
                <p class="text-xs font-semibold text-red-800 leading-tight mt-0.5">{{ sub.admin_notes.Valid ? sub.admin_notes.String : sub.admin_notes }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-100 p-4">
          <button
            @click="fetchHistory(currentPage - 1)"
            :disabled="currentPage === 1"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-bold transition-colors',
              currentPage === 1
                ? 'text-gray-300 cursor-not-allowed bg-gray-50'
                : 'text-primary bg-primary/10 hover:bg-primary/20'
            ]"
          >
            ← Sebelumnya
          </button>
          
          <span class="text-sm font-bold text-gray-500">
            {{ currentPage }} / {{ totalPages }}
          </span>
          
          <button
            @click="fetchHistory(currentPage + 1)"
            :disabled="currentPage === totalPages"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-bold transition-colors',
              currentPage === totalPages
                ? 'text-gray-300 cursor-not-allowed bg-gray-50'
                : 'text-primary bg-primary/10 hover:bg-primary/20'
            ]"
          >
            Selanjutnya →
          </button>
        </div>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-t-3xl sm:rounded-3xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col overflow-hidden animate-slide-up">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center shrink-0">
          <h3 class="text-lg font-black text-gray-900">Detail Pengajuan</h3>
          <button @click="closeDetailModal" class="w-8 h-8 rounded-xl bg-gray-100 hover:bg-gray-200 flex items-center justify-center transition-colors text-gray-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="p-6 overflow-y-auto">
          <div v-if="isLoadingDetail" class="flex justify-center py-12">
            <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <div v-else-if="selectedDetail" class="space-y-5">
            <!-- Status & ID -->
            <div class="bg-gray-50 p-4 rounded-2xl border border-gray-100 flex justify-between items-center">
              <div>
                <p class="text-xs text-gray-400 font-bold uppercase tracking-wide mb-0.5">GROW ID</p>
                <p class="font-black text-gray-900 text-lg">{{ selectedDetail.grow_id }}</p>
                <p class="text-xs font-semibold text-gray-400 mt-0.5">{{ selectedDetail.category_name }}</p>
              </div>
              <span :class="getStatusBadgeClass(selectedDetail.status)" class="text-sm">
                {{ selectedDetail.status }}
              </span>
            </div>

            <!-- Rejection notes -->
            <div v-if="selectedDetail.status === 'REJECTED' && selectedDetail.admin_notes?.Valid" class="bg-red-50 p-4 rounded-2xl border border-red-100">
              <p class="text-xs font-black text-red-500 uppercase tracking-wide mb-1.5">❌ Alasan Penolakan</p>
              <p class="text-sm text-red-800 font-semibold">{{ selectedDetail.admin_notes.String }}</p>
            </div>

            <!-- Approved points -->
            <div v-if="selectedDetail.status === 'APPROVED'" class="bg-green-50 p-4 rounded-2xl border border-green-100 flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-green-100 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" /></svg>
              </div>
              <div>
                <p class="text-xs font-bold text-green-600 uppercase tracking-wide">Poin Didapat</p>
                <p class="text-xl font-black text-green-800">{{ selectedDetail.points_awarded }} <span class="text-sm font-bold">pts</span></p>
              </div>
            </div>

            <!-- Activity Details -->
            <div>
              <p class="text-xs font-black text-gray-400 uppercase tracking-wide mb-2">Detail Aktivitas</p>
              <div class="bg-white border-2 border-gray-100 rounded-2xl overflow-hidden divide-y divide-gray-100">
                <div class="p-3.5 flex justify-between gap-4">
                  <span class="text-sm font-semibold text-gray-500">Nama Aktivitas</span>
                  <span class="text-sm font-bold text-gray-900 text-right">{{ selectedDetail.activity_name === 'Yang lain (Custom)' && selectedDetail.custom_activity_type?.Valid ? selectedDetail.custom_activity_type.String : selectedDetail.activity_name }}</span>
                </div>
                <div v-if="selectedDetail.custom_reference?.Valid" class="p-3.5 flex justify-between gap-4">
                  <span class="text-sm font-semibold text-gray-500">Kegiatan Spesifik</span>
                  <span class="text-sm font-bold text-gray-900 text-right max-w-[200px]">{{ selectedDetail.custom_reference.String }}</span>
                </div>
                <div v-if="selectedDetail.nomor_ss?.Valid" class="p-3.5 flex justify-between gap-4">
                  <span class="text-sm font-semibold text-gray-500">Nomor SS</span>
                  <span class="text-sm font-black text-gray-900 text-right">{{ selectedDetail.nomor_ss.String }}</span>
                </div>
                <div class="p-3.5 flex justify-between gap-4">
                  <span class="text-sm font-semibold text-gray-500">Tanggal Transaksi</span>
                  <span class="text-sm font-bold text-gray-900">{{ formatDateOnly(selectedDetail.activity_date) }}</span>
                </div>
              </div>
            </div>

            <!-- Evidence -->
            <div v-if="selectedDetail.evidence && selectedDetail.evidence.length > 0">
              <p class="text-xs font-black text-gray-400 uppercase tracking-wide mb-2">Bukti Lampiran</p>
              <div class="grid grid-cols-2 gap-3">
                <a v-for="(ev, idx) in selectedDetail.evidence" :key="idx" :href="ev" target="_blank" class="block group relative rounded-2xl overflow-hidden border border-gray-200 bg-gray-50">
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
        
        <div class="p-4 border-t border-gray-100 bg-gray-50/50 shrink-0 flex gap-3">
          <button v-if="selectedDetail?.status === 'REJECTED'" @click="goToRevision" class="btn-game-primary flex-1 py-3 text-sm">
            Revisi Pengajuan
          </button>
          <button @click="closeDetailModal" class="flex-1 py-3 bg-white border-2 border-gray-200 text-gray-700 font-bold rounded-2xl hover:bg-gray-50 transition-colors text-sm">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import apiClient from '../../api/client';

const router = useRouter();
const submissions = ref<any[]>([]);
const categories = ref<any[]>([]);
const isLoading = ref(true);

const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const searchQuery = ref('');
const selectedCategory = ref('');

const showDetailModal = ref(false);
const selectedDetail = ref<any>(null);
const isLoadingDetail = ref(false);

let searchTimeout: ReturnType<typeof setTimeout>;

const onSearch = () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    fetchHistory(1);
  }, 500);
};

const fetchCategories = async () => {
  try {
    const res = await apiClient.get('/categories');
    categories.value = res.data.data.filter((c: any) => c.is_active);
  } catch (err) {
    console.error('Failed to fetch categories', err);
  }
};

const fetchHistory = async (page: number) => {
  try {
    isLoading.value = true;
    currentPage.value = page;
    
    let url = `/submissions/me?page=${page}&limit=10`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (selectedCategory.value) url += `&category_id=${selectedCategory.value}`;
    
    const res = await apiClient.get(url);
    submissions.value = res.data.submissions;
    if (res.data.meta) {
      totalPages.value = res.data.meta.total_pages;
      totalItems.value = res.data.meta.total;
    }
  } catch (err) {
    console.error('Failed to fetch history', err);
  } finally {
    isLoading.value = false;
  }
};

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

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Jakarta' }).replace(/\./g, ':');
};

const formatDateOnly = (dateString: string) => {
  if (!dateString) return '-';
  const parts = dateString.split('T')[0].split('-');
  if (parts.length === 3) {
    const d = new Date(parseInt(parts[0]), parseInt(parts[1])-1, parseInt(parts[2]));
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', timeZone: 'Asia/Jakarta' });
  }
  return dateString;
};

const getStatusIconBg = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'bg-green-100';
    case 'REJECTED': return 'bg-red-100';
    default: return 'bg-amber-100';
  }
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'badge-approved';
    case 'REJECTED': return 'badge-rejected';
    default: return 'badge-pending';
  }
};

onMounted(() => {
  fetchCategories();
  fetchHistory(1);
});
</script>
