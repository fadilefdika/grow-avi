<template>
  <div class="min-h-screen bg-gray-50 pb-12">
    <!-- Hero Header (menggantikan header image statis) -->
    <div class="hero-gradient relative overflow-hidden">
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-10 -right-10 w-48 h-48 rounded-full opacity-10 animate-float" style="background: radial-gradient(circle, #FFA64D, transparent);"></div>
        <div class="absolute bottom-0 left-1/4 w-32 h-32 rounded-full opacity-10" style="background: radial-gradient(circle, #1CB0F6, transparent);"></div>
      </div>
      <div class="px-5 py-10 relative z-10">
        <button @click="$router.push('/dashboard')" class="flex items-center gap-1.5 text-white/70 text-sm font-bold mb-4 hover:text-white transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Kembali
        </button>
        <div class="flex items-center gap-3 mb-2">
          <div class="w-12 h-12 rounded-2xl glass-card flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div>
            <h1 class="text-white text-xl font-black leading-tight">GROW Evidence Form</h1>
            <p class="text-white/60 text-xs font-semibold">Kumpulkan bukti, raih poin!</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-md mx-auto px-4 mt-6">
      <!-- Intro Card (redesign) -->
      <div class="card-game p-5 mb-5 border-l-4 border-primary">
        <div class="flex items-start gap-3 mb-4">
          <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p class="text-sm font-black text-gray-800">Siapkan bukti terbaikmu! 🚀</p>
            <p class="text-xs text-gray-500 font-semibold mt-0.5 leading-relaxed">Isi form ini dengan lengkap dan unggah bukti aktivitasmu. Semakin lengkap, semakin besar peluang poinmu disetujui!</p>
          </div>
        </div>
        <div class="flex items-center gap-2 bg-primary/5 rounded-xl px-3 py-2.5 border border-primary/10">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span class="text-xs font-bold text-primary">Periode: 1 Januari 2026 — 31 Desember 2026</span>
        </div>
        <p class="text-primary font-black mt-3 text-base">#let'sGROW</p>
      </div>

      <!-- Action Buttons -->
      <div class="grid grid-cols-2 gap-3 mb-5">
        <a
          href="https://heyzine.com/flip-book/a0c4abfab8.html#page/1"
          target="_blank"
          rel="noopener noreferrer"
          class="card-game p-4 flex flex-col items-center justify-center gap-2 group hover:border-primary border-2 border-transparent cursor-pointer"
        >
          <div class="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center group-hover:bg-primary transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-primary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
          <span class="text-xs font-black text-gray-700 group-hover:text-primary transition-colors">📖 Guide Book</span>
        </a>
        <a
          href="https://datastudio.google.com/embed/u/0/reporting/e32866a0-4028-403f-87aa-d5e2022b9931/page/T2DvF"
          target="_blank"
          rel="noopener noreferrer"
          class="card-game p-4 flex flex-col items-center justify-center gap-2 group hover:border-secondary border-2 border-transparent cursor-pointer"
        >
          <div class="w-10 h-10 rounded-xl bg-orange-50 flex items-center justify-center group-hover:bg-secondary transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-secondary group-hover:text-white transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
          <span class="text-xs font-black text-gray-700 group-hover:text-secondary transition-colors">📊 Leaderboard</span>
        </a>
      </div>

      <!-- Informational Images -->
      <div class="mb-5 space-y-4">
        <img
          src="/img/rules-grow-avi.png"
          alt="Rules GROW AVI"
          class="w-full rounded-2xl shadow-sm"
        />
        <img
          src="/img/timestamp-grow-avi.png"
          alt="Timestamp Example"
          class="w-full rounded-2xl shadow-sm"
        />
      </div>

      <form @submit.prevent="submitForm">
        <!-- Dropdown Overlay -->
        <div v-if="showCategoryDropdown || showActivityDropdown" @click="closeDropdowns" class="fixed inset-0 z-30"></div>

        <!-- Data Diri Section -->
        <div class="card-game p-6 mb-8 relative overflow-hidden bg-white/90 backdrop-blur-sm z-10">
          <div class="absolute top-0 right-0 w-32 h-32 bg-blue-100 rounded-full blur-3xl opacity-50 -mr-10 -mt-10 pointer-events-none"></div>
          <h2 class="text-lg font-black text-gray-900 mb-5 relative z-10 flex items-center gap-2">
            <span class="text-xl">👤</span> Profil
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 relative z-10">
            <div>
              <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wider">NPK</label>
              <input
                type="text"
                disabled
                :value="auth.user?.npk"
                class="w-full px-4 py-3 rounded-2xl border-2 border-gray-100 bg-gray-50 text-gray-500 font-bold text-sm cursor-not-allowed"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wider">Nama</label>
              <input
                type="text"
                disabled
                :value="auth.user?.name || auth.user?.npk"
                class="w-full px-4 py-3 rounded-2xl border-2 border-gray-100 bg-gray-50 text-gray-500 font-bold text-sm cursor-not-allowed"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wider">Department</label>
              <input
                type="text"
                disabled
                :value="auth.user?.department || 'TBA'"
                class="w-full px-4 py-3 rounded-2xl border-2 border-gray-100 bg-gray-50 text-gray-500 font-bold text-sm cursor-not-allowed"
              />
            </div>
          </div>
        </div>

        <!-- Category Section -->
        <div class="card-game p-6 mb-8 relative bg-white z-40">
          <label class="block text-lg font-black text-gray-900 mb-5 flex items-center gap-2">
            <span class="text-xl">🎯</span> Pilih Kategori <span class="text-red-500 text-2xl leading-none">*</span>
          </label>
          <div class="relative">
            <button type="button" @click="showCategoryDropdown = !showCategoryDropdown; showActivityDropdown = false" class="w-full bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-4 py-3.5 text-sm font-bold flex justify-between items-center transition-colors hover:border-primary/50 focus:outline-none" :class="{'border-primary': showCategoryDropdown}">
              <span class="truncate pr-2">{{ selectedCategoryName || 'Pilih Kategori...' }}</span>
              <svg class="w-5 h-5 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showCategoryDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            <div v-if="showCategoryDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto overflow-x-hidden p-1 animate-fade-in">
              <div v-for="cat in categories" :key="cat.id" @click="selectCategory(cat.id)" class="px-4 py-3 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="form.categoryId === cat.id ? 'text-primary bg-blue-50' : 'text-gray-600'">
                {{ cat.name }}
              </div>
            </div>
          </div>
        </div>

        <!-- Activity Section -->
        <div
          v-if="form.categoryId"
          class="card-game p-6 mb-8 relative bg-white animate-slide-up z-30"
        >
          <div class="mb-5">
            <label class="block text-lg font-black text-gray-900 flex items-center gap-2">
              <span class="text-xl">🏃‍♂️</span> Kegiatan atau Aktivitas <span class="text-red-500 text-2xl leading-none">*</span>
            </label>
            <!-- <div v-if="form.categoryId === 6" class="mt-2 inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 border border-amber-200 rounded-lg text-amber-800 text-xs font-bold">
              <span>⚠️</span> Untuk kegiatan SPORT dll, harus sudah disetujui oleh Perusahaan
            </div> -->
          </div>

          <div class="relative mt-4">
            <button type="button" @click="showActivityDropdown = !showActivityDropdown; showCategoryDropdown = false" class="w-full bg-gray-50 border-2 border-gray-200 text-gray-700 rounded-xl px-4 py-3.5 text-sm font-bold flex justify-between items-center transition-colors hover:border-primary/50 focus:outline-none" :class="{'border-primary': showActivityDropdown}">
              <span class="truncate pr-2">{{ selectedActivityName || 'Pilih Kegiatan / Aktivitas...' }}</span>
              <svg class="w-5 h-5 text-gray-400 shrink-0 transition-transform" :class="{'rotate-180': showActivityDropdown}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </button>
            <div v-if="showActivityDropdown" class="absolute z-50 w-full mt-2 bg-white border border-gray-100 rounded-xl shadow-xl max-h-60 overflow-y-auto overflow-x-hidden p-1 animate-fade-in">
              <div v-for="act in filteredActivities" :key="act.id" @click="selectActivity(act)" class="px-4 py-3 text-sm font-bold rounded-lg hover:bg-blue-50 cursor-pointer transition-colors" :class="form.activityId === act.id ? 'text-primary bg-blue-50' : 'text-gray-600'">
                {{ act.name }}
              </div>
            </div>
          </div>

          <!-- Custom Activity Input -->
          <div v-if="isCustomActivitySelected" class="mt-4 animate-slide-up">
            <label class="block text-sm font-bold text-gray-700 mb-2">
              Sebutkan Kegiatan <span class="text-red-500">*</span>
            </label>
            <input
              type="text"
              v-model="form.customActivityType"
              required
              class="w-full border-2 border-gray-200 rounded-xl focus:border-primary focus:ring-2 focus:ring-primary/20 bg-white px-4 py-3.5 text-sm font-semibold transition-all"
              placeholder="Sebutkan kegiatan..."
            />
          </div>
        </div>

        <!-- Details & Evidence Section -->
        <div
          class="card-game p-6 mb-8 relative bg-white"
        >
          <label class="block text-lg font-black text-gray-900 mb-6 flex items-center gap-2">
            <span class="text-xl">📝</span> Detail & Lampiran
          </label>

          <div class="space-y-6">
            <!-- Nama Kegiatan / Aktivitas -->
            <div v-if="form.categoryId !== 8 && !isSubmitSS">
              <label class="block text-sm font-bold text-gray-700 mb-2">
                Nama Kegiatan/Aktivitas <span class="text-red-500">*</span>
              </label>
              <input
                type="text"
                v-model="form.customActivityName"
                required
                placeholder="Contoh: Lari Pagi 5KM"
                class="w-full px-4 py-3.5 rounded-2xl border-2 border-gray-200 bg-gray-50 focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary font-semibold outline-none transition-all text-sm placeholder-gray-400"
              />
            </div>

            <!-- Activity Date -->
            <div v-if="!isSubmitSS">
              <label class="block text-sm font-bold text-gray-700 mb-2">
                Tanggal Aktivitas <span class="text-red-500">*</span>
              </label>
              <input
                type="date"
                v-model="form.activityDate"
                required
                class="w-full px-4 py-3.5 rounded-2xl border-2 border-gray-200 bg-gray-50 focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary font-semibold outline-none transition-all text-sm"
              />
            </div>

            <!-- Custom Reference (Only for SS / Innovation) -->
            <div v-if="isSubmitSS">
              <label class="block text-sm font-bold text-gray-700 mb-2">
                Nomor SS<span class="text-red-500">*</span>
              </label>
              <input
                type="text"
                v-model="form.customReference"
                :required="isSubmitSS"
                placeholder="Contoh: SS-2026-001"
                class="w-full px-4 py-3.5 rounded-2xl border-2 border-gray-200 bg-gray-50 focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary font-semibold outline-none transition-all text-sm placeholder-gray-400"
              />
            </div>

            <!-- Evidence Upload -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">
                Lampiran File Pendukung <span class="text-red-500">*</span>
              </label>
              <div
                class="mt-1 flex flex-col items-center justify-center px-6 py-8 border-2 border-gray-300 border-dashed rounded-3xl hover:border-primary hover:bg-blue-50 transition-all cursor-pointer bg-gray-50 group"
                @click="triggerFileInput"
              >
                <div class="w-14 h-14 rounded-full bg-blue-100 flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
                  <svg class="h-7 w-7 text-primary" stroke="currentColor" fill="none" viewBox="0 0 48 48" aria-hidden="true">
                    <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <div class="text-center">
                  <span class="text-sm font-bold text-primary">Klik untuk memilih file</span>
                  <p class="text-xs text-gray-500 font-medium mt-1">Maksimal 5 file, (Format: JPG, PNG, WEBP), < 10MB</p>
                </div>
                <input
                  id="file-upload"
                  name="file-upload"
                  type="file"
                  class="sr-only"
                  multiple
                  accept=".jpg,.jpeg,.png,.webp"
                  @change="handleFileUpload"
                  ref="fileInput"
                />
              </div>

              <!-- File Preview List -->
              <div v-if="files.length > 0" class="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div
                  v-for="(file, index) in files"
                  :key="index"
                  class="flex items-center justify-between p-2.5 bg-white rounded-2xl border-2 border-gray-100 shadow-sm animate-pop"
                >
                  <div class="flex items-center space-x-3 overflow-hidden">
                    <div class="w-12 h-12 rounded-xl bg-gray-100 flex-shrink-0 overflow-hidden border border-gray-200">
                      <img :src="file.preview" class="w-full h-full object-cover" />
                    </div>
                    <span class="text-sm font-semibold text-gray-700 truncate max-w-[120px]">{{ file.name }}</span>
                  </div>
                  <button
                    type="button"
                    @click.stop="removeFile(index)"
                    class="w-8 h-8 rounded-full bg-red-50 text-red-500 flex items-center justify-center hover:bg-red-500 hover:text-white transition-colors flex-shrink-0 mr-1"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>



        <!-- Submit Button -->
        <div class="pb-10">
          <button
            type="submit"
            :disabled="isSubmitting"
            class="btn-game-primary w-full py-4 text-lg font-black rounded-2xl flex items-center justify-center gap-2"
          >
            <span v-if="isSubmitting" class="flex items-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses Pengajuan...
            </span>
            <span v-else>🚀 Kirim Pengajuan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "../../stores/auth";
import apiClient from "../../api/client";
import axios from "axios";
import { useToast } from "vue-toastification";
import { useRouter, useRoute } from "vue-router";

const auth = useAuthStore();
const toast = useToast();
const router = useRouter();
const route = useRoute();

const revisiId = computed(() => route.query.revisi_id);
const isEditMode = computed(() => !!revisiId.value);
const previousEvidence = ref<string[]>([]);

const categories = ref<any[]>([]);
const allActivities = ref<any[]>([]);

const fetchMasterData = async () => {
  try {
    const [catRes, actRes] = await Promise.all([
      apiClient.get("/categories"),
      apiClient.get("/activities"),
    ]);
    categories.value = catRes.data.data || [];
    console.log("ini cari categories", categories.value);
    allActivities.value = actRes.data.data || [];
  } catch (err) {
    console.error("Gagal mengambil daftar kategori", err);
  }
};

const fetchSubmissionForEdit = async (id: string) => {
  try {
    const res = await apiClient.get(`/submissions/${id}`);
    const data = res.data.data;
    if (data.status !== 'REJECTED_RESUBMIT') {
      toast.error('Hanya pengajuan yang dikembalikan untuk revisi yang dapat diperbaiki');
      router.push('/dashboard');
      return;
    }
    form.value.categoryId = data.category_id;
    await fetchActivities();
    form.value.activityId = data.activity_id;
    form.value.activityDate = data.activity_date;
    form.value.customActivityType = data.custom_activity_type?.Valid ? data.custom_activity_type.String : '';
    form.value.customActivityName = data.custom_reference?.Valid ? data.custom_reference.String : '';
    form.value.customReference = data.nomor_ss?.Valid ? data.nomor_ss.String : '';
    
    // Load old evidence as File objects so user can preview and delete
    if (data.evidence && data.evidence.length > 0) {
      for (let i = 0; i < data.evidence.length; i++) {
        const url = data.evidence[i];
        const filename = url.split('/').pop() || `lampiran_${i + 1}.jpg`;
        try {
          const resBlob = await axios.get(url, { responseType: 'blob' });
          const file = new File([resBlob.data], filename, { type: resBlob.data.type || 'image/jpeg' });
          const preview = URL.createObjectURL(file);
          files.value.push({ file, preview, name: filename });
        } catch (e) {
          console.warn(`Gagal memuat lampiran lama: ${filename}`, e);
        }
      }
    }
  } catch (err) {
    console.error("Gagal mengambil data pengajuan", err);
    toast.error('Gagal mengambil data pengajuan untuk revisi');
  }
};

const fetchActivities = async () => {
  await fetchMasterData();
};

onMounted(async () => {
  await fetchMasterData();
  if (revisiId.value) {
    await fetchSubmissionForEdit(revisiId.value as string);
  }
});

const form = ref<{
  categoryId: number | "";
  activityId: number | string;
  customActivityType: string;
  customActivityName: string;
  activityDate: string;
  customReference: string;
}>({
  categoryId: "",
  activityId: "",
  customActivityType: "",
  customActivityName: "",
  activityDate: "",
  customReference: "",
});

const files = ref<{ file: File; preview: string; name: string }[]>([]);
const fileInput = ref<any>(null);
const isSubmitting = ref(false);

const filteredActivities = computed(() => {
  if (!form.value.categoryId) return [];
  return allActivities.value.filter(
    (a) => a.category_id === form.value.categoryId,
  );
});

const isSubmitSS = computed(() => {
  const act = allActivities.value.find(a => a.id === form.value.activityId);
  return act && act.name === 'Submit SS';
});

const showCategoryDropdown = ref(false);
const showActivityDropdown = ref(false);

const selectedCategoryName = computed(() => {
  const cat = categories.value.find(c => c.id === form.value.categoryId);
  return cat ? cat.name : '';
});

const selectedActivityName = computed(() => {
  const act = allActivities.value.find(a => a.id === form.value.activityId);
  return act ? act.name : '';
});

const isCustomActivitySelected = computed(() => {
  const act = allActivities.value.find(a => a.id === form.value.activityId);
  return act ? act.is_custom_input : false;
});

const selectCategory = (id: number) => {
  form.value.categoryId = id;
  form.value.activityId = ""; 
  showCategoryDropdown.value = false;
};

const selectActivity = (act: any) => {
  form.value.activityId = act.id;
  if (!act.is_custom_input) {
    form.value.customActivityType = '';
  }
  showActivityDropdown.value = false;
};

const closeDropdowns = () => {
  showCategoryDropdown.value = false;
  showActivityDropdown.value = false;
};

const triggerFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click();
  }
};

