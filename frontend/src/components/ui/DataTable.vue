<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Table Header: Title + Action Button -->
    <div class="px-6 py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-gray-100">
      <div>
        <h2 class="text-base font-bold text-gray-800">{{ title }}</h2>
        <p v-if="subtitle" class="text-xs text-gray-500 mt-0.5">{{ subtitle }}</p>
      </div>
      <div class="flex items-center gap-2">
        <!-- Search -->
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari..."
            class="pl-8 pr-3 py-1.5 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 w-44"
          />
          <svg class="absolute left-2.5 top-2 h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
          </svg>
        </div>
        <!-- Slot for action button (e.g. "Tambah") -->
        <slot name="action" />
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <!-- Row number column -->
            <th class="px-4 py-3 text-center text-xs font-medium text-gray-400 uppercase tracking-wider w-12">#</th>
            <!-- Dynamic columns -->
            <th
              v-for="col in columns"
              :key="col.key"
              class="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider select-none"
              :class="[
                col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : 'text-left',
                col.sortable !== false ? 'cursor-pointer hover:bg-gray-100 transition-colors' : ''
              ]"
              @click="col.sortable !== false && toggleSort(col.key)"
            >
              <span class="flex items-center gap-1" :class="col.align === 'right' ? 'justify-end' : col.align === 'center' ? 'justify-center' : ''">
                {{ col.label }}
                <span v-if="col.sortable !== false" class="text-gray-300">
                  <svg v-if="sortKey === col.key && sortDir === 'asc'" class="h-3.5 w-3.5 text-blue-500" fill="currentColor" viewBox="0 0 20 20"><path d="M5 10l5-5 5 5H5z"/></svg>
                  <svg v-else-if="sortKey === col.key && sortDir === 'desc'" class="h-3.5 w-3.5 text-blue-500" fill="currentColor" viewBox="0 0 20 20"><path d="M15 10l-5 5-5-5h10z"/></svg>
                  <svg v-else class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M7 10l3-3 3 3H7zm0 1l3 3 3-3H7z"/></svg>
                </span>
              </span>
            </th>
            <!-- Actions column if slot is provided -->
            <th v-if="$slots.actions" class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <!-- Empty State -->
          <tr v-if="paginatedRows.length === 0">
            <td :colspan="columns.length + 2" class="px-6 py-10 text-center text-gray-400 text-sm">
              <svg class="h-10 w-10 mx-auto mb-2 text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ emptyText || 'Tidak ada data.' }}
            </td>
          </tr>
          <!-- Data rows -->
          <tr
            v-for="(row, idx) in paginatedRows"
            :key="row.id ?? idx"
            class="hover:bg-blue-50/30 transition-colors"
          >
            <!-- Row number -->
            <td class="px-4 py-3 text-center text-xs text-gray-400 font-medium">
              {{ (currentPage - 1) * pageSize + idx + 1 }}
            </td>
            <!-- Cell slot: passes row and col to parent for custom rendering -->
            <td
              v-for="col in columns"
              :key="col.key"
              class="px-4 py-3 text-sm"
              :class="col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : 'text-left'"
            >
              <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
                {{ row[col.key] ?? '—' }}
              </slot>
            </td>
            <!-- Actions slot -->
            <td v-if="$slots.actions" class="px-4 py-3 text-right text-sm font-medium whitespace-nowrap">
              <slot name="actions" :row="row" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="totalCount > 0" class="px-4 py-3 flex flex-col sm:flex-row items-center justify-between gap-2 border-t border-gray-100 bg-gray-50">
      <p class="text-xs text-gray-500">
        Menampilkan {{ (currentPage - 1) * pageSize + 1 }}–{{ Math.min(currentPage * pageSize, totalCount) }}
        dari {{ totalCount }} data
        <span v-if="searchQuery" class="text-blue-600">(filter aktif)</span>
      </p>
      <div class="flex items-center gap-1">
        <button
          @click="currentPage = 1"
          :disabled="currentPage === 1"
          class="p-1.5 rounded text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Halaman pertama"
        >
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
        </button>
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="p-1.5 rounded text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Sebelumnya"
        >
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
        </button>
        <span class="px-3 py-1 text-xs font-medium text-gray-700 bg-white border border-gray-200 rounded">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="p-1.5 rounded text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Berikutnya"
        >
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
        </button>
        <button
          @click="currentPage = totalPages"
          :disabled="currentPage === totalPages"
          class="p-1.5 rounded text-gray-500 hover:bg-gray-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Halaman terakhir"
        >
          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7"/></svg>
        </button>
        <!-- Page size selector -->
        <select
          v-model="pageSize"
          class="ml-2 text-xs border border-gray-200 rounded px-1.5 py-1 text-gray-600 focus:outline-none focus:ring-1 focus:ring-blue-500"
        >
          <option :value="5">5 / hal</option>
          <option :value="10">10 / hal</option>
          <option :value="25">25 / hal</option>
          <option :value="50">50 / hal</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

