// src/plugins/axios.js
import axios from "axios";
import { jwtDecode } from "jwt-decode";
import dayjs from "dayjs";
import { useAuthStore } from "../stores/auth";

const baseURL = "http://localhost:5000";

export const useAxios = () => {
  const authStore = useAuthStore();

  const axiosInstance = axios.create({
    baseURL,
    headers: { Authorization: `Bearer ${authStore.authTokens?.accessToken}` }
  });

  axiosInstance.interceptors.request.use(async (req) => {
    if (!authStore.authTokens) return req;

    const user = jwtDecode(authStore.authTokens.accessToken);
    const isExpired = dayjs.unix(user.exp).diff(dayjs()) < 1;

    if (!isExpired) return req;

    try {
      const response = await axios.post(`${baseURL}/api/users/refresh/`, {
        token: authStore.authTokens.refreshToken,
      });

      // Update auth tokens in pinia store and localStorage
      const newTokens = response.data;
      authStore.authTokens = newTokens;
      localStorage.setItem("authTokens", JSON.stringify(newTokens));

      // Update request header
      req.headers.Authorization = `Bearer ${newTokens.accessToken}`;
    } catch (error) {
      // Jika refresh token gagal, logout user
      await authStore.logoutUser();
    }

    return req;
  });

  return axiosInstance;
};