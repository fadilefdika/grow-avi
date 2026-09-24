<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Hero Header -->
    <div class="hero-gradient pt-12 pb-28 px-6 relative overflow-hidden">
      <!-- Background decorative elements -->
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <!-- Awan -->
        <img
          src="/img/awan-dashboard.png"
          alt=""
          class="absolute -top-4 left-0 w-20 opacity-90 animate-float opacity-60"
        />
        <!-- Uang -->
        <img
          src="/img/uang-dashboard.png"
          alt=""
          class="absolute -top-6 right-2 w-32 opacity-95 animate-float"
          style="animation-delay: 1s;"
        />
        <!-- Bunga kiri, mentok bawah -->
        <img
          src="/img/bunga-dashboard.png"
          alt=""
          class="absolute -bottom-12 -left-4 w-32 opacity-90 z-10"
        />
        <!-- Bunga kanan, mentok bawah -->
        <img
          src="/img/bunga-dashboard.png"
          alt=""
          class="absolute -bottom-12 -right-4 w-32 opacity-90 -scale-x-100 z-10"
        />

        <!-- Gradient overlay biru → transparan, memudar ke atas -->
        <div
          class="absolute bottom-0 left-0 w-full h-20 bg-gradient-to-t from-[#0D47A1]/90 via-[#1E6FA8]/40 to-transparent z-20"
        ></div>
      </div>

      <div class="flex justify-between items-center gap-3 relative z-10 min-h-[64px]">
        <!-- Name section -->
        <div class="min-w-0 flex-1">
          <p class="text-white/70 text-sm font-semibold mb-0.5">Halo, 👋</p>
          <h1 class="text-white text-xl font-black leading-tight truncate">{{ formatName(auth.user?.name) || auth.user?.npk || 'Karyawan' }}</h1>
          <p class="text-white/60 text-xs font-medium mt-0.5">Siap kumpulkan poin hari ini?</p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0 relative z-10">
          <!-- Avatar -->
          <div class="w-10 h-10 rounded-2xl glass-card flex items-center justify-center shadow-lg">
            <span class="text-white font-black text-base">{{ getInitial(auth.user?.name || auth.user?.npk || 'U') }}</span>
          </div>
          <!-- Logout -->
          <button @click="handleLogout" class="w-9 h-9 rounded-xl glass-card flex items-center justify-center transition-all hover:bg-red-500/40 text-white">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content Overlay -->
    <div class="px-4 -mt-20 relative z-20 space-y-4">

      <!-- Points & Rank Card -->
      <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
        <div class="flex items-stretch divide-x divide-gray-100">
          <!-- Points -->
          <div class="flex-1 p-5">
            <p class="text-gray-500 text-xs font-bold uppercase tracking-wider mb-2">Total Poin GROW</p>
            <div class="flex items-end gap-1.5">
              <span class="text-4xl font-black text-gray-900 leading-none tabular-nums">{{ displayBalance }}</span>
              <span class="text-sm font-bold text-secondary mb-1">pts</span>
            </div>
            <!-- XP Progress Bar -->
            <div class="mt-3">
              <div class="xp-bar-track">
                <div class="xp-bar-fill" :style="{ width: xpPercent + '%' }"></div>
              </div>
              <p class="text-[10px] text-gray-400 font-semibold mt-1">{{ balance }} / {{ nextMilestone }} pts ke milestone berikutnya</p>
            </div>
          </div>
          <!-- Rank -->
          <div class="p-5 text-right cursor-pointer group" @click="$router.push('/leaderboard')">
            <p class="text-gray-500 text-[9px] font-bold uppercase tracking-wider mb-2">
              Peringkat
              <span class="ml-1 text-sm font-bold bg-primary/10 text-primary px-1.5 py-0.5 rounded-full">Lihat</span>
            </p>
            <div class="flex items-center justify-end gap-2">
              <span v-if="balance > 0" class="text-3xl font-black text-gray-900 group-hover:text-primary transition-colors">#{{ rank }}</span>
              <span v-else class="text-lg font-black text-gray-400 group-hover:text-primary transition-colors">Belum ada</span>
            </div>
            <div class="flex items-center justify-end gap-1 mt-1">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-secondary" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
              </svg>
              <span class="text-[10px] font-bold text-gray-400 group-hover:text-primary transition-colors">Leaderboard</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-2 gap-3">
        <!-- Submit -->
        <router-link
          to="/submission"
          id="btn-buat-pengajuan"
          class="card-game p-5 flex flex-col items-center justify-center gap-3 group cursor-pointer border-2 border-transparent hover:border-primary/30"
        >
          <div class="w-14 h-14 rounded-2xl bg-blue-50 flex items-center justify-center group-hover:bg-primary transition-all duration-300 shadow-sm group-hover:shadow-primary/30 group-hover:shadow-md">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-primary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div class="text-center">
            <span class="text-sm font-black text-gray-800 group-hover:text-primary transition-colors">Buat Pengajuan</span>
            <p class="text-[10px] text-gray-400 font-semibold mt-0.5">Kumpulkan bukti aktivitas</p>
          </div>
        </router-link>

        <!-- Rewards -->
        <router-link
          to="/rewards"
          id="btn-tukar-hadiah"
          class="card-game p-5 flex flex-col items-center justify-center gap-3 group cursor-pointer border-2 border-transparent hover:border-secondary/30"
        >
          <div class="w-14 h-14 rounded-2xl bg-orange-50 flex items-center justify-center group-hover:bg-secondary transition-all duration-300 shadow-sm group-hover:shadow-secondary/30 group-hover:shadow-md">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-secondary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
          </div>
          <div class="text-center">
            <span class="text-sm font-black text-gray-800 group-hover:text-secondary transition-colors">Tukar Hadiah</span>
            <p class="text-[10px] text-gray-400 font-semibold mt-0.5">Gunakan poinmu</p>
          </div>
        </router-link>
      </div>

      <!-- Admin Mode Banner -->
      <div v-if="auth.user?.role === 'admin'" class="rounded-2xl p-4 flex items-center justify-between" style="background: linear-gradient(135deg, #1e293b, #334155);">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-yellow-500/20 rounded-xl flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <div>
            <h3 class="text-white font-bold text-sm">Mode Administrator</h3>
            <p class="text-gray-400 text-xs font-medium">Akses panel admin penuh</p>
          </div>
        </div>
        <router-link to="/admin/dashboard" class="px-4 py-2 bg-yellow-400 text-gray-900 font-black rounded-xl text-sm hover:bg-yellow-300 transition-colors shadow-lg">
          Buka Panel
        </router-link>
      </div>

      <!-- Recent Submissions -->
      <div class="card-game p-5">
        <div class="flex justify-between items-center mb-4">
          <div>
            <h3 class="font-black text-gray-800 text-base">Riwayat Terkini</h3>
            <p class="text-xs text-gray-400 font-semibold">5 pengajuan terakhir</p>
          </div>
          <button
            @click="$router.push('/history')"
            v-if="allSubmissions.length > 0"
            class="text-xs font-bold text-primary bg-primary/10 hover:bg-primary/20 px-3 py-1.5 rounded-xl transition-colors"
          >
            Lihat Semua
          </button>
        </div>

        <div v-if="isLoading" class="flex justify-center py-8">
          <svg class="animate-spin h-7 w-7 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <div v-else-if="allSubmissions.length === 0" class="text-center py-8">
          <div class="w-16 h-16 bg-gray-50 rounded-2xl flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-sm font-bold text-gray-400">Belum ada pengajuan.</p>
          <p class="text-xs text-gray-300 font-medium mt-1">Yuk, mulai kumpulkan poinmu!</p>
        </div>

        <div v-else class="space-y-2.5">
          <div
            v-for="sub in displayedSubmissions"
            :key="sub.id"
            @click="openDetail(sub)"
            class="relative p-3 pl-4 rounded-xl border border-gray-100 bg-white cursor-pointer hover:bg-gray-50 hover:border-primary/20 transition-all duration-200 flex flex-col gap-1 overflow-hidden"
          >
            <!-- Accent Line -->
            <div class="absolute left-0 top-0 bottom-0 w-1" :class="{'bg-amber-500': sub.status === 'PENDING', 'bg-green-500': sub.status === 'APPROVED', 'bg-red-500': sub.status === 'REJECTED', 'bg-yellow-500': sub.status === 'REJECTED_RESUBMIT'}"></div>

            <div class="flex items-start justify-between gap-2">
              <div class="overflow-hidden">
                <p class="text-sm font-black text-gray-800 truncate mb-0.5">{{ sub.category_name || 'Memuat...' }}</p>
                <p class="text-[10px] font-bold text-gray-500 truncate">{{ sub.grow_id || 'Pengajuan Baru' }} &middot; Dikirim {{ formatDate(sub.created_at) }}</p>
              </div>
              <div class="shrink-0">
                <span :class="getStatusBadgeClass(sub.status)" class="text-[9px] px-2 py-0.5 rounded font-black whitespace-nowrap">
                  {{ formatStatus(sub.status) }}
                </span>
              </div>
            </div>
            
            <!-- Alasan Reject -->
            <div v-if="(sub.status === 'REJECTED' || sub.status === 'REJECTED_RESUBMIT') && sub.admin_notes && (sub.admin_notes.Valid ? sub.admin_notes.String : sub.admin_notes)" class="mt-1 animate-fade-in">
              <p class="text-[11px] font-semibold text-gray-600 truncate"><span class="font-bold text-red-500">✕ Alasan:</span> {{ sub.admin_notes.Valid ? sub.admin_notes.String : sub.admin_notes }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4 bg-gray-900/60 backdrop-blur-sm">
      <div class="bg-white rounded-t-3xl sm:rounded-3xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col overflow-hidden animate-slide-up">
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center shrink-0">
          <h3 class="text-lg font-black text-gray-900">Detail Pengajuan</h3>
          <button @click="closeDetailModal" class="w-8 h-8 rounded-xl bg-gray-100 hover:bg-gray-200 flex items-center justify-center transition-colors text-gray-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="p-6 overflow-y-auto">
          <div v-if="isLoadingDetail" class="flex justify-center py-12">
            <svg class="animate-spin h-8 w-8 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          
          <div v-else-if="selectedDetail" class="space-y-5">
            <div class="flex items-start justify-between">
              <div>
                <p class="text-xs text-gray-400 font-bold uppercase tracking-wider">GROW ID</p>
                <p class="text-xl font-black text-gray-900 mt-0.5">{{ selectedDetail.grow_id }}</p>
              </div>
              <span :class="getStatusBadgeClass(selectedDetail.status)" class="text-xs font-black px-3 py-1.5 rounded-xl whitespace-nowrap">
                {{ formatStatus(selectedDetail.status) }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100" v-if="selectedDetail.category_name !== 'INNOVATION'">
                <p class="text-[10px] text-gray-500 font-semibold mb-1">Aktivitas</p>
                <p class="text-[11px] font-bold text-gray-800 leading-tight break-words">{{ (selectedDetail.activity_name === 'Yang lain' && selectedDetail.custom_activity_type?.Valid) ? selectedDetail.custom_activity_type.String : (selectedDetail.activity_name || '-') }}</p>
              </div>
              <div class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100">
                <p class="text-[10px] text-gray-500 font-semibold mb-1">Kategori</p>
                <p class="text-[11px] font-bold text-gray-800 leading-tight break-words">{{ selectedDetail.category_name || '-' }}</p>
              </div>
              <div class="bg-blue-50 p-3.5 rounded-2xl border border-blue-100">
                <p class="text-[10px] text-primary/70 font-semibold mb-1">Poin</p>
                <p class="text-base font-bold text-primary">{{ selectedDetail.points_awarded || '-' }} <span class="text-[10px] font-semibold">pts</span></p>
              </div>
              <div class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100">
                <p class="text-[10px] text-gray-500 font-semibold mb-1">Tgl Kegiatan</p>
                <p class="text-[11px] font-bold text-gray-800 leading-tight">{{ formatDate(selectedDetail.activity_date) }}</p>
              </div>
              <div class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100">
                <p class="text-[10px] text-gray-500 font-semibold mb-1">Waktu Submit</p>
                <p class="text-[11px] font-bold text-gray-800 leading-tight">{{ formatDate(selectedDetail.created_at) }} WIB</p>
              </div>
            </div>

            <div v-if="selectedDetail.custom_reference?.Valid" class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100">
              <p class="text-xs text-gray-400 font-bold mb-1">Nama Kegiatan/Aktivitas</p>
              <p class="text-sm font-semibold text-gray-900">{{ selectedDetail.custom_reference.String }}</p>
            </div>

            <div v-if="selectedDetail.nomor_ss?.Valid" class="bg-gray-50 p-3.5 rounded-2xl border border-gray-100">
              <p class="text-xs text-gray-400 font-bold mb-1">Nomor SS</p>
              <p class="text-sm font-black text-gray-900">{{ selectedDetail.nomor_ss.String }}</p>
            </div>

            <div v-if="(selectedDetail.status === 'REJECTED' || selectedDetail.status === 'REJECTED_RESUBMIT') && selectedDetail.admin_notes?.Valid" class="bg-red-50 p-4 rounded-2xl border border-red-200">
              <div class="flex justify-between items-start mb-2">
                <p class="text-xs font-black text-red-600 uppercase flex items-center gap-1.5"><span class="text-[10px]">❌</span> Alasan Penolakan</p>
                <span v-if="selectedDetail.status === 'REJECTED_RESUBMIT'" class="px-2 py-1 bg-yellow-100 text-yellow-700 border border-yellow-200 text-[10px] font-black rounded-md">📝 BISA DIREVISI</span>
                <span v-else class="px-2 py-1 bg-red-100 text-red-700 border border-red-200 text-[10px] font-black rounded-md">⛔ FINAL (TIDAK BISA DIREVISI)</span>
              </div>
              <p class="text-sm text-red-800 font-semibold break-all">{{ selectedDetail.admin_notes.String }}</p>
            </div>

            <!-- Evidence Section -->
            <div v-if="selectedDetail.evidence && selectedDetail.evidence.length > 0">
              <p class="text-sm font-black text-gray-800 mb-3">Bukti Lampiran</p>
              <div class="grid grid-cols-2 gap-3">
                <a v-for="(ev, idx) in selectedDetail.evidence" :key="idx" :href="ev" target="_blank" class="block group relative rounded-2xl overflow-hidden border border-gray-200 bg-gray-50">
                  <div class="aspect-video w-full">
                    <img :src="ev" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                  </div>
                  <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                    <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                  </div>
                </a>
              </div>
            </div>
          </div>
        </div>
        
        <div class="p-4 border-t border-gray-100 bg-gray-50/50 shrink-0 flex gap-3">
          <button v-if="selectedDetail?.status === 'REJECTED_RESUBMIT'" @click="goToRevision" class="btn-game-primary flex-1 py-3 text-sm">
            Revisi Pengajuan
          </button>
          <button @click="closeDetailModal" class="flex-1 py-3 bg-white border-2 border-gray-200 text-gray-700 font-bold rounded-2xl hover:bg-gray-50 transition-colors text-sm">
            Tutup
          </button>
        </div>
      </div>
    </div>



  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import apiClient from '../api/client';

const authStore = useAuthStore();
const router = useRouter();

const auth = useAuthStore();
const balance = ref(0); 
const rank = ref('-');
const allSubmissions = ref<any[]>([]);
const isLoading = ref(true);

// Animated balance counter
const displayBalance = ref(0);

const nextMilestone = computed(() => {
  const milestones = [50, 100, 200, 300, 500, 750, 1000, 1500, 2000];
  return milestones.find(m => m > balance.value) || balance.value + 100;
});

const xpPercent = computed(() => {
  const prev = (() => {
    const milestones = [0, 50, 100, 200, 300, 500, 750, 1000, 1500, 2000];
    const idx = milestones.findIndex(m => m >= nextMilestone.value);
    return idx > 0 ? milestones[idx - 1] : 0;
  })();
  const range = nextMilestone.value - prev;
  const progress = balance.value - prev;
  return Math.min(100, Math.round((progress / range) * 100));
});

const displayedSubmissions = computed(() => {
  const rejectedResubmit = allSubmissions.value.filter(sub => sub.status === 'REJECTED_RESUBMIT');
  const others = allSubmissions.value.filter(sub => sub.status !== 'REJECTED_RESUBMIT');
  return [...rejectedResubmit, ...others].slice(0, 5);
});

const showDetailModal = ref(false);
const selectedDetail = ref<any>(null);
const isLoadingDetail = ref(false);

const openDetail = async (sub: any) => {
  try {
    isLoadingDetail.value = true;
    showDetailModal.value = true;
    selectedDetail.value = null;
    const res = await apiClient.get(`/submissions/${sub.id}`);
    selectedDetail.value = res.data.data;
  } catch (err) {
    console.error("Gagal mengambil detail", err);
    showDetailModal.value = false;
  } finally {
    isLoadingDetail.value = false;
  }
};

const closeDetailModal = () => {
  showDetailModal.value = false;
  selectedDetail.value = null;
};

const goToRevision = () => {
  if (selectedDetail.value) {
    router.push(`/submission?revisi_id=${selectedDetail.value.id}`);
  }
};

const getInitial = (name: string) => {
  if (!name) return 'U';
  return name.charAt(0).toUpperCase();
};

const formatName = (name?: string) => {
  if (!name) return '';
  const words = name.trim().split(' ');
  if (words.length <= 2) return name;
  return words.slice(0, 2).join(' ') + '...';
};

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Jakarta' }).replace(/\./g, ':');
};

const getStatusIconBg = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'bg-green-100';
    case 'REJECTED':
    case 'REJECTED_RESUBMIT': return 'bg-red-100';
    default: return 'bg-amber-100';
  }
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'APPROVED': return 'badge-approved';
    case 'REJECTED': return 'badge-rejected';
    case 'REJECTED_RESUBMIT': return 'bg-amber-100 text-amber-700 border border-amber-200';
    default: return 'badge-pending';
  }
};

