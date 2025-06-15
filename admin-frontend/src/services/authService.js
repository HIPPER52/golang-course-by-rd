import apiClient from './apiClient';

export async function login(email, password) {
  const res = await apiClient.post('/api/auth/signin', {
    email,
    password,
  });

  return res.data;
}
