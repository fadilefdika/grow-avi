<template>
  <div class="min-h-screen bg-gray-50 pb-12">
    <!-- Header Image -->
    <div class="w-full">
      <img
        src="/img/header-grow-avi.png"
        alt="GROW Evidence Form"
        class="w-full h-auto object-cover"
      />
    </div>

    <div class="max-w-md mx-auto px-4 mt-6">
      <!-- Description Section -->
      <div
        class="bg-white rounded-xl shadow-sm p-6 mb-6 border border-gray-100 border-t-8 border-t-primary"
      >
        <h1 class="text-2xl font-bold text-gray-900 mb-4">
          GROW Evidence Form
        </h1>
        <p class="text-gray-700 text-sm leading-relaxed mb-4">
          Formulir ini disusun untuk mendukung proses dokumentasi dan penilaian
          kegiatan dalam program GROW.
        </p>
        <p class="text-gray-700 text-sm leading-relaxed mb-4">
          Setiap karyawan diminta mengunggah bukti aktivitas sesuai kategori
          yang tersedia. Penilaian akan dilakukan berdasarkan kelengkapan,
          relevansi, dan kesesuaian bukti yang diunggah. Harap pastikan seluruh
          data diisi dengan benar dan setiap file pendukung dapat diakses dengan
          baik.
        </p>
        <p class="text-gray-700 text-sm leading-relaxed font-medium">
          Siapkan bukti terbaikmu, unggah dengan rapi, dan kumpulkan POIN-mu🚀
        </p>
        <p class="text-primary font-bold mt-4 text-lg">#let'sGROW</p>
        <div
          class="mt-4 p-3 bg-blue-50 text-primary text-xs rounded-lg font-medium border border-blue-100"
        >
          Periode aktivitas mulai dari 1 Januari 2026 sampai 31 Desember 2026
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="grid grid-cols-2 gap-3 mb-6">
        <a
          href="https://heyzine.com/flip-book/a0c4abfab8.html#page/1"
          target="_blank"
          rel="noopener noreferrer"
          class="flex flex-col items-center justify-center p-3 bg-white rounded-xl shadow-sm border border-gray-100 hover:border-primary transition-colors"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6 text-primary mb-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
            />
          </svg>
          <span class="text-xs font-semibold text-gray-800">Guide Book</span>
        </a>
        <a
          href="https://datastudio.google.com/embed/u/0/reporting/e32866a0-4028-403f-87aa-d5e2022b9931/page/T2DvF"
          target="_blank"
          rel="noopener noreferrer"
          class="flex flex-col items-center justify-center p-3 bg-white rounded-xl shadow-sm border border-gray-100 hover:border-secondary transition-colors"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6 text-secondary mb-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
            />
          </svg>
          <span class="text-xs font-semibold text-gray-800">Leaderboard</span>
        </a>
      </div>

      <!-- Informational Images -->
      <div class="mb-6 space-y-4">
        <img
          src="/img/rules-grow-avi.png"
          alt="Rules GROW AVI"
          class="w-full rounded-xl shadow-sm"
        />
        <img
          src="/img/timestamp-grow-avi.png"
          alt="Timestamp Example"
          class="w-full rounded-xl shadow-sm"
        />
      </div>

      <form @submit.prevent="submitForm">
        <!-- User Info Section -->
        <div
          class="bg-white rounded-xl shadow-sm p-6 mb-4 border-t-8 border-t-gray-500"
        >
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1"
                >NPK</label
              >
              <input
                type="text"
                disabled
                :value="auth.user?.npk"
                class="w-full px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-500 text-sm cursor-not-allowed"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1"
                >Nama</label
              >
              <input
                type="text"
                disabled
                :value="auth.user?.name || auth.user?.npk"
                class="w-full px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-500 text-sm cursor-not-allowed"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1"
                >Department</label
              >
              <input
                type="text"
                disabled
                :value="auth.user?.department || 'TBA'"
                class="w-full px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-500 text-sm cursor-not-allowed"
              />
            </div>
          </div>
        </div>

        <!-- Category Section -->
        <div
          class="bg-white rounded-xl shadow-sm p-6 mb-4 border border-gray-100"
        >
          <label class="block text-base font-medium text-gray-900 mb-4"
            >Kategori Leaderboard <span class="text-red-500">*</span></label
          >
          <div class="space-y-4">
            <label
              v-for="cat in categories"
              :key="cat.id"
              class="flex items-center space-x-3 cursor-pointer group"
            >
              <input
                type="radio"
                v-model="form.categoryId"
                :value="cat.id"
                class="h-5 w-5 text-primary border-gray-300 focus:ring-primary cursor-pointer"
                required
              />
              <span class="text-gray-700 group-hover:text-gray-900 text-sm">{{
                cat.name
              }}</span>
            </label>
          </div>
        </div>

        <!-- Activity Section -->
        <div
          v-if="form.categoryId"
          class="bg-white rounded-xl shadow-sm p-6 mb-4 border border-gray-100"
        >
          <div class="mb-2">
            <label class="block text-base font-medium text-gray-900"
              >Kegiatan atau Aktivitas
              <span class="text-red-500">*</span></label
            >
            <p v-if="form.categoryId === 6" class="text-xs text-gray-500 mt-1">
              Untuk kegiatan SPORT dll, harus sudah disetujui oleh Perusahaan
            </p>
          </div>

          <div class="space-y-4 mt-4">
            <template v-for="act in filteredActivities" :key="act.id">
              <label
                v-if="!act.is_custom_input"
                class="flex items-center space-x-3 cursor-pointer group"
              >
                <input
                  type="radio"
                  v-model="form.activityId"
                  :value="act.id"
                  class="h-5 w-5 text-primary border-gray-300 focus:ring-primary cursor-pointer"
                  required
                />
                <span class="text-gray-700 group-hover:text-gray-900 text-sm">{{
                  act.name
                }}</span>
              </label>

              <!-- Yang lain option -->
              <label
                v-else
                class="flex items-center space-x-3 cursor-pointer group"
              >
                <input
                  type="radio"
                  v-model="form.activityId"
                  :value="act.id"
                  class="h-5 w-5 text-primary border-gray-300 focus:ring-primary cursor-pointer"
                  required
                />
                <div class="flex items-center w-full">
                  <span class="text-gray-700 text-sm whitespace-nowrap mr-2"
                    >{{ act.name }}:</span
                  >
                  <input
                    type="text"
                    v-model="form.customActivityType"
                    :disabled="form.activityId !== act.id"
                    :required="form.activityId === act.id"
                    class="w-full border-b border-gray-300 focus:border-primary focus:outline-none bg-transparent px-1 py-0.5 text-sm transition-colors disabled:opacity-50"
                    placeholder="Sebutkan..."
                  />
                </div>
              </label>
            </template>
          </div>
        </div>

        <!-- Details & Evidence Section -->
        <div
          class="bg-white rounded-xl shadow-sm p-6 mb-6 border border-gray-100"
        >
          <!-- Nama Kegiatan / Aktivitas -->
          <div class="mb-5" v-if="form.categoryId !== 8">
            <label class="block text-sm font-medium text-gray-900 mb-2"
              >Nama Kegiatan/Aktivitas
              <span class="text-red-500">*</span></label
            >
            <input
              type="text"
              v-model="form.customActivityName"
              required
              placeholder="Contoh: Lari Pagi 5KM"
              class="w-full px-4 py-2.5 rounded-md border border-gray-300 focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-sm"
            />
          </div>

          <!-- Activity Date -->
          <div class="mb-5">
            <label class="block text-sm font-medium text-gray-900 mb-2"
              >Tanggal Aktivitas <span class="text-red-500">*</span></label
            >
            <input
              type="date"
              v-model="form.activityDate"
              required
              class="w-full px-4 py-2.5 rounded-md border border-gray-300 focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-sm"
            />
          </div>

          <!-- Custom Reference (Only for SS / Innovation) -->
          <div class="mb-5" v-if="form.categoryId === 4">
            <label class="block text-sm font-medium text-gray-900 mb-2"
              >Nomor SS / Referensi (Opsional)</label
            >
            <input
              type="text"
              v-model="form.customReference"
              placeholder="Contoh: SS-2026-001"
              class="w-full px-4 py-2.5 rounded-md border border-gray-300 focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-sm"
            />
          </div>

          <!-- Evidence Upload -->
          <div>
            <label class="block text-sm font-medium text-gray-900 mb-2"
              >Unggah Bukti Foto <span class="text-red-500">*</span></label
            >
            <div
              class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-lg hover:border-primary transition-colors cursor-pointer bg-gray-50"
              @click="triggerFileInput"
            >
              <div class="space-y-1 text-center">
                <svg
                  class="mx-auto h-10 w-10 text-gray-400"
                  stroke="currentColor"
                  fill="none"
                  viewBox="0 0 48 48"
                  aria-hidden="true"
                >
                  <path
                    d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                <div class="flex text-sm text-gray-600 justify-center">
                  <label
                    for="file-upload"
                    class="relative cursor-pointer rounded-md font-medium text-primary hover:text-blue-800 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-primary"
                  >
                    <span>Tambahkan File</span>
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
                  </label>
                </div>
                <p class="text-xs text-gray-500">Maks 5 file, < 10MB</p>
              </div>
            </div>

            <!-- File Preview List -->
            <div v-if="files.length > 0" class="mt-4 space-y-2">
              <div
                v-for="(file, index) in files"
                :key="index"
                class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200"
              >
                <div class="flex items-center space-x-3 overflow-hidden">
                  <div
                    class="w-10 h-10 rounded bg-gray-200 flex-shrink-0 overflow-hidden"
                  >
                    <img
                      :src="file.preview"
                      class="w-full h-full object-cover"
                    />
                  </div>
                  <span class="text-sm text-gray-700 truncate">{{
                    file.name
                  }}</span>
                </div>
                <button
                  type="button"
                  @click.stop="removeFile(index)"
                  class="text-red-500 hover:text-red-700 p-1"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div
          v-if="errorMessage"
          class="mb-4 p-3 bg-red-50 text-red-700 rounded-lg text-sm border border-red-100"
        >
          {{ errorMessage }}
        </div>

        <!-- Submit Button -->
        <div class="flex justify-end items-center bg-gray-100 p-4 rounded-xl">
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex justify-center py-2 px-6 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-blue-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isSubmitting">Mengirim...</span>
            <span v-else>Kirim</span>
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
import { useToast } from "vue-toastification";
import { useRouter } from "vue-router";

