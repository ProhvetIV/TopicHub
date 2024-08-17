<template>
    <div id="rightSidebar">
        <h2 style="text-align: center;">USERS</h2>
        <div class="Users">
            <div class="singleUserButton"
            v-for="user in users"
            :id="user.username"
            :key="user.username"
            >
                <button v-bind:id="user.username" class="LogReg" @click="userProfile(user.username)">{{ user.username }}</button>
                <button v-bind:id="user.username" class="LogReg" @click="chatUser(user.username)">Chat</button>
                <span class="notification circle" :id="user.username+'-notification'">{{ user.unseen > 0 ? user.unseen : "" }}</span>
                <span v-if="user.online" class="dot"></span>
            </div>
        </div>
        <h2 style="text-align: center;">CHAT</h2>
        <div class="Chat hidden">
            <span class="chatUsername"></span>
            <div class="messages" @scroll="messagesDivScrollEvent"></div>
            <div class="emojiList hidden">
                    <div v-for="emoji in emojis" :key="emoji.id">
                        <div class="emojiTable">
                            <button class="emoji" @click="insertEmoij(emoji)">{{emoji}}</button>
                        </div>
                    </div>
                </div>
            <hr class="chatHr">
            <div class="inputSend">
                <input type="text" class="userInput" placeholder="Type here..." @keyup.enter="msgButtonEvent">
                <button @click="showEmojiList()" class="emojiButton">😏</button>
                <button id="msgButton" class="chatSumbitButton" @click="msgButtonEvent">SEND</button>
            </div>
            <button class="close-chat-button" @click="closeChat">close chat</button>
        </div>
    </div>
</template>

