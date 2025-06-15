<template>
  <form @submit.prevent="register">
    <input
      v-model="name"
      placeholder="Name"
      required
    >
    <input
      v-model="phone"
      placeholder="Phone"
      required
    >
    <button type="submit">
      Start Chat
    </button>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import { registerClient } from '../services/apiClient';

const emit = defineEmits(['registered']);

const name = ref('');
const phone = ref('');

async function register() {
  try {
    const res = await registerClient(name.value, phone.value);
    emit('registered', res.dialog_id);
  } catch (e) {
    alert('Registration failed');
  }
}
</script>