const auth = useAuthStore();
const toast = useToast();
const router = useRouter();

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
    console.error("Failed to load master data", err);
  }
};

onMounted(() => {
  fetchMasterData();
});

const form = ref({
  categoryId: "",
  activityId: "",
  customActivityType: "",
  customActivityName: "",
  activityDate: "",
  customReference: "",
});

const files = ref<any[]>([]);
const fileInput = ref<any>(null);
const isSubmitting = ref(false);
const errorMessage = ref("");

const filteredActivities = computed(() => {
  if (!form.value.categoryId) return [];
  return allActivities.value.filter(
    (a) => a.category_id === form.value.categoryId,
  );
});

const triggerFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click();
  }
};

const handleFileUpload = (event: any) => {
  const selectedFiles = event.target.files;
  if (!selectedFiles) return;

  if (files.value.length + selectedFiles.length > 5) {
    errorMessage.value = "Maksimal 5 file yang diperbolehkan.";
    return;
  }

  errorMessage.value = "";

  for (let i = 0; i < selectedFiles.length; i++) {
    const file = selectedFiles[i];
    if (file.size > 10 * 1024 * 1024) {
      errorMessage.value = "Terdapat file yang ukurannya melebihi 10MB.";
      continue;
    }

    // Create preview URL
    file.preview = URL.createObjectURL(file);
    files.value.push(file);
  }

  // Reset input
  event.target.value = "";
};

