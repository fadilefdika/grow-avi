<template>
  <div class="flex h-screen bg-gray-50 font-sans">
    <AdminSidebar />
    
    <div class="flex-1 flex flex-col h-screen overflow-hidden">
      <!-- Top header -->
      <header class="bg-white border-b border-gray-200 px-8 py-5 flex items-center justify-between shrink-0">
        <div>
          <h1 class="text-2xl font-bold text-gray-800 tracking-tight">Activity Log</h1>
          <p class="text-sm text-gray-500 mt-1">Riwayat semua kegiatan approval, penolakan, dan klaim hadiah.</p>
        </div>
      </header>
      
      <!-- Main Content -->
      <main class="flex-1 overflow-y-auto p-5">
        
        <!-- Filter Panel -->
        <div class="bg-white p-5 rounded-2xl border border-gray-200 shadow-sm mb-6 flex flex-wrap gap-4 items-end">
          <div class="flex-1 min-w-[200px]">
            <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Cari GROW ID / NPK / Nama</label>
            <input type="text" v-model="filterParams.search" placeholder="Cari..." class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @input="onSearchInput" />
          </div>
          
          <div class="relative min-w-[200px]">
            <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Jenis Kegiatan</label>
            <button @click.stop="toggleTypeDropdown" class="w-full text-left bg-gray-50 border border-gray-200 text-gray-700 rounded-lg px-3 py-1.5 text-sm flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showTypeDropdown}">
              <span class="truncate pr-2">{{ getTypeName(filterParams.type) }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showTypeDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showTypeDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto p-1">
              <div @click="selectType('')" class="px-3 py-2 text-sm font-medium rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="!filterParams.type ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Jenis</div>
              <div @click="selectType('APPROVED')" class="px-3 py-2 text-sm font-medium rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="filterParams.type === 'APPROVED' ? 'text-primary bg-blue-50' : 'text-gray-600'">Approval Poin</div>
              <div @click="selectType('REJECTED')" class="px-3 py-2 text-sm font-medium rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="filterParams.type === 'REJECTED' ? 'text-primary bg-blue-50' : 'text-gray-600'">Penolakan</div>
              <div @click="selectType('REDEEM')" class="px-3 py-2 text-sm font-medium rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="filterParams.type === 'REDEEM' ? 'text-primary bg-blue-50' : 'text-gray-600'">Klaim Poin (Redeem)</div>
            </div>
          </div>
          
          <div>
            <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Dari Tanggal</label>
            <input type="date" v-model="filterParams.date_from" class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="applyFilter" />
          </div>
          
          <div>
            <label class="block text-[11px] font-semibold text-gray-500 uppercase tracking-wider mb-1.5">Sampai Tanggal</label>
            <input type="date" v-model="filterParams.date_to" class="w-full px-3 py-1.5 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="applyFilter" />
          </div>
          
          <div>
            <button @click="resetFilter" class="px-4 py-1.5 text-sm text-gray-600 hover:text-gray-900 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors font-medium">Reset</button>
          </div>
        </div>

        <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
          <DataTable
            title="Riwayat Kegiatan"
            :columns="[
              { key: 'grow_id', label: 'GROW ID', sortable: true },
              { key: 'user_name', label: 'Nama', sortable: true },
              { key: 'activity_type', label: 'Jenis', sortable: false },
              { key: 'points', label: 'Poin', sortable: true },
              { key: 'activity_date', label: 'Waktu Submit', sortable: true }
            ]"
            :rows="logData"
            :totalRows="totalLogs"
            :serverSide="true"
            @change="onTableChange"
            @rowClick="openDetail"
          >
            <template #cell-activity_type="{ row }">
              <span v-if="row.activity_type === 'APPROVED'" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 border border-green-200">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
                Approved
              </span>
              <span v-else-if="row.activity_type === 'REJECTED'" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800 border border-red-200">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                Rejected
              </span>
              <span v-else-if="row.activity_type === 'REDEEM'" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800 border border-purple-200">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20"><path d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z"></path><path d="M5 5a2 2 0 00-2 2v8a2 2 0 002 2h8a2 2 0 002-2v-3a1 1 0 10-2 0v3H5V7h3a1 1 0 000-2H5z"></path></svg>
                Claim Point
              </span>
              <span v-else class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                {{ row.activity_type }}
              </span>
            </template>
            
            <template #cell-points="{ row }">
              <span :class="row.activity_type === 'REDEEM' ? 'text-red-600 font-bold' : (row.activity_type === 'APPROVED' ? 'text-green-600 font-bold' : 'text-gray-500 font-bold')">
                {{ row.activity_type === 'REDEEM' ? '-' : (row.activity_type === 'APPROVED' ? '+' : '') }}{{ row.points }}
              </span>
            </template>

            <template #cell-activity_date="{ row }">
              {{ formatDate(row.activity_date) }} WIB
            </template>
            
            <template #actions="{ row }">
              <button @click.stop="openDetail(row)" class="text-primary hover:text-blue-900 font-semibold bg-blue-50 hover:bg-blue-100 px-3 py-1.5 rounded-lg transition-colors text-xs">Lihat Detail</button>
            </template>
          </DataTable>
        </div>
      </main>
    </div>

    <!-- Review Modal -->
    <div v-if="selectedDetail" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm overflow-y-auto">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-4xl overflow-hidden flex flex-col max-h-[90vh]">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <div>
            <h3 class="text-lg font-bold text-gray-900">Detail Riwayat: {{ selectedDetail.grow_id }}</h3>
          </div>
          <button @click="closeDetail" class="text-gray-400 hover:text-gray-600">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        
        <div class="p-6 flex-1 overflow-y-auto">
          <div class="grid grid-cols-1 gap-8" :class="selectedDetail.type === 'submission' ? 'md:grid-cols-2' : ''">
            <div class="space-y-6">
              <div>
                <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Informasi Karyawan</h4>
                <div class="bg-gray-50 p-4 rounded-xl border border-gray-100">
                  <p class="font-bold text-gray-900">{{ selectedDetail.user_name }}</p>
                  <p class="text-sm text-gray-600">{{ selectedDetail.npk }} <span v-if="selectedDetail.department">— {{ selectedDetail.department }}</span></p>
                </div>
              </div>
              
              <!-- Khusus Submission -->
              <template v-if="selectedDetail.type === 'submission'">
                <div>
                  <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Detail Aktivitas</h4>
                  <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 space-y-3">
                    <div>
                      <p class="text-xs text-gray-500">Kategori & Aktivitas</p>
                      <p class="font-medium text-gray-900">
                        {{ selectedDetail.category_name }}<span v-if="selectedDetail.category_name !== 'INNOVATION'"> / {{ selectedDetail.activity_name === 'Yang lain (Custom)' && selectedDetail.custom_activity_type?.Valid ? selectedDetail.custom_activity_type.String : selectedDetail.activity_name }}</span>
                      </p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500">Tanggal Pelaksanaan</p>
                      <p class="font-medium text-gray-900">{{ formatDate(selectedDetail.activity_date) }} WIB</p>
                    </div>
                    <div v-if="selectedDetail.custom_reference?.Valid">
                      <p class="text-xs text-gray-500">Nama Kegiatan/Aktivitas</p>
                      <p class="font-medium text-gray-900">{{ selectedDetail.custom_reference.String }}</p>
                    </div>
                    <div v-if="selectedDetail.nomor_ss?.Valid">
                      <p class="text-xs text-gray-500">Nomor SS</p>
                      <p class="font-bold text-gray-900">{{ selectedDetail.nomor_ss.String }}</p>
                    </div>
                    <div v-if="selectedDetail.status === 'REJECTED'">
                      <p class="text-xs text-red-500 font-bold">Alasan Penolakan</p>
                      <p class="font-medium text-red-700 bg-red-50 p-2 rounded mt-1">{{ selectedDetail.admin_notes || '-' }}</p>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Penilaian Poin</h4>
                  <div class="bg-gray-50 p-4 rounded-xl border border-gray-100">
                    <p class="text-sm text-gray-600">Poin diberikan: <span class="font-bold text-primary text-lg">{{ selectedDetail.points_awarded || 0 }} pts</span></p>
                    <p class="text-xs text-gray-500 mt-1">Poin default aktivitas: {{ selectedDetail.default_points }} pts</p>
                  </div>
                </div>
              </template>

              <!-- Khusus Redeem -->
              <template v-if="selectedDetail.type === 'redemption'">
                <div>
                  <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Detail Penukaran</h4>
                  <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 space-y-3">
                    <div>
                      <p class="text-xs text-gray-500">Hadiah Ditukar</p>
                      <p class="font-medium text-purple-700">{{ selectedDetail.reward_title }}</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500">Waktu Penukaran</p>
                      <p class="font-medium text-gray-900">{{ formatDate(selectedDetail.created_at) }} WIB</p>
                    </div>
                    <div>
                      <p class="text-xs text-gray-500">Status</p>
                      <p class="font-bold uppercase" :class="selectedDetail.status === 'PROCESSED' ? 'text-green-600' : 'text-gray-600'">{{ selectedDetail.status }}</p>
                    </div>
                  </div>
                </div>
                <div>
                  <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Pengeluaran Poin</h4>
                  <div class="bg-red-50 p-4 rounded-xl border border-red-100">
                    <p class="text-sm text-red-900">Poin digunakan: <span class="font-bold text-red-600 text-lg">-{{ selectedDetail.points_spent }} pts</span></p>
                  </div>
                </div>
              </template>
            </div>

            <div v-if="selectedDetail.type === 'submission'">
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Bukti Lampiran ({{ selectedDetail.evidence?.length || 0 }})</h4>
              <div v-if="!selectedDetail.evidence || selectedDetail.evidence.length === 0" class="bg-gray-50 p-8 rounded-xl border border-gray-100 text-center text-gray-500 italic">
                Tidak ada lampiran.
              </div>
              <div v-else class="space-y-4">
                <div v-for="(ev, idx) in selectedDetail.evidence" :key="idx" class="border border-gray-200 rounded-xl overflow-hidden bg-gray-50 relative group">
                  <a :href="ev" target="_blank" class="block relative">
                    <img :src="ev" class="w-full h-auto max-h-64 object-contain mx-auto" />
                    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <span class="text-white bg-black/70 px-3 py-1 rounded-full text-sm font-medium flex items-center">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 21h7a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v11m0 5l4.879-4.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242z"></path></svg>
                        Buka Penuh
                      </span>
                    </div>
                  </a>
                </div>
              </div>
            </div>
            
          </div>
        </div>

        <div class="px-6 py-4 border-t border-gray-100 bg-gray-50 flex justify-end">
          <button @click="closeDetail" class="px-6 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors font-medium text-sm">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';
