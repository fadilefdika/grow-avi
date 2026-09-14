import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null);
  const user = ref<any>(null);

  function setAuth(newToken: string, userData: any) {
    console.log('[DEBUG PINIA] setAuth terpanggil! Data user:', userData);
    token.value = newToken;
    user.value = userData;
  }

  function logout() {
    token.value = null;
    user.value = null;
  }

  return { token, user, setAuth, logout };
});