const handleFileUpload = (event: any) => {
  const selectedFiles = event.target.files;
  if (!selectedFiles) return;

  if (files.value.length + selectedFiles.length > 5) {
    toast.error("Maksimal 5 file yang diperbolehkan.");
    return;
  }

  for (let i = 0; i < selectedFiles.length; i++) {
    const file = selectedFiles[i];
    if (file.size > 10 * 1024 * 1024) {
      toast.error(`File "${file.name}" melebihi 10MB!`);
      continue;
    }

    const preview = URL.createObjectURL(file);
    files.value.push({ file, preview, name: file.name });
  }

  // Reset input so the same file can be re-selected
  event.target.value = "";
};

const removeFile = (index: number) => {
  const removed = files.value[index];
  if (removed.preview.startsWith('blob:')) {
    URL.revokeObjectURL(removed.preview);
  }
  files.value.splice(index, 1);
};

const submitForm = async () => {
  if (!isSubmitSS.value && files.value.length === 0) {
    toast.error("Anda harus mengunggah setidaknya 1 bukti aktivitas.");
    return;
  }

  if (
    form.value.activityId === "other" &&
    !form.value.customActivityName.trim()
  ) {
    toast.error("Mohon isi nama kegiatan (Yang lain).");
    return;
  }

  isSubmitting.value = true;

  try {
    const formData = new FormData();

    if (isSubmitSS.value) {
      form.value.activityDate = new Date().toISOString().split('T')[0];
    }

    formData.append("activity_id", String(form.value.activityId));
    formData.append("activity_date", form.value.activityDate);

    let referenceValue = form.value.customReference;
    let actName = form.value.customActivityName;

    if (form.value.customActivityType) {
      formData.append("custom_activity_type", form.value.customActivityType);
    }
    
    if (actName && form.value.categoryId !== 8) {
      formData.append("custom_reference", actName);
    }
    
    if (referenceValue) {
      formData.append("nomor_ss", referenceValue);
    }

    files.value.forEach((f) => {
      formData.append("evidence", f.file);
    });

    let response;
    if (isEditMode.value) {
      response = await apiClient.put(`/submissions/${revisiId.value}/resubmit`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      toast.success("Pengajuan berhasil direvisi dan dikirim ulang!");
    } else {
      response = await apiClient.post("/submissions", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      toast.success(`Pengajuan berhasil dikirim! GROW ID: ${response.data.grow_id}`);
    }

    isSubmitting.value = false;

    // Reset form
    form.value = {
      categoryId: "",
      activityId: "",
      customActivityType: "",
      customActivityName: "",
      activityDate: "",
      customReference: "",
    };
    files.value.forEach((f) => URL.revokeObjectURL(f.preview));
    files.value = [];

    // Redirect to Dashboard after 1.5 seconds
    setTimeout(() => {
      router.push("/dashboard");
    }, 1500);
  } catch (err: any) {
    isSubmitting.value = false;
    toast.error(err.response?.data?.error || "Terjadi kesalahan saat mengirim pengajuan.");
    console.error(err);
  }
};
</script>

<style scoped>
/* Add a custom accent color for radio buttons to match primary color if needed */
input[type="radio"] {
  accent-color: #0069aa;
}
</style>
