import axios from "axios";
import { getAccessToken, clearAccessToken } from "./tokenStore";

const API_BASE = "/api"; // proxy will forward to :8080

const api = axios.create({
  baseURL: API_BASE,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

api.interceptors.response.use(
  (r) => r,
  (error) => {
    if (error.response?.status === 401) {
      clearAccessToken();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
