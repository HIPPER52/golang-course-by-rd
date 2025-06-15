import apiClient from './apiClient';

export async function fetchMessages(dialogId) {
  const res = await apiClient.get(`client/dialogs/${dialogId}/messages`);
  return res.data;
}
