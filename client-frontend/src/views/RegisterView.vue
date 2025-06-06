<template>
    <div class="register-container">
      <form class="register-form" @submit.prevent="register">
        <h2>Start Chat</h2>
        <input v-model="name" type="text" placeholder="Your Name" required />
        <input v-model="phone" type="tel" placeholder="Phone (+380...)" required />
        <button type="submit">Start</button>
      </form>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import apiClient from '../services/apiClient'
  
  export default {
    setup() {
      const name = ref('')
      const phone = ref('')
      const router = useRouter()
  
      const register = async () => {
        try {
          const res = await apiClient.post('/client/register', {
            name: name.value,
            phone: phone.value,
          })
  
          console.log('Registration response:', res.data)
          const { client, roomID } = res.data
  
          localStorage.setItem('client', JSON.stringify(client))
          localStorage.setItem('dialog_id', roomID)
  
          router.push('/chat')
        } catch (err) {
          alert('Registration failed')
          console.error(err)
        }
      }
  
      return { name, phone, register }
    }
  }
  </script>
  
  <style scoped>
  .register-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f6f8fa;
  }
  
  .register-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 0 12px rgba(0,0,0,0.1);
    width: 100%;
    max-width: 400px;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .register-form h2 {
    margin: 0;
    text-align: center;
  }
  
  .register-form input {
    padding: 0.75rem;
    font-size: 1rem;
    border-radius: 6px;
    border: 1px solid #ccc;
  }
  
  .register-form button {
    padding: 0.75rem;
    font-size: 1rem;
    background-color: #2c7be5;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  
  .register-form button:hover {
    background-color: #1a5fc8;
  }
  </style>
