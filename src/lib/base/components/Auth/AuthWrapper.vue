<template>
  <div class="auth-wrapper">
    <Login 
      v-if="!showRegister" 
      @login-success="handleLoginSuccess"
      @switch-to-register="showRegister = true"
    />
    <Register 
      v-else 
      @register-success="handleRegisterSuccess"
      @switch-to-login="showRegister = false"
    />
  </div>
</template>

<script>
import Login from './Login.vue'
import Register from './Register.vue'

export default {
  name: 'AuthWrapper',
  components: {
    Login,
    Register
  },
  data() {
    return {
      showRegister: false
    }
  },
  methods: {
    handleLoginSuccess() {
      this.$emit('auth-success')
    },
    handleRegisterSuccess() {
      this.showRegister = false
      this.$emit('auth-success')
    }
  }
}
</script>

<style scoped>
.auth-wrapper {
  width: 100%;
  height: 100vh;
}
</style> 