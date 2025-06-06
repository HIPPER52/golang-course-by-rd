import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Chat from '../views/Chat.vue'
import Operators from '../views/Operators.vue'
import Statistics from '../views/Statistics.vue'
import { useUserStore } from '../store/userStore'
import AdminLayout from '../layouts/AdminLayout.vue'

const routes = [
    {
        path: '/',
        component: AdminLayout,
        children: [
          {
            path: 'chat',
            component: Chat,
            meta: { requiresAuth: true },
          },
          {
            path: 'operators',
            component: Operators,
            meta: { requiresAuth: true, requiresAdmin: true },
          },
          {
            path: 'statistics',
            component: Statistics,
            meta: { requiresAuth: true, requiresAdmin: true },
          },
        ],
    },
    {
        path: '/login',
        component: Login,
    }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 🔒 Защита маршрутов
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    return next('/login')
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    return next('/chat')
  }

  next()
})

export default router
