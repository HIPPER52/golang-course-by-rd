<template>
    <div>
      <h2>Welcome to chat</h2>
      <p>Dialog ID: {{ dialogId }}</p>
    </div>
  </template>
  
  <script setup>
  defineProps({ dialogId: String })
  </script>
