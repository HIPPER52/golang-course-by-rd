import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import AdminLayout from '../layouts/AdminLayout.vue';
import Chat from '../views/Chat.vue';
import Operators from '../views/Operators.vue';

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/chat',
      },
      {
        path: 'chat',
        name: 'Chat',
        component: Chat,
      },
      {
        path: 'operators',
        name: 'Operators',
        component: Operators,
        meta: { requiresAdmin: true },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const isLoggedIn = !!localStorage.getItem('token');
  const userRole = localStorage.getItem('role');

  if (to.meta.requiresAuth && !isLoggedIn) {
    return next('/login');
  }

  if (to.meta.requiresAdmin && userRole !== 'admin') {
    return next('/chat');
  }

  next();
});

export default router;
