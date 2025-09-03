import { reactive, toRef } from "vue";
import router from "@/router";
import User from "@/models/User";
import axios from "axios";

const state = reactive({
  isAuthenticated: false,
  user: null,
  User: null,
  loading: false,
  error: ''
});

axios.defaults.withCredentials = true;

async function checkAuth() {
  try {
    const response = await axios.get('/auth/check');
    state.isAuthenticated = response.status === 200;
    return state.isAuthenticated;
  } catch (e) {
    state.isAuthenticated = false;
    return false;
  }
}

async function me() {
  try {
    const response = await axios.get('/auth/user');
    if (response.status === 200) {
      const data = response.data;
      state.user = data;
      state.User = data ? new User(data) : null;
      state.isAuthenticated = true;
      return state.user;
    }
    state.user = null;
    state.User = null;
    state.isAuthenticated = false;
    return null;
  } catch (e) {
    state.user = null;
    state.User = null;
    state.isAuthenticated = false;
    return null;
  }
}

async function login(email, password) {
  state.loading = true;
  state.error = '';
  try {
    await axios.post(
      '/auth/login',
      { email, password },
    );
    await me();
    await router.push('/');
    return true;
  } catch (e) {
    state.error = e.response?.data || e.message || 'Login error';
    return false;
  } finally {
    state.loading = false;
  }
}

async function register(name, email, password) {
  state.loading = true;
  state.error = '';
  try {
    await axios.post(
      '/auth/register',
      new URLSearchParams({ name, email, password }),
      { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
    );
    await me();
    await router.push('/');
    return true;
  } catch (e) {
    state.error = e.response?.data || e.message || 'Register error';
    return false;
  } finally {
    state.loading = false;
  }
}

async function logout() {
  try {
    await axios.get('/auth/logout');
  } finally {
    state.user = null;
    state.User = null;
    state.isAuthenticated = false;
    await router.push('/login');
  }
}

export default function useAuth() {
  return {
    isAuthenticated: toRef(state, 'isAuthenticated'),
    user: toRef(state, 'user'),
    User: toRef(state, 'User'),
    loading: toRef(state, 'loading'),
    error: toRef(state, 'error'),
    checkAuth,
    me,
    login,
    register,
    logout
  };
}
