<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Hero (compact) -->
    <header class="hero-gradient pt-10 pb-12 px-4 relative overflow-hidden">
      <div class="absolute inset-0 overflow-hidden pointer-events-none" aria-hidden="true">
        <div class="absolute -top-8 -right-8 w-32 h-32 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFA64D, transparent);"></div>
      </div>

      <div class="relative z-10 flex items-center gap-3">
        <button
          @click="$router.push('/dashboard')"
          aria-label="Kembali ke dashboard"
          class="shrink-0 w-9 h-9 rounded-xl bg-white/15 text-white flex items-center justify-center hover:bg-white/25 transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-white"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>

        <div class="flex-1 min-w-0">
          <h1 class="text-white text-lg font-black leading-tight truncate">🎁 Katalog Hadiah</h1>
          <p class="text-white/80 text-[11px] font-semibold truncate">Tukarkan poinmu dengan hadiah menarik!</p>
        </div>

        <div class="glass-card rounded-xl px-3 py-1.5 text-right shrink-0">
          <p class="text-white/80 text-[10px] font-bold leading-none">Saldo poin</p>
          <p class="text-white text-lg font-black leading-tight">
            {{ fmt(balance) }}<span class="text-white/80 text-xs font-bold ml-1">pts</span>
          </p>
        </div>
      </div>
    </header>

    <!-- Main content -->
    <main class="px-4 -mt-6 relative z-20">
      <!-- Loading: skeleton rows -->
      <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 gap-2.5" aria-busy="true" aria-label="Memuat hadiah">
        <div v-for="n in 5" :key="n" class="flex items-center gap-3 bg-white rounded-2xl p-3 border border-gray-100 animate-pulse">
          <div class="w-10 h-10 rounded-xl bg-gray-100"></div>
          <div class="flex-1 space-y-2">
            <div class="h-3 w-2/3 rounded bg-gray-100"></div>
            <div class="h-3 w-1/3 rounded bg-gray-100"></div>
          </div>
          <div class="w-16 h-9 rounded-xl bg-gray-100"></div>
        </div>
      </div>

      <!-- Error state -->
      <div v-else-if="loadError" class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm text-center" role="alert">
        <p class="text-3xl mb-2">😵</p>
        <p class="font-black text-gray-800 text-sm">Katalog gagal dimuat</p>
        <p class="text-gray-500 text-xs font-semibold mt-1">Periksa koneksimu, lalu coba lagi.</p>
        <button @click="fetchData()" class="btn-game-orange mt-4 px-5 py-2.5 text-sm">Muat ulang</button>
      </div>

      <!-- Empty state -->
      <div v-else-if="rewards.length === 0" class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm text-center">
        <p class="text-3xl mb-2">🎪</p>
        <p class="font-black text-gray-800 text-sm">Katalog hadiah masih kosong</p>
        <p class="text-gray-500 text-xs font-semibold mt-1">Hadiah baru akan muncul di sini.</p>
      </div>

      <!-- Reward list -->
      <ul v-else class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
        <li
          v-for="reward in sortedRewards"
          :key="reward.id"
          :class="[
            'flex items-center gap-3 bg-white rounded-2xl p-3 transition-shadow',
            reward.stock <= 0
              ? 'border border-gray-100 opacity-60'
              : isEligible(reward)
                ? 'card-eligible shadow-md'
                : 'border border-gray-100 shadow-sm'
          ]"
        >
          <!-- Icon -->
          <div
            :class="[
              'shrink-0 w-10 h-10 rounded-xl bg-gradient-to-br from-orange-100 to-orange-50 border border-orange-100 flex items-center justify-center',
              reward.stock <= 0 ? 'grayscale' : ''
            ]"
            aria-hidden="true"
          >
            <span class="text-xl">{{ getRewardEmoji(reward.title) }}</span>
          </div>

          <!-- Info -->
          <div class="flex-1 min-w-0">
            <h3 class="font-extrabold text-gray-800 text-sm leading-tight break-words">{{ reward.title }}</h3>

            <div class="flex items-center gap-1.5 mt-0.5">
              <span class="text-base font-black text-primary leading-none">{{ fmt(reward.points_required) }}</span>
              <span class="text-[11px] font-bold text-gray-500">pts</span>
              <span
                v-if="reward.stock > 0 && reward.stock <= 5"
                class="text-[10px] font-bold text-orange-600 bg-orange-50 rounded-md px-1.5 py-0.5"
              >
                Sisa {{ reward.stock }}
              </span>
            </div>

            <!-- Progress toward reward -->
            <div v-if="reward.stock > 0 && !isEligible(reward)" class="mt-1.5 flex items-center gap-2">
              <div
                class="flex-1 h-1 bg-gray-100 rounded-full overflow-hidden"
                role="progressbar"
                :aria-valuenow="Math.round(progressPct(reward))"
                aria-valuemin="0"
                aria-valuemax="100"
                :aria-label="`Progres poin untuk ${reward.title}`"
              >
                <div class="h-full bg-orange-400 rounded-full" :style="{ width: progressPct(reward) + '%' }"></div>
              </div>
              <span class="text-[10px] font-bold text-gray-500 whitespace-nowrap">
                Kurang {{ fmt(reward.points_required - balance) }}
              </span>
            </div>
          </div>

          <!-- Action -->
          <span
            v-if="reward.stock <= 0"
            class="shrink-0 px-3 py-2 rounded-xl bg-gray-100 text-gray-500 font-black text-xs"
          >
            Habis
          </span>
          <button
            v-else
            @click="confirmRedeem(reward, $event)"
            :disabled="!isEligible(reward) || isRedeeming"
            :aria-label="`Tukar ${reward.title}`"
            :class="[
              'shrink-0 px-4 py-2 rounded-xl font-black text-sm transition-all duration-200',
              isEligible(reward) ? 'btn-game-orange text-white' : 'bg-gray-100 text-gray-400 cursor-not-allowed'
            ]"
          >
            Tukar
          </button>
        </li>
      </ul>
    </main>

    <!-- Confirm modal -->
    <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-150 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-100 ease-in"
      leave-to-class="opacity-0"
    >
    <div
      v-if="selectedReward"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm"
      @click.self="closeModal"
    >
      <div
        ref="dialogRef"
        role="dialog"
        aria-modal="true"
        aria-labelledby="redeem-title"
        tabindex="-1"
        class="bg-white rounded-3xl shadow-2xl w-full max-w-sm max-h-[90vh] max-h-[90dvh] flex flex-col overflow-hidden outline-none"
        @keydown="onDialogKeydown"
      >
        <!-- Header -->
        <div class="flex items-center gap-3 p-4 border-b border-gray-100">
          <div class="shrink-0 w-10 h-10 rounded-xl bg-orange-50 border border-orange-100 flex items-center justify-center" aria-hidden="true">
            <span class="text-xl">{{ getRewardEmoji(selectedReward.title) }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <h2 id="redeem-title" class="font-black text-gray-900 text-base leading-tight break-words">
              Tukar {{ selectedReward.title }}
            </h2>
            <p class="text-xs font-semibold text-gray-500">{{ fmt(selectedReward.points_required) }} pts</p>
          </div>
          <button
            @click="closeModal"
            aria-label="Tutup"
            class="shrink-0 w-8 h-8 rounded-lg text-gray-400 hover:bg-gray-100 hover:text-gray-600 flex items-center justify-center transition-colors"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Ticket list (scrolls) -->
        <div class="flex-1 min-h-0 overflow-y-auto p-4 space-y-2 custom-scrollbar">
          <div class="flex items-center justify-between">
            <p class="text-xs font-bold text-gray-600">Tiket poin (ID GROW)</p>
            <button
              v-if="sortedPoints.length > 0"
              @click="autoSelect(selectedReward.points_required)"
              class="text-xs font-bold text-primary hover:underline"
            >
              Pilih otomatis
            </button>
          </div>

          <label
            v-for="pt in sortedPoints"
            :key="pt.id"
            class="flex items-center justify-between gap-3 p-2.5 rounded-xl border cursor-pointer transition-colors"
            :class="selectedSubmissionIds.has(pt.id) ? 'bg-blue-50 border-primary' : 'bg-white border-gray-200 hover:border-gray-300'"
          >
            <div class="flex items-center gap-3 min-w-0">
              <input
                type="checkbox"
                :checked="selectedSubmissionIds.has(pt.id)"
                @change="togglePoint(pt.id)"
                class="w-4 h-4 shrink-0 text-primary bg-gray-100 border-gray-300 rounded focus:ring-primary"
              />
              <div class="min-w-0">
                <p class="text-xs font-bold" :class="selectedSubmissionIds.has(pt.id) ? 'text-primary' : 'text-gray-900'">{{ pt.grow_id }}</p>
                <p class="text-[11px] font-semibold text-gray-500 truncate">{{ pt.activity_name }}</p>
              </div>
            </div>
            <span class="shrink-0 text-sm font-black" :class="selectedSubmissionIds.has(pt.id) ? 'text-primary' : 'text-gray-900'">
              {{ fmt(pt.points_awarded) }} pts
            </span>
          </label>

          <div v-if="sortedPoints.length === 0" class="text-center py-4 bg-gray-50 rounded-xl">
            <p class="text-xs font-semibold text-gray-500">Tidak ada tiket poin yang tersedia.</p>
          </div>
        </div>

        <!-- Footer (always visible) -->
        <div class="p-4 pt-3 border-t border-gray-100 bg-white space-y-3">
          <div>
            <div class="flex justify-between items-center">
              <span class="text-xs font-bold text-gray-600">Total terpilih</span>
              <span
                class="text-sm font-black"
                :class="totalSelectedPoints >= selectedReward.points_required ? 'text-green-600' : 'text-red-500'"
              >
                {{ fmt(totalSelectedPoints) }} / {{ fmt(selectedReward.points_required) }} pts
              </span>
            </div>
            <p v-if="shortfall > 0" class="text-[11px] font-semibold text-red-500 mt-1" role="status">
              Kurang {{ fmt(shortfall) }} pts. Pilih tiket lagi.
            </p>
            <p v-else-if="excess > 0" class="text-[11px] font-semibold text-amber-600 mt-1" role="status">
              Kelebihan {{ fmt(excess) }} pts. Kurangi tiket jika ingin menyisakan poin.
            </p>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <button
              @click="closeModal"
              :disabled="isRedeeming"
              class="py-3 bg-gray-100 text-gray-700 font-black rounded-2xl hover:bg-gray-200 transition-colors text-sm disabled:opacity-60"
            >
              Batal
            </button>
            <button
              @click="executeRedeem"
              :disabled="isRedeeming || totalSelectedPoints < selectedReward.points_required"
              class="btn-game-orange py-3 text-sm flex items-center justify-center gap-2 disabled:opacity-60"
            >
              <svg v-if="isRedeeming" class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" aria-hidden="true">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isRedeeming ? 'Memproses...' : 'Tukar sekarang' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue';
import apiClient from '../../api/client';
import { useToast } from 'vue-toastification';

const rewards = ref<any[]>([]);
const balance = ref(0);
const isLoading = ref(true);
const loadError = ref(false);
const isRedeeming = ref(false);
const selectedReward = ref<any>(null);
const availablePoints = ref<any[]>([]);
const selectedSubmissionIds = ref<Set<number>>(new Set());
const dialogRef = ref<HTMLElement | null>(null);
const toast = useToast();

let lastFocusedEl: HTMLElement | null = null;

/* ---------- Helpers ---------- */

const fmt = (n: number): string => Number(n || 0).toLocaleString('id-ID');

const getRewardEmoji = (title: string): string => {
  const t = title?.toLowerCase() || '';
  if (t.includes('grab food') || t.includes('grabfood')) return '🍔';
  if (t.includes('grab')) return '🚗';
  if (t.includes('gojek') || t.includes('go-jek')) return '🛵';
  if (t.includes('shopee')) return '🛍️';
  if (t.includes('tokopedia')) return '🟢';
  if (t.includes('voucher') || t.includes('gift')) return '🎫';
  if (t.includes('makanan') || t.includes('makan') || t.includes('food')) return '🍱';
  if (t.includes('belanja') || t.includes('shopping')) return '🛒';
  if (t.includes('pulsa') || t.includes('telpon')) return '📱';
  return '🎁';
};

const isEligible = (r: any): boolean => r.stock > 0 && balance.value >= r.points_required;

const progressPct = (r: any): number =>
  r.points_required > 0 ? Math.min(100, (balance.value / r.points_required) * 100) : 0;

/* ---------- Derived data ---------- */

// Bisa ditukar dulu, lalu yang paling dekat tercapai, habis di paling bawah.
const sortedRewards = computed(() => {
  const rank = (r: any) => (r.stock <= 0 ? 2 : isEligible(r) ? 0 : 1);
  return [...rewards.value].sort(
    (a, b) => rank(a) - rank(b) || a.points_required - b.points_required
  );
});

// Tiket terlama dulu (FIFO), tapi prioritaskan tiket "Sisa Poin Penukaran"
const sortedPoints = computed(() => {
  const isChange = (p: any) => p.activity_name === 'Sisa Poin Penukaran' ? 0 : 1;
  const ts = (p: any) => Date.parse(p.created_at ?? p.approved_at ?? '') || 0;
  return [...availablePoints.value].sort((a, b) => isChange(a) - isChange(b) || ts(a) - ts(b) || a.id - b.id);
});

const totalSelectedPoints = computed(() =>
  availablePoints.value
    .filter(p => selectedSubmissionIds.value.has(p.id))
    .reduce((sum, p) => sum + p.points_awarded, 0)
);

const shortfall = computed(() =>
  selectedReward.value ? Math.max(0, selectedReward.value.points_required - totalSelectedPoints.value) : 0
);
const excess = computed(() =>
  selectedReward.value ? Math.max(0, totalSelectedPoints.value - selectedReward.value.points_required) : 0
);

/* ---------- Point selection ---------- */

const togglePoint = (id: number) => {
  if (selectedSubmissionIds.value.has(id)) {
    selectedSubmissionIds.value.delete(id);
  } else {
    selectedSubmissionIds.value.add(id);
  }
};

const autoSelect = (required: number) => {
  const next = new Set<number>();
  let sum = 0;
  for (const p of sortedPoints.value) {
    if (sum >= required) break;
    next.add(p.id);
    sum += p.points_awarded;
  }
  selectedSubmissionIds.value = next;
};

/* ---------- Data ---------- */

const fetchData = async (silent = false) => {
  if (!silent) isLoading.value = true;
  loadError.value = false;
  try {
    const [balRes, rwRes, apRes] = await Promise.all([
      apiClient.get('/leaderboard/my-balance'),
      apiClient.get('/rewards'),
      apiClient.get('/available-points')
    ]);
    balance.value = balRes.data.balance || 0;
    rewards.value = rwRes.data.data || [];
    availablePoints.value = apRes.data.data || [];
  } catch (err) {
    if (!silent) loadError.value = true;
    toast.error('Gagal memuat data katalog hadiah.');
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

/* ---------- Modal ---------- */

const confirmRedeem = (reward: any, event?: Event) => {
  lastFocusedEl = (event?.currentTarget as HTMLElement) ?? null;
  selectedReward.value = reward;
  autoSelect(reward.points_required);
};

const closeModal = () => {
  if (isRedeeming.value) return;
  selectedReward.value = null;
};

const executeRedeem = async () => {
  if (!selectedReward.value) return;

  isRedeeming.value = true;

  try {
    const res = await apiClient.post(`/rewards/redeem/${selectedReward.value.id}`, {
      submission_ids: Array.from(selectedSubmissionIds.value)
    });
    toast.success(res.data.message || `Berhasil menukar ${selectedReward.value.title}! 🎉`);

    // Refresh stok dan saldo tanpa memunculkan skeleton lagi
    fetchData(true);
  } catch (err: any) {
    if (err.response && err.response.data && err.response.data.error) {
      toast.error(err.response.data.error);
    } else {
      toast.error('Gagal melakukan penukaran. Coba lagi.');
    }
  } finally {
    isRedeeming.value = false;
    selectedReward.value = null;
  }
};

// Focus trap sederhana di dalam dialog
const onDialogKeydown = (e: KeyboardEvent) => {
  if (e.key !== 'Tab' || !dialogRef.value) return;
  const focusables = dialogRef.value.querySelectorAll<HTMLElement>(
    'button:not([disabled]), input:not([disabled]), [href], [tabindex]:not([tabindex="-1"])'
  );
  if (focusables.length === 0) return;
  const first = focusables[0];
  const last = focusables[focusables.length - 1];
  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault();
    last.focus();
  } else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault();
    first.focus();
  }
};

const onGlobalKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && selectedReward.value) closeModal();
};

// Scroll lock + fokus masuk/keluar dialog
watch(selectedReward, async (val) => {
  if (val) {
    document.body.style.overflow = 'hidden';
    await nextTick();
    dialogRef.value?.focus();
  } else {
    document.body.style.overflow = '';
    lastFocusedEl?.focus?.();
    lastFocusedEl = null;
  }
});

onMounted(() => {
  fetchData();
  document.addEventListener('keydown', onGlobalKeydown);
});

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onGlobalKeydown);
  document.body.style.overflow = '';
});
</script>