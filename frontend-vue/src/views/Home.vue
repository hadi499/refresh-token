<template>
  <div class="profile-container">
  

    <!-- Loading State -->
    <div v-if="loading">
      <p>Loading user data...</p>
    </div> 

    <!-- Data Display with v-for -->
    <div v-else>
     
      <!-- Daftar semua user dari API -->
      <h2>All Users</h2>
      <div class="users-list">
        <div v-for="user in users" :key="user._id" class="user-card">
          <h3>{{ user.name }}</h3>
          <p><strong>Email:</strong> {{ user.email }}</p>
          <p><small>Joined: {{ formatDate(user.createdAt) }}</small></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAxios } from '../plugins/axios';
import { format } from 'date-fns';

// State
const users = ref([]); // Untuk menyimpan array users dari API
const currentUser = ref(null); // Untuk user yang sedang login
const loading = ref(false);
const error = ref(null);

import { useAuthStore } from '../stores/auth';


const authStore = useAuthStore();

// Get axios instance
const api = useAxios();

// Format tanggal
const formatDate = (dateString) => {
  return format(new Date(dateString), 'dd MMM yyyy HH:mm');
};

// Fetch data
const fetchUserData = async () => {
  loading.value = true;
  error.value = null;
  try {
    // Ambil data semua users
    const usersResponse = await api.get('/api/users');
    users.value = usersResponse.data;


  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to fetch user data';
    console.error('Error:', err);
  } finally {
    loading.value = false;
  }
};

// Fetch data on component mount
onMounted(() => {
  fetchUserData();
});
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.error {
  color: red;
  font-weight: bold;
}

.current-user {
  background-color: #f0f7ff;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.users-list {
  display: grid;
  gap: 15px;
}

.user-card {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 8px;
  border-left: 4px solid #4CAF50;
}

.user-card h3 {
  margin-top: 0;
  color: #333;
}

button {
  background-color: #4CAF50;
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}

button:hover {
  background-color: #45a049;
}
</style>