import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import IndexPage from '@/pages/Index/Index'
import MindMapList from '@/components/MindMapList.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: MindMapList
  },
  {
    path: '/index',
    name: 'Index',
    component: IndexPage
  },
  { 
    path: '/edit', 
    name: 'Edit', 
    component: () => import(`./pages/Edit/Index.vue`) 
  },
  { 
    path: '/edit/:id', 
    name: 'EditMap', 
    component: () => import(`./pages/Edit/Index.vue`),
    props: true
  }
]

const router = createRouter({
  history: process.env.NODE_ENV === 'development' ? createWebHistory() : createWebHashHistory(),
  base: '/hyy-vue3-mindmap/',
  routes
})

export default router
