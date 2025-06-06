<template>
    <div class="chat-box">
      <h2>Chat</h2>
      <div class="messages" ref="messagesContainer">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message', msg.sender === 'client' ? 'from-client' : 'from-operator']"
        >
          {{ msg.text }}
        </div>
      </div>
  
      <form class="send-form" @submit.prevent="send">
        <input v-model="text" placeholder="Type your message..." />
        <button>Send</button>
      </form>
    </div>
  </template>
  
  <script setup>
import { ref, onMounted, nextTick } from 'vue'
import { connectSocket, sendSocketMessage, subscribe } from '../socket'

const messages = ref([])
const text = ref('')
const messagesContainer = ref(null)

const client = JSON.parse(localStorage.getItem('client') || '{}')
const dialogId = localStorage.getItem('dialog_id')

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function send() {
  const content = text.value.trim()
  if (!content) return

  sendSocketMessage('message', {
    room_id: dialogId,
    text: content,
  })

  text.value = ''
}

onMounted(() => {
  const dialogId = localStorage.getItem('dialog_id')
  console.log('Connecting to dialog:', dialogId)

  connectSocket(client.id, dialogId)

  subscribe('message', (data) => {
    if (data.room_id !== dialogId) return

    const isFromClient = data.sender_id === client.id

    messages.value.push({
      id: Date.now(),
      text: data.text,
      sender: isFromClient ? 'client' : 'operator',
    })

    nextTick(scrollToBottom)
  })
})
</script>
  
  <style scoped>
.chat-box {
  max-width: 600px;
  height: 80vh;
  margin: 2rem auto;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  padding: 1rem;
}

.chat-box h2 {
  margin-bottom: 1rem;
  font-size: 1.4rem;
  text-align: center;
}

.messages {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
  padding: 0.5rem;
}

.message {
  max-width: 70%;
  padding: 0.6rem 1rem;
  border-radius: 16px;
  word-break: break-word;
  font-size: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #ddd;
}

.from-client {
  align-self: flex-end;
  background-color: #d1eaff;
  border-color: #8ecbff;
  color: black;
}

.from-operator {
  align-self: flex-start;
  background-color: #f0f0f0;
  border-color: #ccc;
  color: black;
}

form {
  display: flex;
  gap: 0.5rem;
}

input[type="text"] {
  flex: 1;
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 1rem;
}

button {
  padding: 0.5rem 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.send-form {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

.send-form input {
  flex: 1;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  border-radius: 6px;
  border: 1px solid #ccc;
}

.send-form button {
  padding: 0.75rem 1.2rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

.send-form button:hover {
  background-color: #0056b3;
}
</style>