const formatStatus = (status: string) => {
  if (status === 'REJECTED_RESUBMIT') return 'REJECTED · REVISI';
  if (status === 'REJECTED') return 'REJECTED · FINAL';
  return status;
};

const animateBalance = (target: number) => {
  const duration = 1000;
  const step = 16;
  const increment = target / (duration / step);
  let current = 0;
  const timer = setInterval(() => {
    current += increment;
    if (current >= target) {
      displayBalance.value = target;
      clearInterval(timer);
    } else {
      displayBalance.value = Math.floor(current);
    }
  }, step);
};

const fetchDashboardData = async () => {
  isLoading.value = true;
  try {
    // Try to fetch leaderboard balance
    try {
      const balRes = await apiClient.get('/leaderboard/my-balance');
      balance.value = balRes.data.balance || 0;
      rank.value = balRes.data.rank || '-';
      animateBalance(balance.value);
    } catch (err) {
      console.warn("Could not fetch balance", err);
    }

    // Fetch recent submissions
    try {
      const subRes = await apiClient.get('/submissions/me');
      allSubmissions.value = subRes.data.submissions || [];
    } catch (err) {
      console.warn("Could not fetch submissions", err);
    }
    
  } finally {
    isLoading.value = false;
  }
};

const handleLogout = async () => {
  try {
    await apiClient.post('/auth/logout');
  } catch (err) {
    console.error('Logout error:', err);
  } finally {
    auth.logout();
    router.push('/login');
  }
};

onMounted(() => {
  fetchDashboardData();
});
</script>
