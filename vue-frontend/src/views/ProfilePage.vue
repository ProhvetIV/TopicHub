<script setup>
import Profile from '../components/Profile.vue';
import { watchEffect, ref } from 'vue';
import RightSidebar from "@/components/RightSidebar.vue";
import { EventBus } from '@/eventBus';

const users = ref([])

    watchEffect(()=>{
        if (EventBus.getUsers){
            users.value = EventBus.getUsers
        }
    }) 
</script>

<template>
 <main>
    <Profile />
    <div id="sidebarContainer">
        <RightSidebar :users="users"/>
    </div>
 </main>
</template>
