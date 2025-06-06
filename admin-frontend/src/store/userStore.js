import { defineStore } from 'pinia'
import { initSocket } from '@/services/socketService'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    role: localStorage.getItem('role') || null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.role === 'admin',
  },
  actions: {
    setUser({ token, role }) {
      this.token = token
      this.role = role
      localStorage.setItem('token', token)
      localStorage.setItem('role', role)

      initSocket(token)
    },
    logout() {
      this.token = null
      this.role = null
      localStorage.clear()
    },
  },
})
