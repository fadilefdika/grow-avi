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
      <main class="flex-1 overflow-y-auto p-8">
        
        <!-- Filter Panel -->
        <div class="bg-white p-5 rounded-2xl border border-gray-200 shadow-sm mb-6 flex flex-wrap gap-4 items-end">
          <div class="flex-1 min-w-[200px]">
            <label class="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Cari GROW ID</label>
            <input type="text" v-model="filterParams.search" placeholder="Contoh: GR2609140001" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @keyup.enter="applyFilter" />
          </div>
          
          <div>
            <label class="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Jenis Kegiatan</label>
            <select v-model="filterParams.type" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="applyFilter">
              <option value="">Semua Jenis</option>
              <option value="APPROVED">Approval Poin</option>
              <option value="REJECTED">Penolakan (Reject)</option>
              <option value="REDEEM">Klaim Poin (Redeem)</option>
            </select>
          </div>
          
          <div>
            <label class="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Dari Tanggal</label>
            <input type="date" v-model="filterParams.date_from" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="applyFilter" />
          </div>
          
          <div>
            <label class="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Sampai Tanggal</label>
            <input type="date" v-model="filterParams.date_to" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:bg-white transition-colors" @change="applyFilter" />
          </div>
          
          <div>
            <button @click="resetFilter" class="px-4 py-2 text-gray-600 hover:text-gray-900 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors font-medium">Reset</button>
          </div>
        </div>

        <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
          <DataTable
            :columns="[
              { key: 'grow_id', label: 'GROW ID', sortable: true },
              { key: 'npk', label: 'NPK', sortable: false },
              { key: 'user_name', label: 'Nama', sortable: true },
              { key: 'activity_type', label: 'Jenis', sortable: false },
              { key: 'points', label: 'Poin', sortable: true },
              { key: 'activity_date', label: 'Waktu Transaksi', sortable: true }
            ]"
            :data="logData"
            :total="totalLogs"
            :loading="isLoading"
            @change="onTableChange"
          >
            <template #activity_type="{ row }">
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
            
            <template #points="{ row }">
              <span :class="row.activity_type === 'REDEEM' ? 'text-red-600 font-bold' : (row.activity_type === 'APPROVED' ? 'text-green-600 font-bold' : 'text-gray-500 font-bold')">
                {{ row.activity_type === 'REDEEM' ? '-' : (row.activity_type === 'APPROVED' ? '+' : '') }}{{ row.points }}
              </span>
            </template>

            <template #activity_date="{ row }">
              {{ formatDate(row.activity_date) }} WIB
            </template>
          </DataTable>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';
import DataTable from '../../components/ui/DataTable.vue';

const logData = ref<any[]>([]);
const totalLogs = ref(0);
const isLoading = ref(true);

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

const onTableChange = (params: any) => {
  tableParams = params;
  fetchLogs();
};

const applyFilter = () => {
  tableParams.page = 1; // Reset to page 1 on filter
  fetchLogs();
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
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
};

onMounted(() => {
  fetchLogs();
});
</script>
