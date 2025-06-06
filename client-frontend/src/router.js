import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from './views/RegisterView.vue'
import ChatView from './views/ChatView.vue'

const routes = [
  { path: '/', component: RegisterView },
  { path: '/chat', component: ChatView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
