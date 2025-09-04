<template>
  <div id="app"> 
    <AuthWrapper v-if="!User" @auth-success="handleAuthSuccess" />
    <router-view v-else-if="User && routes.includes($route.name)" />
    <div v-else-if="!routes.includes($route.name)" class="app-container">
      <BaseLayout/>
    </div>

  </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import AuthWrapper from './components/Auth/AuthWrapper.vue'
import BaseLayout from './components/base/components/Layout/Layout.vue'

import useAuth from '@/composition/useAuth';

const { me, logout, User } = useAuth();

const isAuthenticated = ref(false)
const userInfo = ref(null)

const checkAuthentication = async () => {
  const authenticated = await me()
  if (authenticated) {
    isAuthenticated.value = true
    userInfo.value = authenticated
  }
}

const routes = ref(['EditMap', 'Edit']);

const handleAuthSuccess = async () => {
  isAuthenticated.value = true
  userInfo.value = await me()
}

// const handleLogout = async () => {
//   const success = await logout()
//   if (success) {
//     isAuthenticated.value = false
//     userInfo.value = null
//   }
// }

onMounted(async () => {
  await checkAuthentication()
})
</script>

<style scoped lang="scss">
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