<template>
  <div class="min-h-screen bg-gray-50 flex">
    
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-1 p-8 overflow-auto h-screen">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Leaderboard Keseluruhan</h1>
          <p class="text-gray-500 mt-1">Peringkat poin GROW karyawan</p>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-xl p-6 shadow-sm border border-gray-100 mb-6">
        <div class="flex gap-4 max-w-2xl">
          <div class="flex-1">
            <label class="block text-xs font-medium text-gray-500 mb-2">Filter Kategori</label>
            <select v-model="filterCategory" @change="fetchLeaderboard" class="w-full bg-white border border-gray-300 text-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary">
              <option value="">Semua Kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div class="flex-1">
            <label class="block text-xs font-medium text-gray-500 mb-2">Filter Departemen</label>
            <select v-model="filterDepartment" @change="fetchLeaderboard" class="w-full bg-white border border-gray-300 text-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary">
              <option value="">Semua Departemen</option>
              <option v-for="dept in availableDepartments" :key="dept" :value="dept">{{ dept }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12 bg-white rounded-xl shadow-sm border border-gray-100">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Empty State -->
      <div v-else-if="leaderboard.length === 0" class="bg-white rounded-xl p-12 shadow-sm border border-gray-100 text-center">
        <p class="text-gray-500 font-medium">Belum ada data poin.</p>
      </div>

      <!-- Leaderboard List -->
      <div v-else class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">Rank</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Karyawan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Departemen</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Total Poin</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="entry in leaderboard" :key="entry.npk" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="font-bold text-lg" :class="getRankColor(entry.rank)">#{{ entry.rank }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="font-bold text-gray-900">{{ entry.user_name }}</div>
                <div class="text-xs text-gray-500">{{ entry.npk }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                {{ entry.department }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right">
                <span class="text-lg font-bold text-primary">{{ entry.balance }} pts</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';

const leaderboard = ref<any[]>([]);
const categories = ref<any[]>([]);
const availableDepartments = ref<string[]>([]);
const filterCategory = ref('');
const filterDepartment = ref('');
const isLoading = ref(true);

const fetchCategories = async () => {
  try {
    const res = await apiClient.get('/categories');
    categories.value = res.data.data || [];
  } catch (err) {
    console.error("Gagal load kategori", err);
  }
};

const fetchLeaderboard = async () => {
  isLoading.value = true;
  try {
    let url = '/leaderboard';
    const params = new URLSearchParams();
    if (filterCategory.value) params.append('category_id', filterCategory.value);
    if (filterDepartment.value) params.append('department', filterDepartment.value);
    
    if (params.toString()) {
      url += '?' + params.toString();
    }
    
    const res = await apiClient.get(url);
    leaderboard.value = res.data.data || [];
    
    if (availableDepartments.value.length === 0) {
      const depts = new Set<string>();
      leaderboard.value.forEach(item => {
        if (item.department) depts.add(item.department);
      });
      availableDepartments.value = Array.from(depts).sort();
    }
  } catch (err) {
    console.error("Gagal load leaderboard", err);
  } finally {
    isLoading.value = false;
  }
};

const getRankColor = (rank: number) => {
  if (rank === 1) return 'text-yellow-500';
  if (rank === 2) return 'text-gray-400';
  if (rank === 3) return 'text-amber-700';
  return 'text-gray-500';
};

onMounted(() => {
  fetchCategories();
  fetchLeaderboard();
});
</script>
