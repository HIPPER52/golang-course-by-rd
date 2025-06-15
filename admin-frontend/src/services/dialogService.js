import apiClient from './apiClient';

export async function fetchQueuedDialogs() {
  const res = await apiClient.get('api/operator/dialogs/queued');
  return res.data;
}

export async function fetchActiveDialogs() {
  const res = await apiClient.get('api/operator/dialogs/active');
  return res.data;
}

export async function finishDialogById(dialogId) {
  await apiClient.post(`api/operator/dialogs/${dialogId}/finish`);
}
