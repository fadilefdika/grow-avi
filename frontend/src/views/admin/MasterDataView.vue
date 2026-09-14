<template>
  <div class="min-h-screen bg-gray-50 flex">
    
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-1 p-8 overflow-auto h-screen">
      <div class="flex justify-between items-center mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Master Data</h1>
          <p class="text-gray-500 mt-1">Kelola Kategori dan Aktivitas GROW</p>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex space-x-4 mb-6 border-b border-gray-200">
        <button 
          @click="activeTab = 'categories'" 
          :class="['pb-2 font-medium text-sm transition-colors', activeTab === 'categories' ? 'border-b-2 border-primary text-primary' : 'text-gray-500 hover:text-gray-700']"
        >
          Kategori
        </button>
        <button 
          @click="activeTab = 'activities'" 
          :class="['pb-2 font-medium text-sm transition-colors', activeTab === 'activities' ? 'border-b-2 border-primary text-primary' : 'text-gray-500 hover:text-gray-700']"
        >
          Aktivitas
        </button>
      </div>

      <!-- Error / Success Messages -->
      <div v-if="error" class="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
        {{ error }}
      </div>
      <div v-if="successMsg" class="mb-4 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm">
        {{ successMsg }}
      </div>

      <!-- Categories Tab -->
      <div v-if="activeTab === 'categories'">
        <DataTable
          title="Daftar Kategori"
          :columns="categoryColumns"
          :rows="categories"
          :server-side="true"
          :total-rows="totalCategories"
          @change="onCategoryTableChange"
          empty-text="Tidak ada kategori."
        >
          <template #action>
            <button @click="openCategoryModal()" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors text-sm font-medium whitespace-nowrap">
              + Tambah Kategori
            </button>
          </template>
          <template #row-actions="{ row }">
            <button @click="openCategoryModal(row)" class="text-blue-600 hover:text-blue-900 mr-3">Edit</button>
            <button @click="deleteCategory(row.id)" class="text-red-600 hover:text-red-900">Hapus</button>
          </template>
        </DataTable>
      </div>

      <!-- Activities Tab -->
      <div v-if="activeTab === 'activities'">
        <DataTable
          title="Daftar Aktivitas"
          :columns="activityColumns"
          :rows="activitiesWithCategoryName"
          :server-side="true"
          :total-rows="totalActivities"
          @change="onActivityTableChange"
          empty-text="Tidak ada aktivitas."
        >
          <template #action>
            <button @click="openActivityModal()" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors text-sm font-medium whitespace-nowrap">
              + Tambah Aktivitas
            </button>
          </template>
          <template #cell-name="{ row }">
            <div class="font-medium text-gray-900">{{ row.name }}</div>
            <div v-if="row.is_custom_input" class="text-xs text-orange-600">Butuh Input Custom</div>
          </template>
          <template #cell-default_points="{ value }">
            <span class="font-bold text-gray-900">{{ value }} pts</span>
          </template>
          <template #row-actions="{ row }">
            <button @click="openActivityModal(row)" class="text-blue-600 hover:text-blue-900 mr-3">Edit</button>
            <button @click="deleteActivity(row.id)" class="text-red-600 hover:text-red-900">Hapus</button>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="showCategoryModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">{{ editingCategory ? 'Edit Kategori' : 'Tambah Kategori' }}</h3>
          <button @click="showCategoryModal = false" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <div class="p-6">
          <label class="block text-sm font-medium text-gray-700 mb-1">Nama Kategori</label>
          <input type="text" v-model="categoryForm.name" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary mb-4" />
          
          <div class="flex justify-end gap-3 mt-4">
            <button @click="showCategoryModal = false" class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors font-medium text-sm">Batal</button>
            <button @click="saveCategory" :disabled="isProcessing || !categoryForm.name" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors font-medium text-sm disabled:opacity-50">Simpan</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Activity Modal -->
    <div v-if="showActivityModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">{{ editingActivity ? 'Edit Aktivitas' : 'Tambah Aktivitas' }}</h3>
          <button @click="showActivityModal = false" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Kategori</label>
            <select v-model="activityForm.category_id" class="w-full bg-white border border-gray-300 text-gray-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary">
              <option disabled value="0">Pilih Kategori</option>
              <option v-for="cat in allCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Aktivitas</label>
            <input type="text" v-model="activityForm.name" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Poin Default</label>
            <input type="number" v-model="activityForm.default_points" class="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary" />
          </div>
          <div class="flex items-center">
            <input id="is_custom_input" type="checkbox" v-model="activityForm.is_custom_input" class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded">
            <label for="is_custom_input" class="ml-2 block text-sm text-gray-900">
              Memerlukan Input Keterangan Tambahan (Nomor SS, dsb)
            </label>
          </div>
          
          <div class="flex justify-end gap-3 mt-6">
            <button @click="showActivityModal = false" class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors font-medium text-sm">Batal</button>
            <button @click="saveActivity" :disabled="isProcessing || !activityForm.name || !activityForm.category_id" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-800 transition-colors font-medium text-sm disabled:opacity-50">Simpan</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';
import AdminSidebar from '../../components/layout/AdminSidebar.vue';
import DataTable from '../../components/ui/DataTable.vue';

const router = useRouter();
const authStore = useAuthStore();

const activeTab = ref<'categories' | 'activities'>('categories');
const allCategories = ref<any[]>([]); // For lookup & dropdown
const categories = ref<any[]>([]);
const totalCategories = ref(0);
const activities = ref<any[]>([]);
const totalActivities = ref(0);

