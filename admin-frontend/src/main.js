import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useUserStore } from './store/userStore'
import { initSocket } from './services/socketService'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

const userStore = useUserStore()

if (userStore.token) {
  initSocket(userStore.token)
}

app.mount('#app')
