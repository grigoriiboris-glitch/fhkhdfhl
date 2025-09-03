<template>
  <div class="container">
    <div class="header">
      <h1>Вход в систему</h1>
    </div>
    <div v-if="authError" class="error">{{ authError }}</div>
    <form @submit.prevent="onLogin">
      <div class="form-group">
        <label for="email">Email:</label>
        <input v-model="email" type="email" id="email" required />
      </div>
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input v-model="password" type="password" id="password" required />
      </div>
      <button type="submit" class="submit-btn" :disabled="loading">Войти</button>
    </form>
    <div class="links">
      <router-link to="/register">Регистрация</router-link>
      <router-link to="/">На главную</router-link>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import useAuth from '@/composition/useAuth'

const email = ref('')
const password = ref('')
const { login, error: authError, loading } = useAuth()

async function onLogin() {
  await login(email.value, password.value)
}
</script>
<style scoped>
/* ...скопируйте стили из login.html... */
</style>