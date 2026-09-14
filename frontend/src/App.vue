<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from './stores/auth';
import BottomNav from './components/layout/BottomNav.vue';

const route = useRoute();
const authStore = useAuthStore();

const showNav = computed(() => {
  return authStore.token && route.name !== 'login';
});
</script>

<template>
  <div class="min-h-screen bg-gray-50 font-sans text-gray-900 pb-16">
    <!-- Main Content -->
    <router-view />

    <!-- Bottom Navigation -->
    <BottomNav v-if="showNav" />
  </div>
</template>

<style scoped>
</style>
