<template>
  <div class="operators-page">
    <h2>Operators</h2>

    <table class="operator-table">
      <thead>
        <tr>
          <th>Username</th>
          <th>Email</th>
          <th>Role</th>
          <th>Created</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="operator in operators" :key="operator.id">
          <td>{{ operator.username }}</td>
          <td>{{ operator.email }}</td>
          <td>{{ operator.role }}</td>
          <td>{{ formatDate(operator.created_at) }}</td>
        </tr>
      </tbody>
    </table>

    <form class="add-user-form" @submit.prevent="submit">
      <input v-model="form.username" placeholder="Username" required />
      <input v-model="form.email" placeholder="Email" required type="email" />
      <input v-model="form.password" placeholder="Password" required type="password" />
      <select v-model="form.role">
        <option value="operator">Operator</option>
        <option value="admin">Admin</option>
      </select>
      <button type="submit">Add User</button>
    </form>
  </div>
</template>
  
<script setup>
import { ref, onMounted } from 'vue'
import { fetchOperators, createOperator } from '../services/operatorService'
import { formatDate } from '../utils/formatDate'

const operators = ref([])
const form = ref({
  username: '',
  email: '',
  password: '',
  role: 'operator'
})

async function loadOperators() {
  operators.value = await fetchOperators()
}

async function submit() {
  const dto = {
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      role: form.value.role
  }

  try {
      await createOperator(dto)
      await loadOperators()

      form.value.username = ''
      form.value.email = ''
      form.value.password = ''
      form.value.role = 'operator'
  } catch (err) {
      const errorMessage = err?.response?.data?.error || 'Failed to create operator'
      alert(`Error: ${errorMessage}`)
  }
}

onMounted(loadOperators)
</script>
  
<style scoped>
.operators-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.operator-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 2rem;
}

.operator-table th, .operator-table td {
  padding: 0.75rem;
  border: 1px solid #ddd;
  text-align: left;
}

.add-user-form {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.add-user-form input,
.add-user-form select,
.add-user-form button {
  padding: 0.5rem;
  font-size: 1rem;
}

.add-user-form button {
  background: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

.add-user-form button:hover {
  background: #0056b3;
}
</style>
