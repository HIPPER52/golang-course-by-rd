<template>
    <div class="dialog-layout">
      <div class="dialog-panel">
        <h3>Dialog with {{ dialog.client_name }}</h3>
  
        <div ref="messagesContainer" class="messages">
          <div
            v-for="msg in messages"
            :key="msg.id"
            :class="['message', msg.sender === 'operator' ? 'from-operator' : 'from-client']"
          >
            <p>{{ msg.text }}</p>
          </div>
        </div>
  
        <form class="send-form" @submit.prevent="submit">
          <input v-model="newMessage" placeholder="Type your message..." />
          <button type="submit">Send</button>
        </form>
  
        <button class="close-btn" @click="$emit('close')">Close</button>
      </div>
  
      <div class="client-info">
        <h4>Client Info</h4>
        <p><strong>Name:</strong> {{ dialog.client_name }}</p>
        <p><strong>Phone:</strong> {{ dialog.client_phone || 'N/A' }}</p>
        <p><strong>IP: </strong>
            <a :href="`https://whatismyipaddress.com/ip/${dialog.client_ip}`" target="_blank">
                {{ dialog.client_ip || 'N/A' }}
            </a>
        </p>
      </div>
    </div>
  </template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'
import { fetchMessages, sendMessage } from '../services/messageService'
import { sendEvent } from '../services/socketService'

const props = defineProps({ dialog: Object })
const emit = defineEmits(['close'])

const messages = ref([])
const newMessage = ref('')
const messagesContainer = ref(null)

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function appendMessage(msg) {
  messages.value.push(msg)
  nextTick(scrollToBottom)
}

defineExpose({ appendMessage })

watch(() => props.dialog?.id, async () => {
  if (props.dialog?.id) {
    messages.value = await fetchMessages(props.dialog.id)
    nextTick(scrollToBottom)
  }
})

const submit = () => {
  const text = newMessage.value.trim()
  if (!text) return

  sendEvent('message', {
    room_id: props.dialog.id,
    text,
  })

  messages.value.push({
    id: Date.now(),
    sender: 'operator',
    text,
  })

  newMessage.value = ''
  nextTick(scrollToBottom)
}

const closeDialog = () => {
  sendEvent('dialog_closed', {
    room_id: props.dialog.id,
    info: 'Closed by operator',
  })

  emit('close')
}
</script>

<style scoped>
.dialog-wrapper {
    display: contents;
    height: 100%;
}

.dialog-layout {
  display: flex;
  flex-direction: row;
  height: 100%;
  gap: 1rem;
}

.dialog-panel {
  flex: 3;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 0 6px rgba(0,0,0,0.05);
  color: black;
}

.dialog-panel h3 {
  margin-bottom: 1rem;
  font-size: 1.1rem;
  font-weight: 600;
}

.messages {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.from-operator {
  align-self: flex-end;
  background-color: #d1eaff;
  color: #000;
  padding: 0.6rem 1rem;
  border-radius: 16px;
  max-width: 70%;
  word-break: break-word;
}

.from-client {
  align-self: flex-start;
  background-color: #f0f0f0;
  color: #000;
  padding: 0.6rem 1rem;
  border-radius: 16px;
  max-width: 70%;
  word-break: break-word;
}

.send-form {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.send-form input {
  flex: 1;
  padding: 0.5rem;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  background: white;
  color: black;
  height: 2.5rem;
}

.send-form input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.send-form button {
  height: 2.5rem;
  padding: 0 1rem;
  color: white;
}

.send-form button:hover {
  background-color: #0056b3;
}

.close-btn {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  width: 100%;
  max-width: 200px;
  align-self: center;
}

.close-btn:hover {
  background-color: #c0392b;
}

.client-info {
  flex: 1;
  background: #fff;
  color: black;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 0 6px rgba(0,0,0,0.05);
}

.client-info h4 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.client-info p {
  margin: 0.5rem 0;
  color: #333;
}

</style>
