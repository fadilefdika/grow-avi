<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Hero Header -->
    <div class="hero-gradient pt-12 pb-32 px-6 relative overflow-hidden">
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-8 -right-8 w-40 h-40 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFD700, transparent);"></div>
        <div class="absolute top-16 -left-6 w-28 h-28 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #1CB0F6, transparent); animation-delay: 1.2s;"></div>
      </div>
      <div class="flex justify-between items-center relative z-10">
        <div>
          <button @click="$router.push('/dashboard')" class="flex items-center gap-1.5 text-white/70 text-sm font-bold mb-3 hover:text-white transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali
          </button>
          <h1 class="text-white text-2xl font-black">🏆 Leaderboard</h1>
          <p class="text-white/60 text-xs font-semibold mt-0.5">Siapa yang teratas bulan ini?</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="px-4 -mt-24 relative z-20 space-y-4">

      <!-- Podium Top 3 -->
      <div v-if="!isLoading && top3.length > 0" class="bg-white rounded-3xl shadow-xl border border-gray-100 p-4 overflow-hidden">
        <h2 class="text-xs font-black text-gray-500 uppercase tracking-wider text-center mb-4">🏆 Top 3 Terbaik</h2>

        <!-- Podium Visual — gap & size diturunkan untuk layar kecil -->
        <div class="flex items-end justify-center gap-1.5 xs:gap-3">
          <!-- 2nd Place -->
          <div v-if="top3[1]" class="flex flex-col items-center flex-1 min-w-0">
            <div class="podium-silver w-11 h-11 xs:w-14 xs:h-14 rounded-2xl flex items-center justify-center mb-1.5 shadow-md">
              <span class="text-lg xs:text-2xl font-black text-gray-500">{{ getInitial(top3[1].user_name) }}</span>
            </div>
            <p class="text-[10px] xs:text-xs font-black text-gray-700 truncate w-full text-center px-0.5">{{ formatName(top3[1].user_name) }}</p>
            <p class="text-[9px] xs:text-[10px] font-bold text-gray-400 truncate w-full text-center">{{ top3[1].balance }} pts</p>
            <div class="w-full h-14 xs:h-16 rounded-t-xl mt-1.5 flex items-center justify-center" style="background: linear-gradient(to top, #C0C0C0, #D8D8D8);">
              <span class="text-2xl xs:text-3xl font-black text-white drop-shadow">2</span>
            </div>
          </div>

          <!-- 1st Place -->
          <div v-if="top3[0]" class="flex flex-col items-center flex-1 min-w-0">
            <div class="text-xl xs:text-2xl mb-1 text-center animate-float">👑</div>
            <div class="podium-gold w-12 h-12 xs:w-16 xs:h-16 rounded-2xl flex items-center justify-center mb-1.5 shadow-lg" style="box-shadow: 0 0 16px rgba(255,215,0,0.4);">
              <span class="text-2xl xs:text-3xl font-black text-yellow-700">{{ getInitial(top3[0].user_name) }}</span>
            </div>
            <p class="text-xs xs:text-sm font-black text-gray-800 truncate w-full text-center px-0.5">{{ formatName(top3[0].user_name) }}</p>
            <p class="text-[10px] xs:text-xs font-bold text-yellow-600 text-center">{{ top3[0].balance }} pts</p>
            <div class="w-full h-20 xs:h-24 rounded-t-xl mt-1.5 flex items-center justify-center" style="background: linear-gradient(to top, #FFD700, #FFE766);">
              <span class="text-3xl xs:text-4xl font-black text-white drop-shadow-lg">1</span>
            </div>
          </div>

          <!-- 3rd Place -->
          <div v-if="top3[2]" class="flex flex-col items-center flex-1 min-w-0">
            <div class="podium-bronze w-11 h-11 xs:w-14 xs:h-14 rounded-2xl flex items-center justify-center mb-1.5 shadow-md">
              <span class="text-lg xs:text-2xl font-black text-orange-800">{{ getInitial(top3[2].user_name) }}</span>
            </div>
            <p class="text-[10px] xs:text-xs font-black text-gray-700 truncate w-full text-center px-0.5">{{ formatName(top3[2].user_name) }}</p>
            <p class="text-[9px] xs:text-[10px] font-bold text-gray-400 truncate w-full text-center">{{ top3[2].balance }} pts</p>
            <div class="w-full h-8 xs:h-10 rounded-t-xl mt-1.5 flex items-center justify-center" style="background: linear-gradient(to top, #CD7F32, #D4956A);">
              <span class="text-xl xs:text-2xl font-black text-white drop-shadow">3</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter Card -->
      <div class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100">
        <div class="flex gap-3">
          <!-- Kategori -->
          <div class="flex-1 relative">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Kategori</label>
            <button @click="toggleCategoryDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-[11px] font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showCategoryDropdown}">
              <span class="truncate pr-2">{{ getCategoryName(filterCategory) }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showCategoryDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showCategoryDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto overflow-x-hidden p-1 animate-fade-in">
              <div @click="selectCategory('')" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="!filterCategory ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Kategori</div>
              <div v-for="cat in categories" :key="cat.id" @click="selectCategory(cat.id)" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate transition-colors" :class="filterCategory === cat.id ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ cat.name }}</div>
            </div>
          </div>

          <!-- Departemen -->
          <div class="flex-1 relative hidden xs:block">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Departemen</label>
            <button @click="toggleDepartmentDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-[11px] font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showDepartmentDropdown}">
              <span class="truncate pr-2">{{ filterDepartment || 'Semua Dept.' }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showDepartmentDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showDepartmentDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto overflow-x-hidden p-1 animate-fade-in">
              <div @click="selectDepartment('')" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="!filterDepartment ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Dept.</div>
              <div v-for="dept in availableDepartments" :key="dept" @click="selectDepartment(dept)" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate transition-colors" :class="filterDepartment === dept ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ dept }}</div>
            </div>
          </div>
        </div>
        
        <div class="flex gap-3 mt-3">
          <!-- Aktivitas -->
          <div class="flex-1 relative">
            <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">Aktivitas</label>
            <button @click="toggleActivityDropdown" class="w-full text-left bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-3 py-2.5 text-[11px] font-bold flex justify-between items-center transition-colors hover:border-primary/50" :class="{'border-primary': showActivityDropdown}">
              <span class="truncate pr-2">{{ getActivityName(filterActivity) }}</span>
              <svg class="w-4 h-4 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showActivityDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            
            <div v-if="showActivityDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto overflow-x-hidden p-1 animate-fade-in">
              <div @click="selectActivity('')" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="!filterActivity ? 'text-primary bg-blue-50' : 'text-gray-600'">Semua Aktivitas</div>
              <div v-for="act in filteredActivitiesForDropdown" :key="act.id" @click="selectActivity(act)" class="px-3 py-2.5 text-xs font-bold rounded-lg hover:bg-blue-50 cursor-pointer truncate transition-colors" :class="filterActivity === act.id ? 'text-primary bg-blue-50' : 'text-gray-600'">{{ act.name }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Dropdown Overlay to catch outside clicks -->
      <div v-if="showCategoryDropdown || showDepartmentDropdown" @click="closeDropdowns" class="fixed inset-0 z-40"></div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center py-12">
        <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Empty State -->
      <div v-else-if="leaderboard.length === 0" class="bg-white rounded-2xl p-10 shadow-sm border border-gray-100 text-center">
        <p class="text-3xl mb-3">🏅</p>
        <p class="text-gray-500 font-bold">Belum ada data poin.</p>
        <p class="text-gray-400 text-sm font-medium mt-1">Jadilah yang pertama!</p>
      </div>

      <!-- Full Leaderboard List -->
      <div v-else class="space-y-2">
        <div
          v-for="entry in leaderboard"
          :key="entry.npk"
          :class="[
            'flex items-center gap-4 p-4 rounded-2xl border transition-all',
            entry.npk === currentUserNpk
              ? 'rank-glow'
              : 'bg-white border-gray-100 hover:border-primary/20 hover:shadow-sm'
          ]"
        >
          <!-- Rank Badge -->
          <div class="flex-shrink-0 w-10 text-center">
            <span v-if="entry.rank === 1" class="text-2xl">🥇</span>
            <span v-else-if="entry.rank === 2" class="text-2xl">🥈</span>
            <span v-else-if="entry.rank === 3" class="text-2xl">🥉</span>
            <span v-else class="text-sm font-black" :class="entry.rank <= 10 ? 'text-primary' : 'text-gray-400'">#{{ entry.rank }}</span>
          </div>

          <!-- User Info -->
          <div class="flex-1 overflow-hidden">
            <h3 class="font-black text-gray-900 text-sm truncate flex items-center gap-1.5">
              {{ formatName(entry.user_name) }}
              <span v-if="entry.npk === currentUserNpk" class="inline-flex items-center px-1.5 py-0.5 rounded-lg text-[9px] font-black bg-primary text-white">YOU</span>
            </h3>
            <p class="text-xs font-semibold text-gray-400 truncate">{{ entry.department }}</p>
          </div>
          
          <!-- Points -->
          <div class="text-right flex-shrink-0">
            <span class="text-lg font-black text-primary">{{ entry.balance }}</span>
            <span class="text-xs font-bold text-gray-400 ml-1">pts</span>
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
const activities = ref<any[]>([]);
const availableDepartments = ref<string[]>([]);
const filterCategory = ref('');
const filterDepartment = ref('');
const filterActivity = ref('');
const showCategoryDropdown = ref(false);
const showDepartmentDropdown = ref(false);
const showActivityDropdown = ref(false);
const isLoading = ref(true);

