<template>
    <div class="overlay">
        <div class="overlayContent create-event-overlayContent">
            <div class="create-event">
                <h1>CREATE AN EVENT</h1>
                <h2>WHAT IS THE EVENT CALLED?</h2>
                <input class="event-input" v-model="eventName" type="text" name="event-name" id="event-name" required>
                <h2>WHAT IS IT ABOUT</h2>
                <textarea class="event-input" v-model="eventDescription" name="event-description" id="event-description" cols="30" rows="10" required></textarea>
                <h2>WHEN IS IT HAPPENING?</h2>
                <input class="event-input" v-model="eventDatetime" type="datetime-local" name="event-datetime" id="event-datetime" required>
                <button class="nice-button sidebar-button" @click="createEvent">CREATE EVENT</button>
            </div>
        </div>
        <button id="closeOverlay" class="nice-button" @click="closeViewEventOverlay">Close</button>
    </div>
</template>

<script setup>
    import { ref, defineProps, onMounted } from 'vue'
    import { EventBus } from '@/eventBus';
    import { getWebSocketService } from '@/websocket';
    import { checkIfFilled } from '@/utils';

    const props = defineProps({
        group: {
            type: Object,
            required: true
        }
    })

    function closeViewEventOverlay() {
        EventBus.CreatingEvent = Math.random();
    }

    const eventName = ref('')
    const eventDescription = ref('')
    const eventDatetime = ref()
    function createEvent() {
        const data = {
            title: eventName.value,
            content: eventDescription.value,
            date: eventDatetime.value,
            groupID: props.group.id,
            groupName: props.group.title,
            attendance: 1
        }

        console.log(data)

        if (checkIfFilled(data)) return

        getWebSocketService().sendMessage("createEvent", data)
        closeViewEventOverlay()
    }
</script>

<style scoped>
    .create-event-overlayContent {
	width: 50vw;
	align-items: center;
	justify-content: center;
	flex-direction: row;
}
</style>