interface Column {
  key: string;
  label: string;
  align?: 'left' | 'center' | 'right';
  sortable?: boolean;
  searchable?: boolean; // whether this field is included in search
}

const props = defineProps<{
  title: string;
  subtitle?: string;
  columns: Column[];
  rows: Record<string, any>[];
  emptyText?: string;
  serverSide?: boolean;
  totalRows?: number;
}>();

const emit = defineEmits<{
  (e: 'change', params: { page: number; limit: number; search: string; sort: string; order: 'asc' | 'desc' }): void;
}>();

// ---- Pagination ----
const currentPage = ref(1);
const pageSize = ref(10);

// ---- Search ----
const searchQuery = ref('');

// ---- Sorting ----
const sortKey = ref('');
const sortDir = ref<'asc' | 'desc'>('asc');

const toggleSort = (key: string) => {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortKey.value = key;
    sortDir.value = 'asc';
  }
  currentPage.value = 1;
  emitChange();
};

const emitChange = () => {
  if (props.serverSide) {
    emit('change', {
      page: currentPage.value,
      limit: pageSize.value,
      search: searchQuery.value,
      sort: sortKey.value,
      order: sortDir.value
    });
  }
};

// Debounce search
let searchTimeout: any;
watch(searchQuery, () => {
  currentPage.value = 1;
  if (props.serverSide) {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      emitChange();
    }, 400);
  }
});

watch(pageSize, () => { 
  currentPage.value = 1; 
  emitChange();
});

watch(() => props.rows, () => { 
  if (!props.serverSide) {
    currentPage.value = 1; 
  }
});

watch(currentPage, () => {
  emitChange();
});

// ---- Computed: filter → sort → paginate ----
const filteredRows = computed(() => {
  if (props.serverSide) return props.rows;

  let result = props.rows;

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase();
    result = result.filter(row =>
      props.columns.some(col => {
        if (col.searchable === false) return false;
        const val = row[col.key];
        return val != null && String(val).toLowerCase().includes(q);
      })
    );
  }

  if (sortKey.value) {
    result = [...result].sort((a, b) => {
      const va = a[sortKey.value] ?? '';
      const vb = b[sortKey.value] ?? '';
      if (typeof va === 'number' && typeof vb === 'number') {
        return sortDir.value === 'asc' ? va - vb : vb - va;
      }
      return sortDir.value === 'asc'
        ? String(va).localeCompare(String(vb))
        : String(vb).localeCompare(String(va));
    });
  }

  return result;
});

const totalCount = computed(() => {
  if (props.serverSide) return props.totalRows || 0;
  return filteredRows.value.length;
});

const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / pageSize.value)));

const paginatedRows = computed(() => {
  if (props.serverSide) return props.rows;

  const start = (currentPage.value - 1) * pageSize.value;
  return filteredRows.value.slice(start, start + pageSize.value);
});
</script>
