<template>
  <div class="min-h-screen bg-gray-50 flex">
    
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-1 p-8 overflow-auto h-screen">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Antrean Approval</h1>
          <p class="text-gray-500 mt-1">Verifikasi pengajuan poin GROW dari karyawan</p>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12 bg-white rounded-xl shadow-sm border border-gray-100">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <div v-else>
        <DataTable
          title="Antrean Pengajuan"
          subtitle="Pengajuan dengan status PENDING"
          :columns="submissionColumns"
          :rows="pendingQueue"
          :server-side="true"
          :total-rows="totalSubmissions"
          @change="onQueueTableChange"
          empty-text="Tidak ada antrean pengajuan PENDING saat ini."
        >
          <template #cell-created_at="{ row }">
            <div>{{ formatDate(row.created_at) }}</div>
            <div class="text-xs text-gray-500">{{ row.grow_id }}</div>
          </template>
          <template #cell-user_name="{ row }">
            <div class="font-medium text-gray-900">{{ row.user_name }}</div>
            <div class="text-xs text-gray-500">{{ row.npk }} - {{ row.department }}</div>
          </template>
          <template #cell-activity_name="{ row }">
            <div class="font-medium text-gray-900">{{ row.activity_name }}</div>
            <div class="text-xs text-gray-500">{{ row.category_name }}</div>
            <div v-if="row.custom_reference?.Valid" class="text-xs text-gray-500 italic">Ref: {{ row.custom_reference.String }}</div>
          </template>
          <template #cell-default_points="{ value }">
            <span class="px-2 py-1 bg-blue-50 text-blue-700 rounded-md font-bold text-xs">{{ value }} pts</span>
          </template>
          <template #actions="{ row }">
            <button @click="openDetail(row.id)" class="text-primary hover:text-blue-900 font-semibold bg-blue-50 hover:bg-blue-100 px-3 py-1.5 rounded-lg transition-colors">Review</button>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Review Modal -->
    <div v-if="selectedDetail" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm overflow-y-auto">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-4xl overflow-hidden flex flex-col max-h-[90vh]">
        
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <div>
            <h3 class="text-lg font-bold text-gray-900">Review Pengajuan: {{ selectedDetail.grow_id }}</h3>
          </div>
          <button @click="closeDetail" class="text-gray-400 hover:text-gray-600">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        
        <!-- Modal Body -->
        <div class="p-6 flex-1 overflow-y-auto">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            
            <!-- Left: Info -->
            <div class="space-y-6">
              <div>
                <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Informasi Karyawan</h4>
                <div class="bg-gray-50 p-4 rounded-xl border border-gray-100">
                  <p class="font-bold text-gray-900">{{ selectedDetail.user_name }}</p>
                  <p class="text-sm text-gray-600">{{ selectedDetail.npk }} — {{ selectedDetail.department }}</p>
                </div>
              </div>
              
              <div>
                <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Detail Aktivitas</h4>
                <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 space-y-3">
                  <div>
                    <p class="text-xs text-gray-500">Kategori & Aktivitas</p>
                    <p class="font-medium text-xs text-gray-900">
                      {{ selectedDetail.category_name }}<span v-if="selectedDetail.category_name !== 'INNOVATION'"> - {{ selectedDetail.activity_name === 'Yang lain (Custom)' && selectedDetail.custom_activity_type?.Valid ? selectedDetail.custom_activity_type.String : selectedDetail.activity_name }}</span>
                    </p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Tanggal Pelaksanaan</p>
                    <p class="font-medium text-xs text-gray-900">{{ formatDateOnly(selectedDetail.activity_date) }}</p>
                  </div>
                  <div v-if="selectedDetail.custom_reference?.Valid">
                    <p class="text-xs text-gray-500">Nama Kegiatan/Aktivitas</p>
                    <p class="font-medium text-xs text-gray-900">{{ selectedDetail.custom_reference.String }}</p>
                  </div>
                  <div v-if="selectedDetail.nomor_ss?.Valid">
                    <p class="text-xs text-gray-500">Nomor SS</p>
                    <p class="font-medium text-xs text-gray-900">{{ selectedDetail.nomor_ss.String }}</p>
                  </div>
                </div>
              </div>

              <div>
                <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Penilaian Poin</h4>
                <div class="bg-blue-50 p-4 rounded-xl border border-blue-100">
                  <label class="block text-sm font-medium text-blue-900 mb-2">Poin yang Diberikan</label>
                  <div class="flex items-center">
                    <input type="number" v-model="pointsOverride" @wheel.prevent class="w-24 px-3 py-2 border border-blue-200 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary text-lg font-bold text-primary [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none" />
                    <span class="ml-2 text-sm font-medium text-blue-700">pts (Base: {{ selectedDetail.default_points }})</span>
                  </div>
                  <p class="text-xs text-blue-600 mt-2">Anda bisa mengubah poin default jika aktivitas ini bersifat dinamis (misal: Inovasi).</p>
                </div>
              </div>

              <div v-if="actionType === 'reject'">
                <h4 class="text-xs font-bold text-red-500 uppercase tracking-wider mb-2">Alasan Penolakan</h4>
                <textarea v-model="adminNotes" rows="3" class="w-full px-3 py-2 border border-red-300 rounded-lg shadow-sm focus:outline-none text-xs focus:ring-1 focus:ring-red-500 focus:border-red-500" placeholder="Berikan catatan agar karyawan bisa memperbaiki..."></textarea>
              </div>

            </div>

            <!-- Right: Evidence -->
            <div>
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">Bukti Lampiran ({{ selectedDetail.evidence.length }})</h4>
              <div v-if="selectedDetail.evidence.length === 0" class="bg-gray-50 p-8 rounded-xl border border-gray-100 text-center text-gray-500 italic">
                Tidak ada lampiran.
              </div>
              <div v-else class="space-y-4">
                <div v-for="(ev, idx) in selectedDetail.evidence" :key="idx" class="border border-gray-200 rounded-xl overflow-hidden bg-gray-50 relative group">
                  <a :href="ev" target="_blank" class="block relative">
                    <img :src="ev" class="w-full h-auto max-h-64 object-contain mx-auto" />
                    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <span class="text-white bg-black/70 px-3 py-1 rounded-full text-sm font-medium flex items-center">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 21h7a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v11m0 5l4.879-4.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242z"></path></svg>
                        Buka Penuh
                      </span>
                    </div>
                  </a>
                </div>
              </div>
            </div>
            
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-4 border-t border-gray-100 bg-gray-50 flex items-center justify-between">
          <div>
             <span v-if="error" class="text-sm font-medium text-red-600">{{ error }}</span>
          </div>
          <div class="flex gap-3">
            <template v-if="actionType === 'reject'">
              <button @click="actionType = null" class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors font-medium text-sm">Batal</button>
              <button @click="submitReject" :disabled="!adminNotes || isProcessing" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium text-sm disabled:opacity-50">Kirim Penolakan</button>
            </template>
            <template v-else>
              <button @click="actionType = 'reject'" class="px-4 py-2 border border-red-200 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors font-medium text-sm">Tolak (Reject)</button>
              <button @click="submitApprove" :disabled="isProcessing" class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium text-sm">Setujui (Approve)</button>
            </template>
          </div>
        </div>

      </div>
    </div>

    <!-- Zero Points Confirmation Dialog -->
    <div v-if="showZeroPointsConfirm" class="fixed inset-0 z-[60] flex items-center justify-center p-4 bg-gray-900/70 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-yellow-100 flex items-center justify-center flex-shrink-0">
            <svg class="w-5 h-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div>
            <h3 class="font-bold text-gray-900">Konfirmasi 0 Poin</h3>
            <p class="text-sm text-gray-500">Anda akan memberikan 0 poin untuk pengajuan ini.</p>
          </div>
        </div>
        <p class="text-sm text-gray-600 mb-6">Apakah Anda yakin ingin menyetujui pengajuan ini dengan <span class="font-bold text-red-600">0 poin</span>? Karyawan akan menerima poin tersebut.</p>
        <div class="flex gap-3">
          <button @click="showZeroPointsConfirm = false" class="flex-1 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium text-sm">Batal, Ubah Poin</button>
          <button @click="confirmApproveZero" class="flex-1 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition-colors font-medium text-sm">Ya, Tetap 0 Poin</button>
        </div>
      </div>
    </div>

  </div>
