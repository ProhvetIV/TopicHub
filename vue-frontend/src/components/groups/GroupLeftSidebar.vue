<template>
    <div class="leftSidebar">
        <div class="sidebar-section part-1">
            <h3 class="sidebar-heading">GROUP</h3>
            <h3>{{ group.title }}</h3>
        </div>
        <div class="sidebar-section part-1">
            <h3 class="sidebar-heading">ABOUT</h3>
            <p>{{ group.description }}</p>
        </div>
        <div class="sidebar-section part-2">
            <h1>EVENTS</h1>
            <div class="events">
                <p v-for="event in events" class="event" @click="$emit('view-event', event)">{{ event.title }}</p>
            </div>
        </div>
        <button class="nice-button sidebar-button" @click="$emit('create-event')">CREATE EVENT</button>
        <button class="nice-button sidebar-button" @click="$emit('view-members')">VIEW MEMBERS</button>
    </div>
</template>

<script setup>
    import { defineProps } from 'vue';

    const props = defineProps({
        group: {
            type: Object,
            required: true
        },
        events: {
            type: Array,
            required: false
        }
    })
</script>