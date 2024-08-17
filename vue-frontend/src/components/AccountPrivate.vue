<template>
  <div class="privateAccountText">
    <p>This account is private. Follow to see their profile.</p>
    <button @click="sendFollowRequest">{{followButton}}</button>
  </div>
</template>

<script setup>

import { EventBus } from "@/eventBus";
import { getWebSocketService } from "@/websocket.js";
import { inject, watch, ref, onMounted } from "vue";
import { routeLocationKey } from "vue-router";
import { useRouter } from 'vue-router';

const username = ref(null);
const followButton = ref("follow");
const router = useRouter(); 

const usernameProfile = inject("usernameProfile");

function sendFollowRequest() {
  getWebSocketService().sendMessage("postFollower", {
    beingFollowed: username.value,
    state: followButton.value,
  }) 
}

watch(() => usernameProfile.value, (usernameP) =>{
  if (usernameP){
    username.value = usernameP
    getFollowState();
  }
})

const getFollowState = () => {
  getWebSocketService().sendMessage("getFollowState", {
    beingFollowed: username.value
  })
}

watch(() => EventBus.followState, (followState) =>{
  if (followState){
    followButton.value = followState;
    if (followButton.value === "following"){
       router.push('/profile');
    };
  };
});

// runs only once
onMounted(() => {
  if (usernameProfile.value) {
    username.value = usernameProfile.value;
    getFollowState();
  }
});

</script>

<style scoped>
.privateAccountText {
  margin-top: 25vh;
  border: 2px solid #ccc;
  border-radius: 8px;
  height: 200px;
  width: 300px;
  background-color: #f9f9f9;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.privateAccountText p {
  font-size: 16px;
  color: #333;
  margin-bottom: 20px;
}

.privateAccountText button {
  padding: 10px 20px;
  font-size: 14px;
  color: white;
  background-color: #007bff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.privateAccountText button:hover {
  background-color: #0056b3;
}
</style>
