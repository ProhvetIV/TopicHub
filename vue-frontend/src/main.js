import './assets/style.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { initializeWebSocketService } from './websocket';

// Initialize the WebSocket service with the router
initializeWebSocketService();
/*
import { webSocketService } from './websocket.js';


// Initialize the WebSocket service with the router
webSocketService(router);*/

const app = createApp(App);

app.use(router);
app.mount('#app');