<script setup>
    import { ref, defineProps, watchEffect, watch, onMounted } from 'vue';
    import { getWebSocketService } from '@/websocket';
    import { formatDate, emojis } from '@/utils';
    import { EventBus } from '@/eventBus';
    import { inject } from 'vue'
    import { useRouter } from 'vue-router';
    
    const usernameProfile = inject('usernameProfile')
    const user = ref(null)
    const isPublic = ref(null)
    const router = useRouter(); 
    const followStatus = ref() // checking follow state for private accounts  
    const followersList = ref([])
    const followingList = ref([])
    const chatStatusBool = ref()
    let chatUserStatus = ref()
    let count = 0
    let shouldRoute = ref(false)
  

    const props = defineProps({
        users: {
            type: Array,
            required: true
        }
    })

    const messages = ref([])    

    watchEffect(() => {
      if (EventBus.getUser){
        user.value = EventBus.getUser;
        isPublic.value = user.value[0].isPublic;
      }

      if (EventBus.chatStatus){
        chatUserStatus.value = EventBus.chatStatus
        chatStatusBool.value = chatUserStatus.value[0].isPublic;
      }
      
      if (EventBus.followState){ 
        followStatus.value = EventBus.followState
        setTimeout(() =>{
            if (shouldRoute.value) {
            route()
         }
        }, 10)
     }
    })

    onMounted(() =>{
        const headerUsername = sessionStorage.getItem("username");
        getWebSocketService().sendMessage("getFollowers", {username: headerUsername})
        getWebSocketService().sendMessage("getFollowing", {username: headerUsername})
    }) 
    

    watch(() => EventBus.getChatMessage, (messages) => {
        if (messages){
            console.log("getchat")
            saveMessages(EventBus.getChatMessage);
            callback();
        }
    })

    watch(() => EventBus.getFollowers, (followers) =>{
        if (followers){
            followersList.value = followers;
            console.log(followersList.value);
        }
    })

    watch(() => EventBus.getFollowing, (Following) =>{
    if (Following){
        followingList.value = Following;
        console.log("following",followingList.value);

        }
    }) 

    function showEmojiList(){
        let list = document.querySelector(".emojiList")
        if (list.classList.contains("hidden")){
            list.classList.remove("hidden");
        }else{
            list.classList.add("hidden");
        }
    }

    function insertEmoij(emoji){
       let inputValue =  document.querySelector(".userInput");
       inputValue.value += emoji
    }

    // getting user data for profile
    function userProfile(username){
        followStatus.value = ""
        usernameProfile.value = username
        getWebSocketService().sendMessage("getUser",{
            username: username
        })
        getWebSocketService().sendMessage('getFollowState',{
            beingFollowed: username
        })

        shouldRoute.value = true
    }

    function route(){

        if (isPublic.value === true && followStatus.value === "following"){ 
            console.log('Redirecting to /profile');
            router.push('/profile');
        } else if (isPublic.value === false){
            console.log('Redirecting to /profile');
            router.push('/profile');
        } else if (isPublic.value === true && followStatus.value === "follow"  || followStatus.value === "pending"){
            console.log("redireceting to private");
            router.push('/accountprivate');
        }

        shouldRoute.value = false;
        EventBus.followState = null;
    }

    // Event for clicking on a user in the users list. 
    function chatUser(user) {
        const headerUsername = sessionStorage.getItem("username");
        getWebSocketService().sendMessage("chatStatus",{ username: user});
        if (user === headerUsername) {
            console.log("early return in chatUser")
            return;
	    }
            setTimeout(() => {
                if ((followersList.value.some((followerUser) => followerUser.follower === user )  || followingList.value.some((followingUser) => followingUser.beingFollowed === user)) || chatStatusBool.value === false){
            const username = document.querySelector(".chatUsername");
            const messagesDiv = document.querySelector(".messages");
            const userNotification = document.getElementById(`${user}-notification`);
            let inputValue =  document.querySelector(".userInput");
            userNotification.textContent = "";
            const data = {
                user: headerUsername,
                fromUser: user,
            };
            getWebSocketService().sendMessage("rmNotifications", data);
            const chatDiv = document.querySelector(".Chat");
            chatDiv.classList.remove("hidden");
            username.innerHTML = "";
            inputValue.value = ""
            messagesDiv.innerHTML = "";
            username.textContent = user;
            getWebSocketService().sendMessage("getChatMessage", username);
            count++;
            resetMessagesArray();
        }}, 10);
    }

    function closeChat() {
        const chatDiv = document.querySelector(".Chat");
	    chatDiv.classList.add("hidden");
    }
    //send message
    function msgButtonEvent() {
        const regex = /^\s*$/;
        const headerUsername = sessionStorage.getItem("username");
        const username = document.querySelector(".chatUsername");
        const userInput = document.querySelector(".userInput");
    	let userMessage = userInput.value;
        
    	let messageData = {
    		recieverUsername: username.textContent,
    		content: userMessage,
    	};
    	// checking empty input
    	userInput.value = "";
    	if (regex.test(userMessage)) {
    		return;
    	}
    
    	getWebSocketService().sendMessage("postChatMessage", messageData);
    	showMessages(
    		{
    			senderUsername: headerUsername,
    			content: userMessage,
    			creation_date: Date.now(),
    		},
    		false
    	);
    	scrollToBottom(".messages");
    }

    // if somebody writes to you than it goes to last messages
    function scrollToBottom(selector) {
	    const messagesContainer = document.querySelector(selector);
	    messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    // chat message
    function showMessages(message, before = true) {
        let username = document.querySelector(".chatUsername");
        if (username === null){
            return
        }
        localStorage.setItem('chatUsername', username.textContent)
        const headerUsername = sessionStorage.getItem('username')
    	if (
    		(count > 0 && message.senderUsername === headerUsername) ||
    		message.senderUsername === localStorage.getItem('chatUsername') ||
    		message.recieverUsername === localStorage.getItem('chatUsername')
    	) {
            const messagesDiv = document.querySelector(".messages");
    		const messageData = document.createTextNode(message.content);
    		const senderName = document.createTextNode(message.senderUsername + ": ");
    		const messageDateTime = document.createTextNode(formatDate(message.creation_date));
    		// adds space between each text item
    		let spanElement = document.createElement("span");
            spanElement.classList.add("highlighted")
    		spanElement.appendChild(messageData);

    		let spanDateTime = document.createElement("span");
    		spanDateTime.classList.add("dateTime")
    		spanDateTime.appendChild(messageDateTime);

            let spanSender = document.createElement("span")
            spanSender.classList.add("senderSpan")
            spanSender.appendChild(senderName)

    		let br = document.createElement("br");
    		let singleMessage = document.createElement("div");
            singleMessage.classList.add("singleMessage")
            singleMessage.appendChild(spanSender);
    		singleMessage.appendChild(spanElement);
    		singleMessage.appendChild(spanDateTime);
    		singleMessage.appendChild(br);
    		if (before) {
    			messagesDiv.insertBefore(singleMessage, messagesDiv.firstChild);
    		} else {
    			messagesDiv.appendChild(singleMessage, messagesDiv.firstChild);
    		}

    		scrollToBottom(".messages");
    	}
    }

    function saveMessages(messagesss) {
        console.log(messagesss)
        messagesss.forEach((message) => {
            messages.value.push(message);
        })
    	console.log(messages.value);
    }

    function throttle(fn, wait) {
    	let time = Date.now();
    	return function () {
    		if (Date.now() - time >= wait) {
    			fn.apply(this, arguments);
    			time = Date.now();
    		}
    	};
    }

//
     function callback() {
        if (messages.value.length > 0) {
    		messages.value.reverse().forEach((message) => {
    			showMessages(message);
    		});
    		messages.value = [];
    		return;
    	}
    } 

    function resetMessagesArray() {
    	messages.value.length = 0;
    }

    /* const throttledCallback = throttle(callback, 3000);
    //messagesDiv.addEventListener('scroll', () => throttledCallback())
    //messagesDiv.addEventListener('scroll', throttle(callback, 3000))

    function messagesDivScrollEvent() {
        const messagesDiv = document.querySelector(".messages");
        //console.log("blblbl", messagesDiv.scrollTop);
    	if (messagesDiv.scrollTop === 0) {
    		callback(); //throttledCallback();
    	}
    } */
    //messagesDiv.addEventListener("scroll", throttle(callback, 3000));

    function notification(data) {
        const users = document.querySelectorAll(".singleUserButton button");
        const usersContainer = document.querySelector(".Users");

        users.forEach((user) => {
            let singleUserButton = user.closest(".singleUserButton"); // Find the closest parent div with class 'singleUserButton'
            let spanCounter = singleUserButton.querySelector(".notification");
            if (user.textContent === data.senderUsername) {
                let counter = parseInt(spanCounter.textContent || 0);
                counter++;
                spanCounter.textContent = counter;
                spanCounter.classList.add("circle");
                usersContainer.insertBefore(singleUserButton, usersContainer.firstChild);
            }
        });
    }

</script>

<style scoped>
.emojiButton{
    background-color: transparent;
    border: transparent;
    cursor: pointer;
    font-size: 2em;
}

.emojiList {
    position: relative;
    background-color: rgb(12, 9, 9);
    border: 2px solid #000000;
    border-radius: 5px;
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    z-index: 1000;
}

.emojiTable {
    margin: 5px;
}

.emoji{
    cursor: pointer;
    background-color: transparent;
    border: transparent;
}
</style>