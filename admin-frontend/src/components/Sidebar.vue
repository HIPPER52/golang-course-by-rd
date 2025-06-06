<template>
    <aside class="sidebar">
      <ul class="menu">
        <li><router-link to="/chat">Chat</router-link></li>
        <li v-if="isAdmin"><router-link to="/operators">Operators</router-link></li>
        <li v-if="isAdmin"><router-link to="/statistics">Statistics</router-link></li>
      </ul>
      <button class="logout-btn" @click="logout">Log Out</button>
    </aside>
  </template>
  
  <script setup>
  import { useUserStore } from '../store/userStore'
  import { useRouter } from 'vue-router'
  
  const userStore = useUserStore()
  const router = useRouter()
  
  console.log('User Store:', userStore)
  const isAdmin = userStore.isAdmin
  
  function logout() {
    userStore.logout()
    router.push('/login')
  }
  </script>
  
  <style scoped>
  .sidebar {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: #f2f2f2;
    padding: 1rem;
  }
  .sidebar ul {
    list-style: none;
    padding: 0;
  }
  .sidebar li {
    margin-bottom: 1rem;
  }
  .sidebar a {
    text-decoration: none;
    color: #333;
    font-weight: bold;
  }
  .sidebar a.router-link-active {
    color: #007bff;
  }
  .sidebar button {
    margin-top: 2rem;
    padding: 0.5rem;
    background: #e74c3c;
    color: white;
    border: none;
    cursor: pointer;
    border-radius: 4px;
  }
  .logout-btn {
    background-color: #e74c3c;
    color: white;
    border: none;
    padding: 0.5rem;
    border-radius: 4px;
    cursor: pointer;
}
  </style>
