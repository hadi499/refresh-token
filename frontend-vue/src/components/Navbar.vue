<template>
  <nav class="bg-blue-600 p-4 text-white shadow-md sticky top-0 z-50">
    <div class="w-full mx-auto flex justify-between items-center md:px-8">
      <!-- Left Side -->
      <div class="flex items-center space-x-8">
        <router-link to="/"
          class="text-xl font-medium hover:text-yellow-300 transition-colors duration-200 inline-block"
          active-class="text-yellow-300">
          MyApp
        </router-link>

        
      </div>

      <!-- Right Side -->
      <div class="flex items-center space-x-4">
        <template v-if="!authStore.authTokens">
          <router-link to="/login"
            class="px-4 py-2 rounded hover:bg-blue-500 transition-colors duration-200 font-medium"
            active-class="bg-blue-700">
            Login
          </router-link>          
        </template>
        <template v-else>
          <span class=" mr-2 truncate">{{ authStore.user?.email }}</span>
          <button @click="logout"
            class="px-4 py-2 rounded bg-red-500 hover:bg-red-600 transition-colors duration-200 font-medium">
            Logout
          </button>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const logout = async () => {
  try {
    await authStore.logoutUser();
    router.push('/login');
  } catch (error) {
    console.error('Logout error:', error);
  }
};
</script>

<style scoped>
/* Tambahan agar transisi tetap halus tanpa mengubah ukuran */
.router-link-active {
  transition: color 0.2s ease;
}
</style>
