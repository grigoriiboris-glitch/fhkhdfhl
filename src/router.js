import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import IndexPage from '@/pages/Index/Index'
import MindMapList from '@/components/MindMapList.vue'
import Login from '@/components/Auth/Login.vue'
import Register from '@/components/Auth/Register.vue'
import useAuth from './composition/useAuth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: MindMapList,
     meta: {
      title: 'Мои карты',
      roles: [1, 2,3],
    },
  },
  {
    path: '/index',
    name: 'Index',
    component: IndexPage,
     meta: {
      title: 'Index',
      roles: [1, 2,3],
    },
  },
  { 
    path: '/edit', 
    name: 'Edit', 
    component: () => import(`./pages/Edit/Index.vue`),
     meta: {
      title: 'Edit',
      roles: [1, 2,3],
    },
  },
  { 
    path: '/edit/:id', 
    name: 'EditMap', 
    component: () => import(`./pages/Edit/Index.vue`),
    props: true,
     meta: {
      title: 'EditMap',
      roles: [1, 2,3],
    },
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
     meta: {
      title: 'Login',
      roles: [1, 2,3],
    },
  },
  {
    path: '/profile',
    name: 'profile',
    component: Login,
     meta: {
      title: 'profile',
      roles: [1, 2,3],
    },
  },
  {
    path: '/tariffs',
    name: 'tariffs',
    component: Login,
     meta: {
      title: 'tariffs',
      roles: [1, 2,3],
    },
  },
  
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: {
      title: 'Register',
      roles: [1, 2,3],
    },
  }
]

const router = createRouter({
  history: process.env.NODE_ENV === 'development' ? createWebHistory() : createWebHashHistory(),
  base: '/hyy-vue3-mindmap/',
  routes
})

router.beforeEach(async (to, from, next) => {
  const { User, login, me, logout } = useAuth();

  const dataRoute = ['login', 'register', 'forgot'];

  if (!dataRoute.includes(to.name) && !User.value) {
    await me();
    if (!User.value) {
      next('/login')
    }
  }

  if (to.name === "login" && User.value) {
    log(User.value)
    next('/');
  }

  // if (to.meta.roles && !to.meta.roles.includes(User.value.role_id)) {
  //   next('/');
  // }

  // if (to.meta.show === 'admin' && !User.value.isSuperAdmin()) {
  //   next('/');
  // }

  next();
});

export default router;
