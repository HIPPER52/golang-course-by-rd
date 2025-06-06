import apiClient from './apiClient'

export async function fetchMessages(dialogId) {
  const res = await apiClient.get(`api/operator/dialogs/${dialogId}/messages`)
  return res.data
}

export async function sendMessage(dialogId, text) {
  const res = await apiClient.post(`/dialogs/${dialogId}/messages`, { text })
  return res.data
}
