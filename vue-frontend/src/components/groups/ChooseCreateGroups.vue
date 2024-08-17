<template>
    <div class="container">
        <div class="side">
            <h1 class="less-margin">CHOOSE A GROUP</h1>
            <div class="side-internal-container">
                <h3>GROUPS</h3>
                <div class="groups-list border">
                    <p v-for="group in groups" @click="selectedGroup = group">{{ group.title }}</p>
                </div>
            </div>
            <div class="side-internal-container">
                <h3>DESCRIPTION</h3>
                <div class="description-container border">
                    <p v-if="selectedGroup">{{ selectedGroup.description }}</p>
                </div>
            </div>
            <div class="choose-group-buttons">
                <button v-if="isSelectedGroupMember === false" class="nice-button g" @click="requestToJoinGroup">{{ joinButton }}</button>
                <button v-if="isSelectedGroupMember === true" class="nice-button g" @click="enterGroup">ENTER</button>
                <button v-if="isSelectedGroupMember === true" class="nice-button g" @click="leaveGroup">LEAVE GROUP</button>
            </div>
        </div>
        <div class="separator">
            <div class="line"></div>
            <h2 class="or">OR</h2>
            <div class="line"></div>
        </div>
        <div class="side">
            <h1 class="less-margin">CREATE A NEW GROUP</h1>
            <div class="side-internal-container">
                <h3>NAME</h3>
                <input v-model="groupName" type="text" name="group-name" id="group-name" placeholder="Name me!" required>
            </div>
            <div class="side-internal-container">
                <h3>DESCRIPTION</h3>
                <textarea v-model="groupDescription" name="group-description" id="group-description" cols="30" rows="10" placeholder="Describe me!" required></textarea>
            </div>
            <button class="nice-button g" @click="createGroup"> CREATE GROUP</button>
        </div>
    </div>
</template>

<script setup>
    import { ref, onMounted, watchEffect, watch, computed } from 'vue'
    import { getWebSocketService } from '@/websocket.js'
    import { EventBus } from '@/eventBus';
    import { checkIfFilled } from '@/utils';

    const groups = ref([])
    const userGroups = ref([])

    watchEffect(() => {
        if (EventBus.GetGroups) {
            groups.value = EventBus.GetGroups
        }

        if (EventBus.GetUserGroups) {
            userGroups.value = EventBus.GetUserGroups
        }
    })

    watch(() => EventBus.UpdateUserGroups, () => {
        getWebSocketService().sendMessage("getUserGroups")
    })

    // watch(() => EventBus.UpdateUserGroups)

    onMounted(() => {
        getWebSocketService().sendMessage("getGroups")
        getWebSocketService().sendMessage("getUserGroups")
    })


    const joinButton = ref('JOIN A GROUP')
    const description = ref('')
    const selectedGroup = ref(null)
    const isSelectedGroupMember = ref(false)

    watch([selectedGroup, userGroups], () => {
        console.log("working")

        isSelectedGroupMember.value = false;
        description.value = selectedGroup.value.description;

        if (userGroups.value.length < 1) return;
        
        for (let i = 0; i < userGroups.value.length; i++) {
            if (selectedGroup.value.id === userGroups.value[i].group_id && userGroups.value[i].state == 0) {
                console.log(`SelectedGroup: ${selectedGroup.value.id}, userGroups: ${userGroups.value[i].state}`)
                isSelectedGroupMember.value = false
                joinButton.value = "REQUEST PENDING"
                return
            }
            if (selectedGroup.value.id === userGroups.value[i].group_id && userGroups.value[i].state != 0) {
                console.log(`SelectedGroup: ${selectedGroup.value.id}, userGroups: ${userGroups.value[i].state}`)
                isSelectedGroupMember.value = true
                return
            }
        }

        joinButton.value = "JOIN GROUP"
    })

    const groupName = ref("")
    const groupDescription = ref("")
    function createGroup() {
        const data = {
            title: groupName.value,
            description: groupDescription.value
        };

        if (!checkIfFilled(data)) {
            getWebSocketService().sendMessage("createGroup", data)
        } else {
            console.log("Required fields not filled.")
        }
    }

    function enterGroup() {
        const data = {
            groupID: selectedGroup.value.id
        }

        getWebSocketService().sendMessage("getGroup", data)
    }

    function leaveGroup() {
        const data = {
            groupID: selectedGroup.value.id
        }
        getWebSocketService().sendMessage("deleteGroupMember", data)
        isSelectedGroupMember.value = false;
    }

    function requestToJoinGroup() {
        if (joinButton.value == "REQUEST PENDING") return;

        const data = {
            recipient: selectedGroup.value.creator_id, // recipient in this case is group creator
            groupID: selectedGroup.value.id,
            groupName: selectedGroup.value.title,
            change: "request",
            state: 0
        }

        getWebSocketService().sendMessage("requestToJoinGroup", data)
        joinButton.value = "REQUEST PENDING"
    }
</script>