const filteredActivitiesForDropdown = computed(() => {
  if (!filterCategory.value) return activities.value;
  return activities.value.filter(a => a.category_id === filterCategory.value);
});

const top3 = computed(() => leaderboard.value.slice(0, 3));

const getInitial = (name?: string) => {
  if (!name) return 'U';
  return name.charAt(0).toUpperCase();
};

const formatName = (name?: string) => {
  if (!name) return '';
  const words = name.trim().split(' ');
  if (words.length <= 2) return name;
  return words.slice(0, 2).join(' ') + '...';
};

const getCategoryName = (id: string) => {
  if (!id) return 'Semua Kategori';
  const cat = categories.value.find((c: any) => c.id === id);
  return cat ? cat.name : 'Semua Kategori';
};

const getActivityName = (id: string) => {
  if (!id) return 'Semua Aktivitas';
  const act = activities.value.find((c: any) => c.id === id);
  return act ? act.name : 'Semua Aktivitas';
};

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

const closeDropdowns = () => {
  showCategoryDropdown.value = false;
  showDepartmentDropdown.value = false;
  showActivityDropdown.value = false;
};

const selectCategory = (id: string) => {
  filterCategory.value = id;
  // If activity does not belong to new category, reset activity filter
  if (filterActivity.value) {
    const act = activities.value.find((a: any) => a.id === filterActivity.value);
    if (act && act.category_id !== filterCategory.value && filterCategory.value !== '') {
      filterActivity.value = '';
    }
  }
  closeDropdowns();
  fetchLeaderboard();
};

const selectActivity = (act: any) => {
  if (!act) {
    filterActivity.value = '';
  } else {
    filterActivity.value = act.id;
    // Auto sync category
    if (act.category_id) {
      filterCategory.value = act.category_id;
    }
  }
  closeDropdowns();
  fetchLeaderboard();
};

const selectDepartment = (dept: string) => {
  filterDepartment.value = dept;
  closeDropdowns();
  fetchLeaderboard();
};

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
    const data = res.data.data || [];
    // Jangan masukkan user yang poinnya 0
    leaderboard.value = data.filter((entry: any) => entry.balance > 0);
    
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

onMounted(() => {
  fetchCategories();
  fetchLeaderboard();
});
</script>
