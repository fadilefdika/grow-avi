<template>
  <div class="min-h-screen bg-gray-50 flex">
    
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-1 p-8 overflow-auto h-screen">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Manajemen Hadiah</h1>
          <p class="text-gray-500 mt-1">Kelola katalog hadiah dan pantau riwayat penukaran</p>
        </div>
        <!-- <div class="text-right">
          <div class="inline-flex items-center px-4 py-2 bg-blue-50 text-primary rounded-lg border border-blue-100 font-medium">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            Nilai Tukar: 100 pts = Rp 50.000
          </div>
        </div> -->
      </div>


        <DataTable
          title="Katalog Hadiah"
          :columns="rewardColumns"
          :rows="rewards"
          :server-side="true"
          :total-rows="totalRewards"
          @change="onRewardTableChange"
          empty-text="Katalog hadiah kosong."
        >
          <template #action>
            <button @click="openRewardModal()" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors text-sm font-medium whitespace-nowrap">
              + Tambah Hadiah
            </button>
          </template>
          <template #cell-points_required="{ value }">
            <span class="font-bold text-primary">{{ value }} pts</span>
          </template>


          <template #row-actions="{ row }">
            <button @click="openRewardModal(row)" class="text-blue-600 hover:text-blue-900 mr-3">Edit</button>
            <button @click="deleteReward(row.id)" class="text-red-600 hover:text-red-900">Hapus</button>
          </template>
        </DataTable>


    </div>

    <!-- Reward Modal -->
    <div v-if="showRewardModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">{{ editingReward ? 'Edit Hadiah' : 'Tambah Hadiah' }}</h3>
          <button @click="showRewardModal = false" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Hadiah</label>
            <input type="text" v-model="rewardForm.title" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi (Opsional)</label>
            <textarea v-model="rewardForm.description" rows="2" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Poin Diperlukan</label>
            <input type="number" min="0" @wheel.prevent v-model="rewardForm.points_required" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none" />
          </div>

          
          <div class="flex justify-end gap-3 mt-6">
            <button @click="showRewardModal = false" class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors font-medium text-sm">Batal</button>
            <button @click="saveReward" :disabled="isProcessing || !rewardForm.title || rewardForm.points_required <= 0" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors font-medium text-sm disabled:opacity-50">Simpan</button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';
import DataTable from '../../components/ui/DataTable.vue';
import { useToast } from 'vue-toastification';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const rewardColumns = [
  { key: 'title', label: 'Nama Hadiah' },
  { key: 'points_required', label: 'Poin Diperlukan' }
];

const rewards = ref<any[]>([]);
const totalRewards = ref(0);

const rewardParams = ref({ page: 1, limit: 10, search: '', sort: '', order: 'asc' });
const isProcessing = ref(false);

const showRewardModal = ref(false);
const editingReward = ref<any>(null);
const rewardForm = ref({ title: '', description: '', points_required: 0, reward_type: 'CUSTOM' });

const fetchRewards = async () => {
  try {
    const res = await apiClient.get('/rewards', { params: rewardParams.value });
    rewards.value = res.data.data || [];
    totalRewards.value = res.data.total || 0;
  } catch (err) {
    console.error("Gagal memuat reward", err);
  }
};

const onRewardTableChange = (params: any) => {
  rewardParams.value = params;
  fetchRewards();
};



const openRewardModal = (reward?: any) => {
  if (reward) {
    editingReward.value = reward;
    rewardForm.value = {
      title: reward.title,
      description: reward.description || '',
      points_required: reward.points_required,
      reward_type: reward.reward_type || 'CUSTOM'
    };
  } else {
    editingReward.value = null;
    rewardForm.value = { title: '', description: '', points_required: 0, reward_type: 'CUSTOM' };
  }
  showRewardModal.value = true;
};

const saveReward = async () => {
  isProcessing.value = true;
  try {
    if (editingReward.value) {
      await apiClient.put(`/admin/rewards/${editingReward.value.id}`, rewardForm.value);
      toast.success('Hadiah berhasil diupdate');
    } else {
      await apiClient.post('/admin/rewards', rewardForm.value);
      toast.success('Hadiah berhasil ditambahkan');
    }
    showRewardModal.value = false;
    fetchRewards();
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal menyimpan hadiah');
  } finally {
    isProcessing.value = false;
  }
};

const deleteReward = async (id: number) => {
  if (!confirm('Yakin ingin menghapus hadiah ini?')) return;
  try {
    await apiClient.delete(`/admin/rewards/${id}`);
    toast.success('Hadiah berhasil dihapus');
    fetchRewards();
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal menghapus hadiah');
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Jakarta' }).replace(/\./g, ':');
};

onMounted(() => {
  fetchRewards();
});
</script>
