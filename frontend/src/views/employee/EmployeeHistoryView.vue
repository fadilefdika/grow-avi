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

      <!-- Filter Panel -->
      <div class="bg-white p-5 rounded-2xl border border-gray-100 shadow-sm mb-6 flex flex-wrap gap-4 items-end">
        <div class="flex-1 min-w-[200px]">
          <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Cari GROW ID</label>
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Masukkan GROW ID..."
              class="w-full pl-9 pr-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-all"
              @input="onSearch"
            />
            <svg class="h-4 w-4 text-gray-400 absolute left-3 top-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          </div>
        </div>

        <div class="w-full sm:w-48 relative">
          <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Kategori</label>
          <select
            v-model="selectedCategory"
            @change="onFilterChange"
            class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-all appearance-none"
          >
            <option value="">Semua Kategori</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
          <svg class="w-4 h-4 text-gray-400 absolute right-3 top-[28px] pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
        </div>
        
        <div>
          <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Dari Tanggal</label>
          <input type="date" v-model="filterDateFrom" class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="onFilterChange" />
        </div>
        
        <div>
          <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Sampai Tanggal</label>
          <input type="date" v-model="filterDateTo" class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="onFilterChange" />
        </div>

        <div class="flex gap-2">
          <button @click="resetFilter" class="px-4 py-1.5 text-sm text-gray-600 hover:text-gray-900 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors font-medium">Reset</button>
        </div>
      </div>

      <!-- Table Content -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <DataTable
          title="Riwayat Pengajuan"
          :columns="[
            { key: 'grow_id', label: 'GROW ID', sortable: false },
            { key: 'category_name', label: 'Kategori', sortable: false },
            { key: 'status', label: 'Status', sortable: false },
            { key: 'activity_date', label: 'Tanggal Kegiatan', sortable: false },
            { key: 'created_at', label: 'Waktu Submit', sortable: false }
          ]"
          :rows="submissions"
          :totalRows="totalItems"
          :serverSide="true"
          @change="onTableChange"
          @rowClick="openDetail"
        >
          <template #cell-status="{ row }">
            <span :class="getStatusBadgeClass(row.status)" class="inline-flex items-center px-2.5 py-1 rounded-full text-[10px] font-black whitespace-nowrap">
              {{ formatStatus(row.status) }}
            </span>
          </template>
          
          <template #cell-activity_date="{ row }">
            {{ formatDateOnly(row.activity_date) }}
          </template>

          <template #cell-created_at="{ row }">
            {{ formatDate(row.created_at) }} WIB
          </template>

          <template #actions="{ row }">
            <button @click.stop="openDetail(row)" class="text-primary hover:text-blue-900 font-semibold bg-blue-50 hover:bg-blue-100 px-3 py-1.5 rounded-lg transition-colors text-xs whitespace-nowrap">Lihat Detail</button>
          </template>
        </DataTable>
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
            <div v-if="(selectedDetail.status === 'REJECTED' || selectedDetail.status === 'REJECTED_RESUBMIT') && selectedDetail.admin_notes?.Valid" class="bg-red-50 p-4 rounded-2xl border border-red-100">
              <p class="text-xs font-black text-red-500 uppercase tracking-wide mb-1.5">❌ Alasan Penolakan</p>
              <p class="text-sm text-red-800 font-semibold mb-3">{{ selectedDetail.admin_notes.String }}</p>
              <button v-if="selectedDetail.status === 'REJECTED_RESUBMIT'" @click.stop="$router.push({ path: '/submit', query: { revisi_id: selectedDetail.id } })" class="w-full px-4 py-2 bg-orange-500 text-white rounded-xl text-sm font-bold hover:bg-orange-600 transition-colors">
                Perbaiki Pengajuan Sekarang
              </button>
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
                  <span class="text-sm font-bold text-gray-900 text-right">{{ selectedDetail.activity_name === 'yang lain' && selectedDetail.custom_activity_type?.Valid ? selectedDetail.custom_activity_type.String : selectedDetail.activity_name }}</span>
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
import DataTable from '../../components/ui/DataTable.vue';

const router = useRouter();
const submissions = ref<any[]>([]);
const categories = ref<any[]>([]);
const isLoading = ref(true);

const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const searchQuery = ref('');
const selectedCategory = ref('');
const filterDateFrom = ref('');
const filterDateTo = ref('');

const showDetailModal = ref(false);
const selectedDetail = ref<any>(null);
const isLoadingDetail = ref(false);

let searchTimeout: ReturnType<typeof setTimeout>;

const onSearch = () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    tableParams.page = 1;
    fetchHistory(1);
  }, 500);
};

let tableParams = {
  page: 1,
  limit: 10
};

const onTableChange = (params: any) => {
  tableParams.page = params.page;
  tableParams.limit = params.limit;
  fetchHistory(params.page);
};

const onFilterChange = () => {
  tableParams.page = 1;
  fetchHistory(1);
};

const resetFilter = () => {
  searchQuery.value = '';
  selectedCategory.value = '';
  filterDateFrom.value = '';
  filterDateTo.value = '';
  tableParams.page = 1;
  fetchHistory(1);
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
    
    let url = `/submissions/me?page=${page}&limit=${tableParams.limit}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (selectedCategory.value) url += `&category_id=${selectedCategory.value}`;
    if (filterDateFrom.value) url += `&date_from=${filterDateFrom.value}`;
    if (filterDateTo.value) url += `&date_to=${filterDateTo.value}`;
    
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
    case 'REJECTED': 
    case 'REJECTED_RESUBMIT': return 'badge-rejected';
    default: return 'badge-pending';
  }
};

const formatStatus = (status: string) => {
  if (status === 'REJECTED_RESUBMIT') return 'REJECTED (REVISI)';
  return status;
};

onMounted(() => {
  fetchCategories();
  fetchHistory(1);
});
</script>