import DataTable from '../../components/ui/DataTable.vue';

const route = useRoute();
const logData = ref<any[]>([]);
const totalLogs = ref(0);
const isLoading = ref(true);

const selectedDetail = ref<any>(null);

const filterParams = ref({
  search: '',
  type: '',
  date_from: '',
  date_to: ''
});

let tableParams = {
  page: 1,
  limit: 10,
  sort: 'activity_date',
  order: 'desc'
};

const showTypeDropdown = ref(false);

const toggleTypeDropdown = () => {
  showTypeDropdown.value = !showTypeDropdown.value;
};

const selectType = (type: string) => {
  filterParams.value.type = type;
  showTypeDropdown.value = false;
  applyFilter();
};

const getTypeName = (type: string) => {
  if (type === 'APPROVED') return 'Approval Poin';
  if (type === 'REJECTED') return 'Penolakan';
  if (type === 'REDEEM') return 'Klaim Poin (Redeem)';
  return 'Semua Jenis';
};

if (typeof window !== 'undefined') {
  window.addEventListener('click', (e) => {
    const target = e.target as HTMLElement;
    if (!target.closest('.relative')) {
      showTypeDropdown.value = false;
    }
  });
}

const fetchLogs = async () => {
  isLoading.value = true;
  try {
    const res = await apiClient.get('/admin/activity-log', { 
      params: { ...tableParams, ...filterParams.value } 
    });
    logData.value = res.data.data || [];
    totalLogs.value = res.data.total || 0;
  } catch (err) {
    console.error("Gagal load activity log", err);
  } finally {
    isLoading.value = false;
  }
};

