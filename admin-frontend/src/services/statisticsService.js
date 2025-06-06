import apiClient from './apiClient'

export async function fetchStatistics() {
  const res = await apiClient.get('api/admin/stats')
  return res.data
}
