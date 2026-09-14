import axios from 'axios';
import { useAuthStore } from '../stores/auth';
import router from '../router';

// Create axios instance
const apiClient = axios.create({
  baseURL: '/api',
  withCredentials: true, // Important for cookies (refresh token)
});

// Request interceptor to add access token
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    if (authStore.token && config.headers) {
      config.headers.Authorization = `Bearer ${authStore.token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle 401 and silent refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // Attempt to refresh the token
        const response = await axios.post('/api/auth/refresh', {}, { withCredentials: true });
        
        if (response.data && response.data.access_token) {
          const authStore = useAuthStore();
          authStore.setAuth(response.data.access_token, authStore.user);
          
          // Retry the original request
          originalRequest.headers.Authorization = `Bearer ${response.data.access_token}`;
          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, force logout
        const authStore = useAuthStore();
        authStore.logout();
        router.push('/login');
      }
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;
