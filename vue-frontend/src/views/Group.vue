<template>
    <main>
        <GroupLeftSidebar :group="groupData" @view-members="viewingMembers" @view-event="selectingEvent" @create-event="creatingEvent" :events="events"/>
        <PageContent :posts="posts" :selectedEvent="selectedEvent" :isCreatingEvent="isCreatingEvent" :isViewingMembers="isViewingMembers" :group="groupData" :events="events"/>
        <div id="sidebarContainer">
            <RightSidebar :users="users"/>
            <groupChat :group="groupData"/>
        </div>
    </main>
</template>

<script setup>
    import groupChat from "@/components/groupChat.vue";
    import GroupLeftSidebar from "../components/groups/GroupLeftSidebar.vue"
    import PageContent from "@/components/PageContent.vue";
    import RightSidebar from "@/components/RightSidebar.vue";
    import { useRoute } from "vue-router";
    import { ref, onMounted, watch } from "vue";
    import { EventBus } from "@/eventBus";
    import { getWebSocketService } from "@/websocket";

    const route = useRoute()
    const queryArr = ref([])
    const groupData = ref({})
    const groupDataString = route.query.groupData;
      
    if (groupDataString) {
        try {
            queryArr.value = JSON.parse(groupDataString);
            groupData.value = queryArr.value[0]
        } catch (e) {
            console.error("Failed to parse groupData: ", e)
        }
    }
    
    onMounted(() => {
        const groupPostData = {
            groupID: groupData.value.id
        }
        getWebSocketService().sendMessage('getGroupPosts', groupPostData);
        getWebSocketService().sendMessage('getUsers');

        const groupEventsData = {
            groupID: groupData.value.id
        }
        getWebSocketService().sendMessage("getEvents", groupEventsData)
    })
    
    const posts = ref([])
    const users = ref([])
    
    watch(() => EventBus.UserPostInGroup, (newPost) => {
        posts.value.push(newPost[0])
    })

    watch(() => EventBus.getPosts, (newPosts) => {
        posts.value = newPosts
    })

    watch(() => EventBus.getGroupPosts, (newPosts) => {
        if (newPosts === null) {
            return
        }
        posts.value = newPosts
    })

    watch(() => EventBus.getUsers, (newUsers) => {
        users.value = newUsers
    })

    watch(() => EventBus.ViewMembers, () => {
        isViewingMembers.value = false
    })

    const selectedEvent = ref(null)
    watch(() => EventBus.ViewEvent, () => {
        selectedEvent.value = null
    })

    watch(() => EventBus.CreatingEvent, () => {
        isCreatingEvent.value = false
    })

    watch(() => EventBus.UserPost, (newComment) => {
        posts.value.push(newComment[0])
    })

    const events = ref([])
    watch(() => EventBus.GetEvents, (newEvents) => {
        events.value = newEvents
    })

    watch(() => EventBus.CreateEvent, (newEvent) => {
        if (events.value === null) {
            events.value = newEvent
        } else {
            events.value.push(newEvent[0])
        }
        selectedEvent.value = newEvent[0]
    })

    function selectingEvent(event) {
        selectedEvent.value = event
    }

    const isViewingMembers = ref(false)
    function viewingMembers() {
        isViewingMembers.value = true
    }

    const isCreatingEvent = ref(false)
    function creatingEvent() {
        isCreatingEvent.value = true
    }
</script>
