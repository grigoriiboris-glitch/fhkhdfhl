<template>
  <div class="container">
    <div class="header">
      <h1>Регистрация</h1>
    </div>
    <div v-if="authError" class="error">{{ authError }}</div>
    <form @submit.prevent="onRegister">
      <div class="form-group">
        <label for="name">Имя:</label>
        <input v-model="name" type="text" id="name" required />
      </div>
      <div class="form-group">
        <label for="email">Email:</label>
        <input v-model="email" type="email" id="email" required />
      </div>
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input v-model="password" type="password" id="password" required minlength="6" />
      </div>
      <div class="form-group">
        <label for="confirm_password">Подтвердите пароль:</label>
        <input v-model="confirm_password" type="password" id="confirm_password" required minlength="6" />
      </div>
      <button type="submit" class="submit-btn" :disabled="loading">Зарегистрироваться</button>
    </form>
    <div class="links">
      <router-link to="/login">Уже есть аккаунт? Войти</router-link>
      <router-link to="/">На главную</router-link>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import useAuth from '@/composition/useAuth'

const name = ref('')
const email = ref('')
const password = ref('')
const confirm_password = ref('')
const { register, error: authError, loading } = useAuth()

async function onRegister() {
  if (password.value !== confirm_password.value) {
    authError.value = 'Пароли не совпадают'
    return
  }
  await register(name.value, email.value, password.value)
}
</script>
<style scoped>
/* ...скопируйте стили из register.html... */
</style>