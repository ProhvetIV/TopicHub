<script setup>
    import LoginForm from '../components/LoginForm.vue';
    import { initializeWebSocketService } from '@/websocket';

    // Initialize the WebSocket service with the router
    initializeWebSocketService();
</script>

<template>
  <main>
    <LoginForm />
  </main>
</template>
