<template>
    <div class="statistics-page">
      <h2>Operator Statistics</h2>
  
      <table class="statistics-table">
        <thead>
          <tr>
            <th>Username</th>
            <th>Email</th>
            <th>Role</th>
            <th>Total Dialogs</th>
            <th>Closed Dialogs</th>
            <th>Avg. Duration (min)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in stats" :key="item.id">
            <td>{{ item.username }}</td>
            <td>{{ item.email }}</td>
            <td>{{ item.role }}</td>
            <td>{{ item.total_dialogs }}</td>
            <td>{{ item.archived_dialogs }}</td>
            <td>{{ item.average_duration_minutes || '—' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { fetchStatistics } from '../services/statisticsService'
  
  const stats = ref([])
  
  async function loadStats() {
    stats.value = await fetchStatistics()
  }
  
  onMounted(loadStats)
  </script>
  
  <style scoped>
  .statistics-page {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .statistics-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1.5rem;
  }
  
  .statistics-table th,
  .statistics-table td {
    padding: 0.75rem;
    border: 1px solid #ccc;
    text-align: left;
  }
  </style>