const openDetail = async (row: any) => {
  if (row.source_table === 'submission') {
    try {
      const res = await apiClient.get(`/admin/submissions/${row.id}`);
      selectedDetail.value = res.data.data;
      selectedDetail.value.type = 'submission';
    } catch (err: any) {
      console.error(err);
      alert("Gagal memuat detail pengajuan.");
    }
  } else if (row.source_table === 'redemption') {
    try {
      const res = await apiClient.get(`/admin/redemptions/${row.id}`);
      selectedDetail.value = res.data.data;
      selectedDetail.value.type = 'redemption';
    } catch (err: any) {
      console.error(err);
      alert("Gagal memuat detail penukaran.");
    }
  }
};

const closeDetail = () => {
  selectedDetail.value = null;
};

const onTableChange = (params: any) => {
  tableParams = params;
  fetchLogs();
};

const applyFilter = () => {
  tableParams.page = 1; // Reset to page 1 on filter
  fetchLogs();
};

let searchTimeout: ReturnType<typeof setTimeout> | null = null;
const onSearchInput = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    applyFilter();
  }, 500);
};

const resetFilter = () => {
  filterParams.value = {
    search: '',
    type: '',
    date_from: '',
    date_to: ''
  };
  applyFilter();
};

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Jakarta' }).replace(/\./g, ':');
};

onMounted(() => {
  if (route.query.search) {
    filterParams.value.search = route.query.search as string;
  }
  fetchLogs();
});
</script>
