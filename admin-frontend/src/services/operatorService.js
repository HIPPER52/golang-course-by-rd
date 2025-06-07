import apiClient from './apiClient';

export async function fetchOperators() {
  const res = await apiClient.get('api/admin/users');
  return res.data;
}

export async function createOperator(data) {
  const res = await apiClient.post('api/auth/signup', data);
  return res.data;
}
