<template>
    <header>
        <button v-if="username" @click="homePage" class="nice-button">Home</button>
        <button v-if="username" @click="post" class="nice-button">Post</button>
        <button v-if="!username" @click="login" class="nice-button">Login</button>
        <button v-if="!username" @click="register" class="nice-button">Register</button>
        <button v-if="username" @click="groups" class="nice-button">Groups</button>

        <div v-if="username" class="notifications-wrapper" @mouseleave="hideNotifications">
            <button :class="{'dingo'  : notifications.length > 0}" class="nice-button" @mouseover="showNotifications" @click="logNotifications">Notifications</button>
            <div v-if="showNotificationsDropdown" class="notifications-dropdown">
                <div v-for="notification in notifications" :key="notification.id" class="notification-item">
                    <template v-if="notification.type === 0">
                        <p>{{ notification.actor }} is now following you</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <template v-if="notification.type === 1">
                        <p>{{ notification.actor }} wants to follow you</p>
                        <button @click="acceptNotification(notification.id, notification.actor)" class="nice-button">Accept</button>
                        <button @click="declineNotification(notification.id, notification.actor)" class="nice-button">Decline</button>
                    </template>
                    <template v-if="notification.type === 2">
                        <p>{{ notification.actor }} unfollowed you</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <!-- Type 3 is maybe "accepted from group or denied from group" -->
                    <template v-if="notification.type === 3">
                        <p>Your membership request to join group <b>{{ notification.actor }}</b> has been accepted, enjoy!</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <template v-if="notification.type === 4">
                        <p>Your membership request to join group "{{ notification.actor }}" has been declined, better luck next time!</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <template v-if="notification.type === 5">
                        <p>{{ notification.actor }} wants to join group {{ notification.group_name }}</p>
                        <button @click="handleMembership(notification, 'accept')" class="nice-button">Accept</button>
                        <button @click="handleMembership(notification, 'deny')" class="nice-button">Decline</button>
                    </template>
                    <template v-if="notification.type === 6">
                        <p>{{ notification.actor }} wants you to join {{ notification.group_name }}</p>
                        <button @click="handleMembership(notification, 'join')" class="nice-button">Accept</button>
                        <button @click="handleMembership(notification, 'decline')" class="nice-button">Decline</button>
                    </template>
                    <template v-if="notification.type === 7">
                        <p>{{ notification.actor }} has accepted your invitation to join {{ notification.group_name }}</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <template v-if="notification.type === 8">
                        <p>{{ notification.actor }} has declined your invitation to join {{ notification.group_name }}</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                    <template v-if="notification.type === 9">
                        <p>A new event has just been created in {{ notification.actor }}! Don't miss out!</p>
                        <button @click="removeNotification(notification.id)" class="nice-button">Dismiss</button>
                    </template>
                </div>
            </div>
        </div>

        <button v-if="username" @click="logout" class="nice-button">Log out</button>

        <span v-if="username" @click="profile(username)" class="LogReg nice-button" id="headerUsername">{{ username }}</span> 
    </header>
</template>
  
<script setup>
    import { removeSession, isSessionExpired, sendMessageWhenOpen } from '@/utils.js';
    import { onMounted, ref, watchEffect, inject, watch } from 'vue';
    import { EventBus } from '@/eventBus';
    import { useRouter } from 'vue-router';
    import { getWebSocketService } from '@/websocket.js';
    
    const username = ref(null);//ref(sessionStorage.getItem('username'));//ref(null);
    const notifications = ref([]);
    const showNotificationsDropdown = ref(false);

    const router = useRouter();
    const usernameProfile = inject('usernameProfile')
    const login = () => {
        router.push('/login');
    };

    const register = () => {
        router.push('/register');
    };

    const post = () => {
        router.push('/post');
    };

    const homePage = () => {
        router.push('/home');
        };
    
    const profile = (username) => {
        usernameProfile.value = username
        router.push('/profile')
    };

    const groups = () => {
        router.push('/groups')
    }

    const logout = () => {
        getWebSocketService().sendMessage('logout');
        //removeSession();
        username.value = null;
        router.push('/')
    };

    const showNotifications = () => {
        showNotificationsDropdown.value = true;
    };

    const hideNotifications = () => {
        showNotificationsDropdown.value = false;
    };

    const removeNotification = (id) => {
        notifications.value = notifications.value.filter(notification => notification.id !== id);
        getWebSocketService().sendMessage('removeNotification', { id });
    };

    const acceptNotification = (id, actor) => {
        notifications.value = notifications.value.filter(notification => notification.id !== id);
        getWebSocketService().sendMessage('acceptFollowRequest', { id, actor, username });   
    };

    const declineNotification = (id, actor) => {
        notifications.value = notifications.value.filter(notification => notification.id !== id);
        getWebSocketService().sendMessage('declineFollowRequest', { id, actor, username });       
    };

    const handleMembership = (notificationData, action) => {
        console.log(notificationData)
        const data = {
            groupID: notificationData.group_id,
            groupName: notificationData.group_name,
            username: notificationData.actor,
            change: ""
        }

        if (action == "accept") {
            data.change = action;
            data.state = 1;
            getWebSocketService().sendMessage("acceptGroupMember", data)
        };

        if (action == "deny") {
            data.change = action;
            getWebSocketService().sendMessage("denyGroupMember", data)
        }

        if (action == "join") {
            data.username = sessionStorage.getItem("username");
            data.change = action;
            data.state = 1;
            getWebSocketService().sendMessage("acceptGroupMember", data)
            data.username = notificationData.actor;
        }

        if (action == "decline") {
            data.change = action;
            getWebSocketService().sendMessage("deleteGroupMember", data)
        }

        getWebSocketService().sendMessage("groupMemberRequestNotificationReply", data)
        removeNotification(notificationData.id)
    }

    const logNotifications = () => {
        console.log(notifications.value)
    }

    onMounted(() => {
        isSessionExpired();
        username.value = sessionStorage.getItem('username'); 
        console.log("username: ", username.value);      
        if (username.value) {
            sendMessageWhenOpen('getNotifications');
        }       
    });

    watch(username, (newVal) => {
        if (newVal) {
            if (username.value) {
                router.push('/home');
                sendMessageWhenOpen('getNotifications');
            }
        }
    });
    watch(notifications, (newVal) => {
        console.log("new notifications: ", newVal)
    });
    watchEffect(() => {
        if (EventBus.usrn) {
            username.value = sessionStorage.getItem('username');
        }
        if (EventBus.notifications) {
            notifications.value = EventBus.notifications;
            //notifications.value.push(EventBus.notifications);
            console.log("notifications: ", notifications.value);
            EventBus.notifications = null; // needed for refresh, would rewrite seen notifications in frontend
        }
        /*if (EventBus.notification) {
            //notifications.value.push(EventBus.notification);
            notifications.value.push(EventBus.notification);
            console.log("all notifications: ", notifications)
            //notifications.value = EventBus.notification;
        }*/
    });
</script>

<style scoped>
.notifications-wrapper {
  position: relative;
  display: inline-block;
}

.notifications-dropdown {
  display: block;
  position: absolute;
  background-color: white;
  border: 1px solid #ccc;
  min-width: 200px;
  z-index: 1;
  top: 100%;
  left: 0;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.notification-item {
  padding: 12px;
  border-bottom: 1px solid #eee;
}

.dingo{
    background-color: red;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-item p {
  margin: 0;
}
</style>