</template>



<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { useToast } from 'vue-toastification';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';
import DataTable from '../../components/ui/DataTable.vue';

const router = useRouter();
const authStore = useAuthStore();

const submissionColumns = [
  { key: 'created_at', label: 'Tgl Submit' },
  { key: 'user_name', label: 'Karyawan' },
  { key: 'activity_name', label: 'Aktivitas' },
  { key: 'default_points', label: 'Poin Base' }
];

const pendingQueue = ref<any[]>([]);
const totalSubmissions = ref(0);
const queueParams = ref({ page: 1, limit: 10, search: '', sort: '', order: 'asc' });
const isLoading = ref(true);

const selectedDetail = ref<any>(null);
const pointsOverride = ref<number>(0);
const actionType = ref<'reject' | null>(null);
const adminNotes = ref('');
const isProcessing = ref(false);
const error = ref('');
const toast = useToast();
const showZeroPointsConfirm = ref(false);

const fetchQueue = async () => {
  isLoading.value = true;
  try {
    const res = await apiClient.get('/admin/submissions/pending', { params: queueParams.value });
    pendingQueue.value = res.data.data || [];
    totalSubmissions.value = res.data.total || 0;
  } catch (err) {
    console.error("Gagal load antrean", err);
  } finally {
    isLoading.value = false;
  }
};