const removeFile = (index: number) => {
  URL.revokeObjectURL(files.value[index].preview);
  files.value.splice(index, 1);
};

const submitForm = async () => {
  if (files.value.length === 0) {
    errorMessage.value = "Anda harus mengunggah setidaknya 1 bukti aktivitas.";
    return;
  }

  if (
    form.value.activityId === "other" &&
    !form.value.customActivityName.trim()
  ) {
    errorMessage.value = "Mohon isi nama kegiatan (Yang lain).";
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = "";

  try {
    const formData = new FormData();
    formData.append("activity_id", form.value.activityId);
    formData.append("activity_date", form.value.activityDate);

    let referenceValue = form.value.customReference;
    
    // Gabungkan Nama Kegiatan Custom ke dalam Custom Reference 
    // agar datanya tersimpan di backend (kolom custom_reference)
    if (form.value.customActivityName && form.value.categoryId !== 8) {
      if (referenceValue) {
        referenceValue = `${form.value.customActivityName} (${referenceValue})`;
      } else {
        referenceValue = form.value.customActivityName;
      }
    }
    
    if (referenceValue) {
      formData.append("custom_reference", referenceValue);
    }

    files.value.forEach((file) => {
      formData.append("evidence", file);
    });

    const response = await apiClient.post("/submissions", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    isSubmitting.value = false;
    toast.success(`Pengajuan berhasil dikirim! GROW ID: ${response.data.grow_id}`);

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
    errorMessage.value =
      err.response?.data?.error || "Terjadi kesalahan saat mengirim pengajuan.";
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
