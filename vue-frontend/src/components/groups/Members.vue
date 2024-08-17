<template>
    <div class="overlay">
        <div class="overlayContent" id="members-overlayContent">
            <div class="members-left">
                <h2>MEMBERS</h2>
                <div class="members-list good-border-background">
                    <p v-for="member in members" class="true-member">{{ member.username }}</p>
                </div>
            </div>
            <div class="members-right">
                <h2>INVITE TO GROUP</h2>
                <div class="invite-list good-border-background">
                    <div v-for="user in inviteUsersList" class="invite-container">
                        <p>{{ user.name }}</p>
                        <button class="nice-button" @click="inviteToGroup(user.name)">INVITE</button>
                    </div>
                </div>
            </div>
        </div>
        <button id="closeOverlay" class="nice-button" @click="closeMembersOverlay">Close</button>
    </div>
</template>

<script setup>
    import { ref, defineProps, onMounted, watchEffect, watch, computed } from 'vue'
    import { EventBus } from '@/eventBus';
    import { getWebSocketService } from '@/websocket';

    const props = defineProps({
        group: {
            type: Object,
            required: true
        }
    })

    onMounted(() => {
        const data = {
            groupID: props.group.id
        }
        getWebSocketService().sendMessage("getGroupMembers", data)

        const followersAndFollowingData = {
            username: sessionStorage.getItem("username")
        }
        getWebSocketService().sendMessage("getFollowers", followersAndFollowingData)
        getWebSocketService().sendMessage("getFollowing", followersAndFollowingData)
    })
    
    const allMembers = ref([])
    const allMembersNames = ref([])
    const members = ref([])
    function sortMembers() {
        members.value.length = 0;
        allMembersNames.value.length = 0;
        allMembers.value.forEach((member) => {
            if (member.state !== 0) {
                members.value.push(member)
            }
            allMembersNames.value.push(member.username)
        })
    }

    // this is for removing duplicate names
    const inviteUsersList = computed(() => {
        const peopleToInvite = []

        for (let i = 0; i < followerList.value.length; i++) {
            if (!allMembersNames.value.includes(followerList.value[i].follower)) {
                peopleToInvite.push(followerList.value[i])
                peopleToInvite[peopleToInvite.length-1]["name"] = followerList.value[i].follower
            }
        }

        for (let i = 0; i < followingList.value.length; i++) {
            if (!allMembersNames.value.includes(followingList.value[i].beingFollowed)) {
                peopleToInvite.push(followingList.value[i])
                peopleToInvite[peopleToInvite.length-1]["name"] = followingList.value[i].beingFollowed
            }
        }

        return peopleToInvite;
    })

    watchEffect(() => {
        if (EventBus.GetGroupMembers) {
            allMembers.value = EventBus.GetGroupMembers
            console.log(allMembers.value)
            
            sortMembers()
            console.log(allMembers.value)
        }
    })

    watch(() => EventBus.UpdateMembers, () => {
        const data = {
            groupID: props.group.id
        }
        console.log("updateing members")
        getWebSocketService().sendMessage("getGroupMembers", data)
    })

    const followerList = ref([])
    watch(() => EventBus.getFollowers, (Followers) => {
        if (Followers) {
            followerList.value = Followers;
            
            sortMembers()
            
            console.log(followerList.value)
            followerList.value = followerList.value.filter((follower) => {
                return !allMembersNames.value.includes(follower.follower)
            })
            console.log(followerList.value)
        }
    });

    const followingList = ref([])
    watch(() => EventBus.getFollowing, (Following) =>{
        if (Following){
            followingList.value = Following;
            
            sortMembers()
            
            console.log(followingList.value)
            followingList.value = followingList.value.filter((following) => {
                console.log(allMembersNames.value, following.beingFollowed)
                return !allMembersNames.value.includes(following.beingFollowed)
            })
            console.log(followingList.value)
        }
    })

    function inviteToGroup(person) {
        const data = {
            recipient: person,
            groupID: props.group.id,
            groupName: props.group.title,
            change: "invite",
            state: 0
        }

        getWebSocketService().sendMessage("requestToJoinGroup", data)
    }


    function closeMembersOverlay() {
        EventBus.ViewMembers = Math.random()
    }
</script>
