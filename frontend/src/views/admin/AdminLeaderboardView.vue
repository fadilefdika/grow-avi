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
          <!-- Kategori -->
          <div class="flex-1 relative">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Kategori</label>
            <button @click="toggleCategoryDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-sm font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showCategoryDropdown}">
              <span class="truncate pr-2">{{ getCategoryName(filterCategory) }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showCategoryDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showCategoryDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto p-1">
              <div @click="selectCategory('')" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer" :class="!filterCategory ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Kategori</div>
              <div v-for="cat in categories" :key="cat.id" @click="selectCategory(cat.id)" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate" :class="filterCategory === cat.id ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ cat.name }}</div>
            </div>
          </div>
          
          <!-- Departemen -->
          <div class="flex-1 relative">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Departemen</label>
            <button @click="toggleDepartmentDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-sm font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showDepartmentDropdown}">
              <span class="truncate pr-2">{{ filterDepartment || 'Semua Departemen' }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showDepartmentDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showDepartmentDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto p-1">
              <div @click="selectDepartment('')" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer" :class="!filterDepartment ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Departemen</div>
              <div v-for="dept in availableDepartments" :key="dept" @click="selectDepartment(dept)" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate" :class="filterDepartment === dept ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ dept }}</div>
            </div>
          </div>
          <!-- Aktivitas -->
          <div class="flex-1 relative">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Aktivitas</label>
            <button @click="toggleActivityDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-sm font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showActivityDropdown}">
              <span class="truncate pr-2">{{ getActivityName(filterActivity) }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showActivityDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showActivityDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto p-1">
              <div @click="selectActivity('')" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer" :class="!filterActivity ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Aktivitas</div>
              <div v-for="act in filteredActivitiesForDropdown" :key="act.id" @click="selectActivity(act)" class="px-3 py-2.5 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate" :class="filterActivity === act.id ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ act.name }}</div>
            </div>
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
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">Peringkat</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Karyawan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Departemen</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Total Poin</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="entry in leaderboard" :key="entry.npk" @click="$router.push(`/admin/activity-log?search=${entry.npk}`)" class="hover:bg-gray-50 transition-colors cursor-pointer">
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
                <span class="text-lg font-bold text-primary">{{ entry.balance }} poin</span>
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
const activities = ref<any[]>([]);
const availableDepartments = ref<string[]>([]);
const filterCategory = ref('');
const filterActivity = ref('');
const filterDepartment = ref('');
const isLoading = ref(true);

const showCategoryDropdown = ref(false);
const showDepartmentDropdown = ref(false);
const showActivityDropdown = ref(false);

import { computed } from 'vue';

const filteredActivitiesForDropdown = computed(() => {
  if (!filterCategory.value) return activities.value;
  return activities.value.filter(a => a.category_id === filterCategory.value);
});

const toggleCategoryDropdown = () => {
  showCategoryDropdown.value = !showCategoryDropdown.value;
  showDepartmentDropdown.value = false;
  showActivityDropdown.value = false;
};

const toggleDepartmentDropdown = () => {
  showDepartmentDropdown.value = !showDepartmentDropdown.value;
  showCategoryDropdown.value = false;
  showActivityDropdown.value = false;
};

const toggleActivityDropdown = () => {
  showActivityDropdown.value = !showActivityDropdown.value;
  showCategoryDropdown.value = false;
  showDepartmentDropdown.value = false;
};

const selectCategory = (id: string) => {
  filterCategory.value = id;
  if (filterActivity.value) {
    const act = activities.value.find((a: any) => a.id === filterActivity.value);
    if (act && act.category_id !== filterCategory.value && filterCategory.value !== '') {
      filterActivity.value = '';
    }
  }
  showCategoryDropdown.value = false;
  fetchLeaderboard();
};

const selectActivity = (act: any) => {
  if (!act) {
    filterActivity.value = '';
  } else {
    filterActivity.value = act.id;
    if (act.category_id) {
      filterCategory.value = act.category_id;
    }
  }
  showActivityDropdown.value = false;
  fetchLeaderboard();
};

const selectDepartment = (dept: string) => {
  filterDepartment.value = dept;
  showDepartmentDropdown.value = false;
  fetchLeaderboard();
};

const getCategoryName = (id: string) => {
  if (!id) return 'Semua Kategori';
  const c = categories.value.find((c: any) => c.id == id);
  return c ? c.name : 'Semua Kategori';
};

const getActivityName = (id: string) => {
  if (!id) return 'Semua Aktivitas';
  const act = activities.value.find((c: any) => c.id === id);
  return act ? act.name : 'Semua Aktivitas';
};

// Close dropdowns when clicking outside (simple implementation by listening on window)
if (typeof window !== 'undefined') {
  window.addEventListener('click', (e) => {
    const target = e.target as HTMLElement;
    if (!target.closest('.relative')) {
      showCategoryDropdown.value = false;
      showDepartmentDropdown.value = false;
      showActivityDropdown.value = false;
    }
  });
}

const fetchCategories = async () => {
  try {
    const [catRes, actRes] = await Promise.all([
      apiClient.get('/categories'),
      apiClient.get('/activities')
    ]);
    categories.value = catRes.data.data || [];
    activities.value = actRes.data.data || [];
  } catch (err) {
    console.error("Gagal load master data", err);
  }
};

const fetchLeaderboard = async () => {
  isLoading.value = true;
  try {
    let url = '/leaderboard';
    const params = new URLSearchParams();
    if (filterCategory.value) params.append('category_id', filterCategory.value);
    if (filterActivity.value) params.append('activity_id', filterActivity.value);
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