const onQueueTableChange = (params: any) => {
  queueParams.value = params;
  fetchQueue();
};

const openDetail = async (id: number) => {
  try {
    const res = await apiClient.get(`/admin/submissions/${id}`);
    selectedDetail.value = res.data.data;
    console.log(selectedDetail.value.evidence);
    
    pointsOverride.value = res.data.data.default_points;
    actionType.value = null;
    adminNotes.value = '';
    error.value = '';
  } catch (err: any) {
    console.error(err);
    const errorMsg = err.response?.data?.error || "Gagal load detail pengajuan";
    alert(errorMsg);
  }
};

const closeDetail = () => {
  selectedDetail.value = null;
};

const submitApprove = async () => {
  if (!selectedDetail.value) return;
  // Show zero-points confirmation if needed
  if (Number(pointsOverride.value) === 0) {
    showZeroPointsConfirm.value = true;
    return;
  }
  await doApprove();
};

const confirmApproveZero = async () => {
  showZeroPointsConfirm.value = false;
  await doApprove();
};

const doApprove = async () => {
  if (!selectedDetail.value) return;
  isProcessing.value = true;
  error.value = '';
  try {
    await apiClient.post(`/admin/submissions/${selectedDetail.value.id}/approve`, {
      points_override: pointsOverride.value
    });
    toast.success('Pengajuan berhasil di-approve!');
    closeDetail();
    fetchQueue();
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || 'Terjadi kesalahan saat approve';
    error.value = errorMsg;
    toast.error(errorMsg);
  } finally {
    isProcessing.value = false;
  }
};

const submitReject = async () => {
  if (!selectedDetail.value) return;
  isProcessing.value = true;
  error.value = '';
  try {
    await apiClient.post(`/admin/submissions/${selectedDetail.value.id}/reject`, {
      admin_notes: adminNotes.value
    });
    toast.success('Pengajuan telah ditolak.');
    closeDetail();
    fetchQueue();
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || 'Terjadi kesalahan saat reject';
    error.value = errorMsg;
    toast.error(errorMsg);
  } finally {
    isProcessing.value = false;
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
};

// Date-only formatter (for activity_date which has no meaningful time component)
const formatDateOnly = (dateString: string) => {
  if (!dateString) return '';
  // Parse just the date portion to avoid timezone shifting
  const [year, month, day] = dateString.split('T')[0].split('-').map(Number);
  const date = new Date(year, month - 1, day);
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
};

onMounted(() => {
  fetchQueue();
});
</script>
