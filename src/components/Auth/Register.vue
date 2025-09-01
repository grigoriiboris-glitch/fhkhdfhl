<template>
  <div class="register-container">
    <div class="register-form">
      <h2>{{ $t('auth.register') }}</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="name">{{ $t('auth.name') }}</label>
          <input
            type="text"
            id="name"
            v-model="name"
            required
            :placeholder="$t('auth.namePlaceholder')"
          />
        </div>
        <div class="form-group">
          <label for="email">{{ $t('auth.email') }}</label>
          <input
            type="email"
            id="email"
            v-model="email"
            required
            :placeholder="$t('auth.emailPlaceholder')"
          />
        </div>
        <div class="form-group">
          <label for="password">{{ $t('auth.password') }}</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
            :placeholder="$t('auth.passwordPlaceholder')"
          />
        </div>
        <div class="form-group">
          <label for="confirmPassword">{{ $t('auth.confirmPassword') }}</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            required
            :placeholder="$t('auth.confirmPasswordPlaceholder')"
          />
        </div>
        <div class="form-actions">
          <button type="submit" :disabled="loading || !passwordsMatch">
            {{ loading ? $t('auth.registering') : $t('auth.register') }}
          </button>
          <button type="button" @click="$emit('switch-to-login')">
            {{ $t('auth.haveAccount') }}
          </button>
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div v-if="!passwordsMatch && confirmPassword" class="error-message">
          {{ $t('auth.passwordsNotMatch') }}
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    return {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
      loading: false,
      error: ''
    }
  },
  computed: {
    passwordsMatch() {
      return this.password === this.confirmPassword
    }
  },
  methods: {
    async handleRegister() {
      if (!this.passwordsMatch) {
        return
      }

      this.loading = true
      this.error = ''
      
      try {
        const response = await fetch('/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: new URLSearchParams({
            name: this.name,
            email: this.email,
            password: this.password
          }),
          credentials: 'include'
        })

        if (response.ok) {
          this.$emit('register-success')
        } else {
          const data = await response.text()
          this.error = data || this.$t('auth.registerError')
        }
      } catch (err) {
        this.error = this.$t('auth.networkError')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.register-form h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.form-actions button {
  flex: 1;
  padding: 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.form-actions button[type="submit"] {
  background-color: #667eea;
  color: white;
}

.form-actions button[type="submit"]:hover:not(:disabled) {
  background-color: #5a6fd8;
}

.form-actions button[type="submit"]:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.form-actions button[type="button"] {
  background-color: #f8f9fa;
  color: #667eea;
  border: 1px solid #667eea;
}

.form-actions button[type="button"]:hover {
  background-color: #e9ecef;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  text-align: center;
}
</style> 