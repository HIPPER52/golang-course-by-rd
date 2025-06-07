import axios from 'axios';
import { useUserStore } from '../store/userStore';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  const userStore = useUserStore();
  if (userStore.token) {
    config.headers['X-User-Token'] = userStore.token;
  }
  return config;
});

export default apiClient;
