<template>
    <div v-if="group.id != ''" id="rightSidebar">
        <h2 style="text-align: center;">GROUP CHAT</h2>
        <div ref="chatDiv" class="groupChat">
            <span class="chatUsername"></span>
            <div ref="messagesDiv" class="groupMessages" @scroll="messagesDivScrollEvent"></div>
            <div class="groupEmojiList hidden">
                    <div v-for="emoji in emojis" :key="emoji.id">
                        <div class="groupEmojiTable">
                            <button class="groupEmoji" @click="insertEmoij(emoji)">{{emoji}}</button>
                        </div>
                    </div>
                </div>
            <hr class="chatHr">
            <div class="inputSend">
                <input type="text" class="groupchatuserInput" placeholder="Type here..." @keyup.enter="msgButtonEvent">
                <button @click="showEmojiList()" class="emojiButton">😏</button>
                <button id="msgButton" class="chatSumbitButton" @click="msgButtonEvent">SEND</button>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref, defineProps, watchEffect, watch } from 'vue';
    import { getWebSocketService } from '@/websocket';
    import { formatDate, emojis} from '@/utils';
    import { EventBus } from '@/eventBus';
    import { onMounted, inject } from 'vue'
    import { useRouter } from 'vue-router';

    const usernameProfile = inject('usernameProfile');
    const user = ref(null);
    const isPublic = ref(null);
    const router = useRouter(); 
    const followStatus = ref(); // checking follow state for private accounts  
    let count = 0;
    let shouldRoute = ref(false);
    const messages = ref([]);

    const props = defineProps({
        group: {
            type: Object,
            require: true
        }
    })
    
    const messagesDiv = ref(null)
    const chatDiv = ref(null)
    

    let chatUserStatus = ref()
    const chatStatusBool = ref()

    onMounted(() => {
        chatDiv.value.classList.remove("hidden");
        messagesDiv.value.innerHTML = "";
        //username.textContent = user;
        getWebSocketService().sendMessage("getGroupChatMessage", {username: props.group.id});
        count++;
        resetMessagesArray();
    })

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

    watch(() => EventBus.getGroupChatMessage, (messages) => {
        if (messages){
            console.log("getchat: ", EventBus.getGroupChatMessage)
            saveMessages(EventBus.getGroupChatMessage);
            callback();
        }
    })

    watch(() => EventBus.messageGroupMember, (data) => {
        console.log("data: ", data)
        if (data.content !== null) { // if group chat message
            const messageObj = {
                recieverUsername: data.recieverUsername,
               	content: data.content,
                senderUsername: data.senderUsername,
                creation_date: data.creation_date,
            };
            console.log("messageObj: ", messageObj)
            showMessages(messageObj, false);
        } else { // if group post
            
        }
    })

    function showEmojiList(){
        console.log("hi")
        let list = document.querySelector(".groupEmojiList")
        if (list.classList.contains("hidden")){
            list.classList.remove("hidden");
        }else{
            list.classList.add("hidden");
        }
    }

    function insertEmoij(emoji){
       let inputValue =  document.querySelector(".groupchatuserInput");
       inputValue.value += emoji
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
    
    //send message
    function msgButtonEvent() {
        console.log("message send button pressed")
        const regex = /^\s*$/;
        const headerUsername = sessionStorage.getItem("username");
        //const username = document.querySelector(".chatUsername");
        const userInput = document.querySelector(".groupchatuserInput");
    	let userMessage = userInput.value;
        
    	let messageData = {
            senderUsername: headerUsername,
    		recieverUsername: props.group.id,
    		content: userMessage,
            groupID: props.group.id,
            creation_date: Date.now(),
    	};
        console.log("userInput.value: ", userInput.value)
    	// checking empty input
    	userInput.value = "";
    	if (regex.test(userMessage)) {
            console.log("empty string")
    		return;
    	}
        
    	getWebSocketService().sendMessage("postChatMessage", messageData);
    	/*showMessages(
    		{
    			senderUsername: headerUsername,
    			content: userMessage,
    			creation_date: Date.now(),
    		},
    		false
    	);
    	scrollToBottom(".groupMessages");*/
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
    	if (count > 0 || message.recieverUsername === props.group.id) {
            const messagesDiv = document.querySelector(".groupMessages");
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

    		let br = document.createElement("br");
    		let singleMessage = document.createElement("div");
            singleMessage.classList.add("singleMessage")
    		singleMessage.appendChild(senderName);
    		singleMessage.appendChild(spanElement);
    		singleMessage.appendChild(spanDateTime);
    		singleMessage.appendChild(br);
    		if (before) {
    			messagesDiv.insertBefore(singleMessage, messagesDiv.firstChild);
    		} else {
    			messagesDiv.appendChild(singleMessage, messagesDiv.firstChild);
    		}

    		scrollToBottom(".groupMessages");
    	}
    }

    function saveMessages(messagesss) {
        //console.log(messagesss)
        messagesss.forEach((message) => {
            messages.value.push(message);
        })
    	//console.log("messages.value: ", messages.value);
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

    /*function notification(data) {
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
    }*/

</script>

<style scoped>
.emojiButton{
    background-color: transparent;
    border: transparent;
    cursor: pointer;
    font-size: 2em;
}

.groupEmojiList {
    position: relative; /* Changed to relative */
    background-color: rgb(12, 9, 9);
    border: 2px solid #000000;
    border-radius: 5px;
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    z-index: 1000;
}

.groupEmojiTable {
    margin: 5px;
}

.groupEmoji{
    cursor: pointer;
    background-color: transparent;
    border: transparent;
}
</style>