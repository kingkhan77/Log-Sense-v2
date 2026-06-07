import axios from 'axios';

const api = axios.create({ baseURL: '/api/v1' });

api.interceptors.request.use(cfg => {
  const token = localStorage.getItem('token');
  if (token) cfg.headers.Authorization = `Bearer ${token}`;
  return cfg;
});

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);

export default api;

export function jwtPayload() {
  try {
    return JSON.parse(atob(localStorage.getItem('token').split('.')[1]));
  } catch {
    return {};
  }
}
