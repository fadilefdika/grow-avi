<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Top Header -->
    <div class="bg-primary pt-12 pb-16 px-6 rounded-b-[2rem] shadow-md relative">
      <div class="flex items-center gap-4 relative z-10">
        <button @click="$router.push('/dashboard')" class="w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition-colors text-white border border-white/20">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-white text-xl font-bold">Riwayat Pengajuan</h1>
          <p class="text-white/80 text-xs mt-0.5">Semua data riwayat pengajuan poin Anda</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="px-4 -mt-8 relative z-20 space-y-4">
      
      <!-- Filter Card -->
      <div class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 flex flex-col gap-3">
        <div class="flex gap-2">
          <div class="flex-1 relative">
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Cari GROW ID..." 
              class="w-full pl-9 pr-3 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors"
              @input="onSearch"
            />
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400 absolute left-3 top-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <select v-model="selectedCategory" @change="fetchHistory(1)" class="w-32 px-2 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors">
            <option value="">Semua Kategori</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
      </div>

      <!-- List Content -->
      <div class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 min-h-[400px]">
        <div v-if="isLoading" class="flex justify-center py-12">
          <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        
        <div v-else-if="submissions.length === 0" class="text-center py-12">
          <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-sm text-gray-500 font-medium">Tidak ada data riwayat ditemukan.</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="sub in submissions" :key="sub.id" @click="openDetail(sub)" class="flex items-center justify-between p-3 rounded-xl border border-gray-100 bg-gray-50/50 cursor-pointer hover:bg-blue-50 transition-colors group">
            <div class="flex items-center gap-3">
              <div :class="getStatusIconBg(sub.status)" class="w-10 h-10 rounded-full flex items-center justify-center shrink-0">
                <svg v-if="sub.status === 'PENDING'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                <svg v-else-if="sub.status === 'APPROVED'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </div>
              <div class="overflow-hidden">
                <p class="text-sm font-semibold text-gray-800 truncate">{{ sub.grow_id || 'Pengajuan Baru' }}</p>
                <p class="text-[11px] text-gray-500 truncate">{{ sub.category_name }} • {{ formatDate(sub.created_at) }}</p>
              </div>
            </div>
            <div class="text-right shrink-0">
              <span :class="getStatusTextColor(sub.status)" class="text-xs font-bold px-2 py-1 rounded-md" :style="{ backgroundColor: getStatusBgColor(sub.status) }">
                {{ sub.status }}
              </span>
            </div>
          </div>
        </div>

        <!-- Pagination Controls -->
        <div v-if="totalPages > 1" class="mt-6 flex items-center justify-between border-t border-gray-100 pt-4">
          <button 
            @click="fetchHistory(currentPage - 1)" 
            :disabled="currentPage === 1"
            class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors"
            :class="currentPage === 1 ? 'text-gray-300 cursor-not-allowed bg-gray-50' : 'text-gray-700 bg-white border border-gray-200 hover:bg-gray-50'"
          >
            Sebelumnya
          </button>
          
          <span class="text-sm font-medium text-gray-600">
            Hal {{ currentPage }} dari {{ totalPages }}
          </span>
          
          <button 
            @click="fetchHistory(currentPage + 1)" 
            :disabled="currentPage === totalPages"
            class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors"
            :class="currentPage === totalPages ? 'text-gray-300 cursor-not-allowed bg-gray-50' : 'text-gray-700 bg-white border border-gray-200 hover:bg-gray-50'"
          >
            Selanjutnya
          </button>
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
            <div>
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Status & ID</h4>
              <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 flex justify-between items-center">
                <div>
                  <p class="font-bold text-gray-900">{{ selectedDetail.grow_id }}</p>
                  <p class="text-xs text-gray-500 mt-0.5">Kategori: {{ selectedDetail.category_name }}</p>
                </div>
                <span :class="getStatusTextColor(selectedDetail.status)" class="text-sm font-bold px-3 py-1.5 rounded-lg" :style="{ backgroundColor: getStatusBgColor(selectedDetail.status) }">
                  {{ selectedDetail.status }}
                </span>
              </div>
            </div>

            <div v-if="selectedDetail.status === 'REJECTED' && selectedDetail.admin_notes?.Valid">
              <h4 class="text-xs font-bold text-red-400 uppercase tracking-wider mb-2">Catatan Penolakan / Revisi</h4>
              <div class="bg-red-50 p-4 rounded-xl border border-red-100 text-red-800 text-sm font-medium">
                {{ selectedDetail.admin_notes.String }}
              </div>
            </div>

            <div v-if="selectedDetail.status === 'APPROVED'">
              <h4 class="text-xs font-bold text-green-500 uppercase tracking-wider mb-2">Penilaian Poin</h4>
              <div class="bg-green-50 p-4 rounded-xl border border-green-100">
                <p class="text-sm text-green-900">Poin Diberikan: <span class="font-bold text-xl">{{ selectedDetail.points_awarded }}</span> pts</p>
              </div>
            </div>

            <div>
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Detail Aktivitas</h4>
              <div class="bg-white border border-gray-200 rounded-xl overflow-hidden divide-y divide-gray-100">
                <div class="p-3 flex justify-between">
                  <span class="text-sm text-gray-500">Nama Aktivitas</span>
                  <span class="text-sm font-medium text-gray-900 text-right">{{ selectedDetail.activity_name === 'Yang lain (Custom)' && selectedDetail.custom_activity_type?.Valid ? selectedDetail.custom_activity_type.String : selectedDetail.activity_name }}</span>
                </div>
                <div v-if="selectedDetail.custom_reference?.Valid" class="p-3 flex justify-between">
                  <span class="text-sm text-gray-500">Kegiatan Spesifik</span>
                  <span class="text-sm font-medium text-gray-900 text-right max-w-[200px]">{{ selectedDetail.custom_reference.String }}</span>
                </div>
                <div v-if="selectedDetail.nomor_ss?.Valid" class="p-3 flex justify-between">
                  <span class="text-sm text-gray-500">Nomor SS</span>
                  <span class="text-sm font-bold text-gray-900 text-right">{{ selectedDetail.nomor_ss.String }}</span>
                </div>
                <div class="p-3 flex justify-between">
                  <span class="text-sm text-gray-500">Tanggal Transaksi</span>
                  <span class="text-sm font-medium text-gray-900">{{ formatDateOnly(selectedDetail.activity_date) }}</span>
                </div>
              </div>
            </div>

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
          <button v-if="selectedDetail?.status === 'REJECTED'" @click="goToRevision" class="flex-1 py-2.5 px-4 bg-primary text-white font-medium rounded-lg hover:bg-blue-700 transition-colors">
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
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
};

const formatDateOnly = (dateString: string) => {
  if (!dateString) return '-';
  const parts = dateString.split('T')[0].split('-');
  if (parts.length === 3) {
    const d = new Date(parseInt(parts[0]), parseInt(parts[1])-1, parseInt(parts[2]));
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
  }
  return dateString;
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

onMounted(() => {
  fetchCategories();
  fetchHistory(1);
});
</script>
