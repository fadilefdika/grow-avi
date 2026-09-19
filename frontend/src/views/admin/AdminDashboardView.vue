<template>
  <div class="min-h-screen bg-gray-50 flex">
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-1 p-8 overflow-auto h-screen">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dasbor Admin</h1>
          <p class="text-gray-500 mt-1">Selamat datang,  Admin</p>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- Card 1: Pending (most urgent) -->
        <div class="bg-white rounded-xl p-6 shadow-sm border border-orange-100 flex items-center justify-between cursor-pointer hover:border-orange-400 transition-colors" @click="$router.push('/admin/submissions')">
          <div>
            <p class="text-sm font-medium text-gray-500 mb-1">Menunggu Approval</p>
            <p class="text-3xl font-bold" :class="stats.pendingSubmissions > 0 ? 'text-orange-600' : 'text-gray-400'">{{ stats.pendingSubmissions }}</p>
            <p class="text-xs text-gray-400 mt-1">pengajuan perlu ditinjau</p>
          </div>
          <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="stats.pendingSubmissions > 0 ? 'bg-orange-50 text-orange-600' : 'bg-gray-50 text-gray-400'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
        </div>
        
        <!-- Card 2: Approved this month -->
        <div class="bg-white rounded-xl p-6 shadow-sm border border-gray-100 flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 mb-1">Disetujui Bulan Ini</p>
            <p class="text-3xl font-bold text-green-600">{{ stats.approvedThisMonth }}</p>
            <p class="text-xs text-gray-400 mt-1">dari {{ stats.approvedThisMonth + stats.rejectedThisMonth }} pengajuan diproses</p>
          </div>
          <div class="w-12 h-12 bg-green-50 rounded-full flex items-center justify-center text-green-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
        </div>

        <!-- Card 3: Rejection rate -->
        <div class="bg-white rounded-xl p-6 shadow-sm border border-gray-100 flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 mb-1">Tingkat Penolakan</p>
            <div class="flex items-baseline gap-1">
              <p class="text-3xl font-bold" :class="stats.rejectionRate > 30 ? 'text-red-600' : stats.rejectionRate > 10 ? 'text-yellow-600' : 'text-gray-700'">{{ stats.rejectionRate.toFixed(1) }}</p>
              <span class="text-lg font-bold text-gray-400">%</span>
            </div>
            <p class="text-xs mt-1" :class="stats.rejectionRate > 30 ? 'text-red-500 font-medium' : 'text-gray-400'">
              {{ stats.rejectionRate > 30 ? '⚠️ Di atas batas wajar' : stats.rejectionRate > 10 ? 'Perlu diperhatikan' : 'Dalam batas normal' }}
            </p>
          </div>
          <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="stats.rejectionRate > 30 ? 'bg-red-50 text-red-600' : stats.rejectionRate > 10 ? 'bg-yellow-50 text-yellow-600' : 'bg-gray-50 text-gray-500'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
          <h2 class="text-lg font-bold text-gray-800">5 Karyawan Teratas</h2>
          <router-link to="/admin/leaderboard" class="text-sm text-primary font-medium hover:underline">Lihat Lengkap</router-link>
        </div>
        <div v-if="isLoading" class="p-8 text-center text-gray-500">Memuat data...</div>
        <div v-else-if="stats.topUsers.length === 0" class="p-8 text-center text-gray-500">Belum ada data poin.</div>
        <ul v-else class="divide-y divide-gray-100">
          <li v-for="(user, index) in stats.topUsers" :key="user.npk" class="px-6 py-4 flex items-center justify-between hover:bg-gray-50">
            <div class="flex items-center gap-4">
              <span class="w-8 text-center font-bold text-gray-400">#{{ index + 1 }}</span>
              <div>
                <p class="text-sm font-bold text-gray-900">{{ user.user_name }}</p>
                <p class="text-xs text-gray-500">{{ user.npk }} - {{ user.department }}</p>
              </div>
            </div>
            <div class="text-right">
              <span class="text-lg font-bold text-primary">{{ user.balance }} poin</span>
            </div>
          </li>
        </ul>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';

const router = useRouter();
const authStore = useAuthStore();

const isLoading = ref(true);
const stats = ref({
  pendingSubmissions: 0,
  approvedThisMonth: 0,
  rejectedThisMonth: 0,
  rejectionRate: 0,
  totalRedemptionsMonth: 0,
  topUsers: [] as any[]
});

const fetchStats = async () => {
  isLoading.value = true;
  try {
    const res = await apiClient.get('/admin/users/stats');
    const d = res.data.data;
    stats.value = {
      pendingSubmissions: d.pending_submissions || 0,
      approvedThisMonth: d.approved_this_month || 0,
      rejectedThisMonth: d.rejected_this_month || 0,
      rejectionRate: d.rejection_rate || 0,
      totalRedemptionsMonth: d.total_redemptions_month || 0,
      topUsers: d.top_users || []
    };
  } catch (err) {
    console.error("Gagal load stats", err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchStats();
});
</script>