const error = ref('');
const successMsg = ref('');
const isProcessing = ref(false);

const catParams = ref({ page: 1, limit: 10, search: '', sort: '', order: 'asc' });
const actParams = ref({ page: 1, limit: 10, search: '', sort: '', order: 'asc' });

const categoryColumns = [
  { key: 'name', label: 'Nama Kategori' }
];

const activityColumns = [
  { key: 'name', label: 'Aktivitas' },
  { key: 'category_name', label: 'Kategori' },
  { key: 'default_points', label: 'Poin Base' }
];

const activitiesWithCategoryName = computed(() =>
  activities.value.map(act => ({
    ...act,
    category_name: getCategoryName(act.category_id)
  }))
);

const showCategoryModal = ref(false);
const editingCategory = ref<any>(null);
const categoryForm = ref({ name: '' });

const showActivityModal = ref(false);
const editingActivity = ref<any>(null);
const activityForm = ref({ category_id: 0, name: '', default_points: 0, is_custom_input: false });

const fetchAllCategories = async () => {
  try {
    const res = await apiClient.get('/categories');
    allCategories.value = res.data.data || [];
  } catch (err) {
    console.error("Gagal memuat semua kategori", err);
  }
};

const fetchCategories = async () => {
  try {
    const res = await apiClient.get('/categories', { params: catParams.value });
    categories.value = res.data.data || [];
    totalCategories.value = res.data.total || 0;
  } catch (err) {
    console.error("Gagal memuat kategori", err);
  }
};

const fetchActivities = async () => {
  try {
    const res = await apiClient.get('/activities', { params: actParams.value });
    activities.value = res.data.data || [];
    totalActivities.value = res.data.total || 0;
  } catch (err) {
    console.error("Gagal memuat aktivitas", err);
  }
};

const fetchData = async () => {
  await fetchAllCategories();
  fetchCategories();
  fetchActivities();
};

const onCategoryTableChange = (params: any) => {
  catParams.value = params;
  fetchCategories();
};

const onActivityTableChange = (params: any) => {
  actParams.value = params;
  fetchActivities();
};

const getCategoryName = (id: number) => {
  const cat = allCategories.value.find(c => c.id === id);
  return cat ? cat.name : '-';
};

const showMessage = (msg: string, isError = false) => {
  if (isError) {
    error.value = msg;
    successMsg.value = '';
  } else {
    successMsg.value = msg;
    error.value = '';
  }
  setTimeout(() => {
    error.value = '';
    successMsg.value = '';
  }, 5000);
};

// --- Category Logic ---
const openCategoryModal = (cat?: any) => {
  if (cat) {
    editingCategory.value = cat;
    categoryForm.value.name = cat.name;
  } else {
    editingCategory.value = null;
    categoryForm.value.name = '';
  }
  showCategoryModal.value = true;
};

const saveCategory = async () => {
  isProcessing.value = true;
  try {
    if (editingCategory.value) {
      await apiClient.put(`/admin/categories/${editingCategory.value.id}`, categoryForm.value);
      showMessage('Kategori berhasil diupdate');
    } else {
      await apiClient.post('/admin/categories', categoryForm.value);
      showMessage('Kategori berhasil ditambahkan');
    }
    showCategoryModal.value = false;
    await fetchAllCategories();
    fetchCategories();
  } catch (err: any) {
    showMessage(err.response?.data?.error || 'Gagal menyimpan kategori', true);
  } finally {
    isProcessing.value = false;
  }
};

const deleteCategory = async (id: number) => {
  if (!confirm('Yakin ingin menghapus kategori ini?')) return;
  try {
      await apiClient.delete(`/admin/categories/${id}`);
      showMessage('Kategori berhasil dihapus');
      await fetchAllCategories();
      fetchCategories();
    } catch (err: any) {
    showMessage(err.response?.data?.error || 'Gagal menghapus kategori', true);
  }
};

// --- Activity Logic ---
const openActivityModal = (act?: any) => {
  if (act) {
    editingActivity.value = act;
    activityForm.value = {
      category_id: act.category_id,
      name: act.name,
      default_points: act.default_points || 0,
      is_custom_input: act.is_custom_input
    };
  } else {
    editingActivity.value = null;
    activityForm.value = { category_id: 0, name: '', default_points: 0, is_custom_input: false };
  }
  showActivityModal.value = true;
};

const saveActivity = async () => {
  isProcessing.value = true;
  try {
    if (editingActivity.value) {
      await apiClient.put(`/admin/activities/${editingActivity.value.id}`, activityForm.value);
      showMessage('Aktivitas berhasil diupdate');
    } else {
      await apiClient.post('/admin/activities', activityForm.value);
      showMessage('Aktivitas berhasil ditambahkan');
    }
    showActivityModal.value = false;
    fetchActivities();
  } catch (err: any) {
    showMessage(err.response?.data?.error || 'Gagal menyimpan aktivitas', true);
  } finally {
    isProcessing.value = false;
  }
};

const deleteActivity = async (id: number) => {
  if (!confirm('Yakin ingin menghapus aktivitas ini?')) return;
  try {
    await apiClient.delete(`/admin/activities/${id}`);
    showMessage('Aktivitas berhasil dihapus');
    fetchActivities();
  } catch (err: any) {
    showMessage(err.response?.data?.error || 'Gagal menghapus aktivitas', true);
  }
};

onMounted(() => {
  fetchData();
});
</script>
