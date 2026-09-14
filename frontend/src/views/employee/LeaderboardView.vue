<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Top Header -->
    <div class="bg-primary pt-12 pb-24 px-6 rounded-b-[2rem] shadow-md relative">
      <div class="flex justify-between items-center relative z-10">
        <div>
          <button @click="$router.push('/dashboard')" class="text-white mb-2 flex items-center text-sm opacity-80 hover:opacity-100">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali
          </button>
          <h1 class="text-white text-xl font-bold">Leaderboard</h1>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="px-5 -mt-16 relative z-20 space-y-5">
      
      <!-- Filters -->
      <div class="bg-white rounded-2xl p-5 shadow-lg border border-gray-100">
        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-xs font-medium text-gray-500 mb-1">Kategori</label>
            <select v-model="filterCategory" @change="fetchLeaderboard" class="w-full bg-gray-50 border border-gray-200 text-gray-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary">
              <option value="">Semua Kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div class="flex-1">
            <label class="block text-xs font-medium text-gray-500 mb-1">Departemen</label>
            <select v-model="filterDepartment" @change="fetchLeaderboard" class="w-full bg-gray-50 border border-gray-200 text-gray-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary">
              <option value="">Semua Departemen</option>
              <option v-for="dept in availableDepartments" :key="dept" :value="dept">{{ dept }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Empty State -->
      <div v-else-if="leaderboard.length === 0" class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 text-center">
        <p class="text-gray-500 font-medium">Belum ada data poin.</p>
      </div>

      <!-- Leaderboard List -->
      <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="divide-y divide-gray-100">
          <div 
            v-for="entry in leaderboard" 
            :key="entry.npk" 
            :class="['p-4 flex items-center gap-4 transition-colors', entry.npk === currentUserNpk ? 'bg-blue-50/50' : 'hover:bg-gray-50']"
          >
            <!-- Rank Badge -->
            <div class="flex-shrink-0 w-8 text-center font-extrabold text-lg" :class="getRankColor(entry.rank)">
              #{{ entry.rank }}
            </div>
            
            <!-- User Info -->
            <div class="flex-1 overflow-hidden">
              <h3 class="font-bold text-gray-900 text-sm truncate flex items-center gap-2">
                {{ entry.user_name }}
                <span v-if="entry.npk === currentUserNpk" class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-primary text-white">
                  Anda
                </span>
              </h3>
              <p class="text-xs text-gray-500 truncate">{{ entry.department }}</p>
            </div>
            
            <!-- Points -->
            <div class="text-right flex-shrink-0">
              <span class="text-lg font-bold text-primary">{{ entry.balance }}</span>
              <span class="text-xs font-medium text-gray-500 ml-1">pts</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';

const authStore = useAuthStore();
const currentUserNpk = computed(() => authStore.user?.npk);

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
    
    // Extract unique departments for the dropdown if not already set
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
