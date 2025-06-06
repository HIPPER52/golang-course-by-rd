import { defineStore } from 'pinia'
import { initSocket } from '@/services/socketService'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    operatorId: localStorage.getItem('operator_id') || null,
    role: localStorage.getItem('role') || null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.role === 'admin',
  },
  actions: {
    setUser({ token, role, operator_id }) {
      this.token = token
      this.role = role
      this.operatorId = operator_id
      localStorage.setItem('token', token)
      localStorage.setItem('role', role)
      localStorage.setItem('operator_id', operator_id)

      initSocket(token)
    },
    logout() {
      this.token = null
      this.role = null
      localStorage.clear()
    },
  },
})
