<template>
    <div class="chat-page">
      <div class="chat-sidebar">
        <div class="dialog-list">
            <h2>Active Dialogs</h2>
            <ul>
                <li v-for="dialog in activeDialogs" :key="dialog.id">
                <span>{{ dialog.client_name }}</span>
                <div>
                    <button @click="openDialog(dialog)">Open</button>
                    <button @click="finish(dialog.id)">Finish</button>
                </div>
                </li>
            </ul>
        </div>

        <div class="dialog-list">
            <h2>Queued Dialogs</h2>
            <ul>
                <li v-for="dialog in queuedDialogs" :key="dialog.id">
                <span>{{ dialog.client_name }}</span>
                <div>
                    <button @click="take(dialog.id)">Take</button>
                </div>
                </li>
            </ul>
        </div>
      </div>
  
      <div class="chat-main">
        <DialogPanel
          v-if="selectedDialog"
          ref="dialogPanelRef"
          :dialog="selectedDialog"
          @close="selectedDialog = null"
        />
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted, onUnmounted } from 'vue'
  import DialogPanel from '../components/DialogPanel.vue'
  import Sidebar from '../components/Sidebar.vue'
  import { sendEvent } from '../services/socketService'
  import {
    fetchQueuedDialogs,
    fetchActiveDialogs,
    finishDialogById
  } from '../services/dialogService'
  import { subscribe, unsubscribe } from '../services/socketService'
  
  const queuedDialogs = ref([])
  const activeDialogs = ref([])
  const selectedDialog = ref(null)
  const dialogPanelRef = ref(null) // ⏎
  
  const loadDialogs = async () => {
    queuedDialogs.value = await fetchQueuedDialogs()
    activeDialogs.value = await fetchActiveDialogs()
  }
  
  const take = (id) => {
    sendEvent('dialog_taken', {
        room_id: id,
        info: 'Taken by operator',
    })
  }
  
  const finish = (id) => {
    sendEvent('dialog_closed', {
        room_id: id,
        info: 'Closed by operator',
    })

    if (selectedDialog.value?.id === id) {
        selectedDialog.value = null
    }
  }
  
  const openDialog = (dialog) => {
    selectedDialog.value = dialog
  }
  
  async function handleSocketEvent(payload) {
    const { event, data } = payload

    switch (event) {
        case 'dialog_taken':
            queuedDialogs.value = queuedDialogs.value.filter(d => d.id !== data.room_id)

            const alreadyExists = Array.isArray(activeDialogs.value) &&
                      activeDialogs.value.some(d => d.id === data.room_id)
            if (!alreadyExists) {
                const takenDialog = queuedDialogs.value.find(d => d.id === data.room_id)

                if (takenDialog) {
                    activeDialogs.value.push(takenDialog)
                } else {
                    try {
                        const updatedList = await fetchActiveDialogs()
                        activeDialogs.value = updatedList
                    } catch (e) {
                        console.error('Failed to refresh active dialogs', e)
                    }
                }
            }
            break
        case 'dialog_closed':
            activeDialogs.value = activeDialogs.value.filter(d => d.id !== data.room_id)
            break

        case 'message':
            if (selectedDialog.value?.id !== data.room_id) break

            const operatorId = localStorage.getItem('operator_id')

            if (data.sender_id === operatorId) {
                break
            }

            dialogPanelRef.value?.appendMessage?.({
                id: Date.now(),
                sender: data.sender,
                text: data.text,
            })
            break
    }
    }
  
  onMounted(async () => {
    await loadDialogs()
    subscribe('message', handleSocketEvent)
    subscribe('dialog_taken', handleSocketEvent)
    subscribe('dialog_closed', handleSocketEvent)
  })
  
  onUnmounted(() => {
    unsubscribe('message', handleSocketEvent)
    unsubscribe('dialog_taken', handleSocketEvent)
    unsubscribe('dialog_closed', handleSocketEvent)
  })
  </script>


<style scoped>
.chat-view {
  padding: 2rem;
}

.dialog-lists {
  display: flex;
  gap: 2rem;
}

.dialog-list {
  flex: 1;
  background: #f9f9f9;
  padding: 1rem;
  border-radius: 8px;
}

.dialog-list h2 {
  margin-bottom: 1rem;
}

.dialog-list ul {
  list-style: none;
  padding: 0;
}

.dialog-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  margin-bottom: 0.5rem;
  padding: 0.5rem;
  border-radius: 4px;
}

.chat-page {
  display: flex;
  height: 100vh;
  overflow: hidden;
  font-family: 'Inter', sans-serif;
}

.chat-sidebar {
  width: 280px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1rem;
  background: #f5f5f5;
  border-right: 1px solid #ddd;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  overflow: hidden;
}

.chat-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  max-width: 800px;
  margin: 0 auto;
}

.dialog-list {
  margin-bottom: 2rem;
  background: #f9f9f9;
  padding: 1rem;
  border-radius: 8px;
}

.dialog-list h2 {
  margin-bottom: 1rem;
  font-size: 1.2rem;
  color: #333;
}

.dialog-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.dialog-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-bottom: 0.5rem;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  transition: background 0.2s ease;
}

.dialog-list li:hover {
  background: #f1f1f1;
}

.dialog-list span {
  font-weight: 500;
  color: #333;
}

.dialog-list button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 0.4rem 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.dialog-list button + button {
  margin-left: 0.5rem;
}

.dialog-list button:hover {
  background-color: #0056b3;
}
</style>
