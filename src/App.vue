<template>
  <div id="app">
    <AuthWrapper 
      v-if="!isAuthenticated" 
      @auth-success="handleAuthSuccess"
    />
    <div v-else class="app-container">
      <Header>
        <template #right>
          <div class="user-info">
            <span>{{ userInfo?.name || userInfo?.email }}</span>
            <button @click="handleLogout" class="logout-btn">
              {{ $t('auth.logout') }}
            </button>
          </div>
        </template>
      </Header>
      <main class="app-main">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script>
import AuthWrapper from './components/Auth/AuthWrapper.vue'
import Header from './pages/Index/components/Header.vue'
import { checkAuth, getUserInfo, logout } from './utils/auth'

export default {
  name: 'App',
  components: {
    AuthWrapper,
    Header
  },
  data() {
    return {
      isAuthenticated: false,
      userInfo: null
    }
  },
  async mounted() {
    await this.checkAuthentication()
  },
  methods: {
    async checkAuthentication() {
      const authenticated = await checkAuth()
      if (authenticated) {
        this.isAuthenticated = true
        this.userInfo = await getUserInfo()
      }
    },
    async handleAuthSuccess() {
      this.isAuthenticated = true
      this.userInfo = await getUserInfo()
    },
    async handleLogout() {
      const success = await logout()
      if (success) {
        this.isAuthenticated = false
        this.userInfo = null
      }
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
  margin: 0;
  padding: 0;
}

.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.header-content h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info span {
  font-weight: 500;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.app-main {
  flex: 1;
  overflow: hidden;
}
</style>
