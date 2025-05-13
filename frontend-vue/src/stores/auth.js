// src/stores/auth.js
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { jwtDecode } from 'jwt-decode';

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter();

  // State (menggunakan ref dari Vue Composition API)
  const authTokens = ref(localStorage.getItem("authTokens")
    ? JSON.parse(localStorage.getItem("authTokens"))
    : null
  );

  const user = ref(localStorage.getItem("authTokens")
    ? jwtDecode(JSON.parse(localStorage.getItem("authTokens")).accessToken)
    : null
  );

  const loading = ref(true);

  // Actions
  const loginUser = async (formData) => {
    try {
      const response = await fetch("http://localhost:5000/api/users/auth", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password,
        }),
      });

      const data = await response.json();

      if (response.status === 200) {
        authTokens.value = data;
        user.value = jwtDecode(data.accessToken);
        localStorage.setItem("authTokens", JSON.stringify(data));
        router.push("/");
        return true;
      } else {
        return false;
      }
    } catch (error) {
      console.error("Login error:", error);
      return false;
    }
  };

  const logoutUser = async () => {
    try {
      if (authTokens.value?.refreshToken) {
        await fetch("http://localhost:5000/api/users/logout", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ token: authTokens.value.refreshToken }),
        });
      }

      authTokens.value = null;
      user.value = null;
      localStorage.removeItem("authTokens");
      router.push("/login");
    } catch (error) {
      console.error("Logout error:", error);
    }
  };

  // Efek samping (mirip dengan useEffect)
  const initializeStore = () => {
    if (authTokens.value) {
      user.value = jwtDecode(authTokens.value.accessToken);
    }
    loading.value = false;
  };

  // Memanggil initialize
  initializeStore();

  return {
    // State yang dapat diakses
    user,
    authTokens,
    loading,

    // Actions
    loginUser,
    logoutUser,
  };
});