<template>
    <div class="overlay">
        <div class="overlayContent" id="event-view-overlayContent">
            <div class="view-event-container">
                <div class="left-event-description">
                    <h2 class="event-header">EVENT</h2>
                    <h2 class="event-title-content">{{ event.title }}</h2>
                    <h2 class="event-header">DATE</h2>
                    <h2 class="event-title-content">{{ convertDateIntoSthReadable(event.date) }}</h2>
                </div>
                <div class="line"></div>
                <div class="right-event-description">
                    <h2 class="event-header">DESCRIPTION</h2>
                    <h2 class="event-title-content">{{ event.content }}</h2>
                    <div class="attendance-buttons">
                        <button class="nice-button" @click="updateAttendance(1)">GOING</button>
                        <button class="nice-button" @click="updateAttendance(-1)">NOT GOING</button>
                        <button class="nice-button" @click="updateAttendance(0)">MAYBE</button>
                    </div>
                </div>
            </div>
            <h2 class="attendees-header">PEOPLE ATTENDING</h2>
            <div class="attendees">
                <div class="attendee-list">
                    <h4>GOING</h4>
                    <div>
                        <p v-for="g in going" class="attendee">{{ g.username }}</p>
                    </div>
                </div>
                <div class="attendee-list">
                    <h4>NOT GOING</h4>
                    <div>
                        <p v-for="n in notGoing" class="attendee">{{ n.username }}</p>
                    </div>
                </div>
                <div class="attendee-list">
                    <h4>MAYBE</h4>
                    <div>
                        <p v-for="m in maybe" class="attendee">{{ m.username }}</p>
                    </div>
                </div>
            </div>
        </div>
        <button id="closeOverlay" class="nice-button" @click="closeViewEventOverlay">Close</button>
    </div>
</template>

<script setup>
    import { ref, defineProps, onMounted, watch } from 'vue'
    import { EventBus } from '@/eventBus';
    import { getWebSocketService } from '@/websocket';


    const props = defineProps({
        event: {
            type: Object,
            required: true
        }
    })

    const attendees = ref([])

    onMounted(() => {
        const data = {
            eventID: props.event.id,
        }
        getWebSocketService().sendMessage("getEventAttendees", data)
    })

    const going = ref([])
    const notGoing = ref([])
    const maybe = ref([])

    watch(() => EventBus.EventAttendees, (newAttendees) => {
        attendees.value = newAttendees
        console.log("newAttendees: ", attendees.value)
        
        if (attendees.value != null) {
            if (going.value != null) going.value.length = 0;
            if (notGoing.value != null) notGoing.value.length = 0;
            if (maybe.value != null) maybe.value.length = 0;

            for (let i = 0; i < attendees.value.length; i++) {
                if (attendees.value[i].attendance === 1) {
                    going.value.push(attendees.value[i])
                } else if (attendees.value[i].attendance === -1) {
                    notGoing.value.push(attendees.value[i])
                } else {
                    maybe.value.push(attendees.value[i])
                }
            }
        } else {
            if (going.value != null) going.value.length = 0;
            if (notGoing.value != null) notGoing.value.length = 0;
            if (maybe.value != null) maybe.value.length = 0;
        }
    })

    function updateAttendance(attendance) {
        const data = {
            attendance: attendance,
            eventID: props.event.id
        }

        const {isAttending, theAttendance} = checkForUserAttendance()

        if (isAttending && theAttendance == attendance) {
            getWebSocketService().sendMessage("deleteEventAttendee", data)
            return
        }

        if (isAttending) {
            getWebSocketService().sendMessage("updateEventAttendee", data)
            return
        }

        getWebSocketService().sendMessage("addEventAttendee", data)
    }

    watch(() => EventBus.UpdatedAttendance, () => {
        const data = {
            eventID: props.event.id,
        }
        getWebSocketService().sendMessage("getEventAttendees", data)
    })

    function checkForUserAttendance() {
        if (attendees.value == null) return {isAttending: false, theAttendance: 2};

        for (let i = 0; i < attendees.value.length; i++) {
            console.log(`${attendees.value[i].username} vs ${sessionStorage.getItem("username")}`)
            if (attendees.value[i].username == sessionStorage.getItem("username")) {
                return {isAttending: true, theAttendance: attendees.value[i].attendance};
            };
        }

        return {isAttending: false, theAttendance: 2};
    }

    function convertDateIntoSthReadable() {
        const dateString = "2024-06-09T19:14:00Z";

        // Create a new Date object
        const date = new Date(dateString);

        // Format the date
        const readableDate = date.toLocaleString('en-US', {
            weekday: 'long', 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric', 
            hour: 'numeric', 
            minute: 'numeric', 
            second: 'numeric', 
            hour12: true 
        });

        return readableDate;
    }

    function closeViewEventOverlay() {
        EventBus.ViewEvent = Math.random()
    }
</script>
