<template>
  <div class="pageContent">
    <div class="postingWrapper">
      <div class="inputContainer">
        <label for="postTitle">Title</label>
        <br />
        <input type="text" id="postTitle" v-model="title" required />
      </div>
      <div class="inputContainer">
        <label for="postContent">Post</label>
        <br />
        <input type="text" id="postContent" v-model="post" required />
        <br />
        <br />
        <label for="postImage">Image</label>
        <br />
        <button class="nice-button" @click="triggerFileInput">Upload File</button>
        <input type="file" ref="fileInput" @change="showImage" accept=".jpeg, .jpg, .png, .gif, .webp" style="display: none;"/>
        <br />
        <img v-if="imageSrc" :src="imageSrc" class="preview-img"/>
      </div>
      <div class="postStatus">
        <button @click="handlePostStatus('public')" :class="{'active' : activeFilter === 'public'}" class="postStatusButton">Public</button>
        <button @click="handlePostStatus('private')" :class="{'active' : activeFilter === 'private'}" class="postStatusButton">Private</button>
        <button @click="handlePostStatus('almostPrivate')" :class="{'active' : activeFilter === 'almostPrivate'}" class="postStatusButton">Almost private</button>
      </div>
      <div class="almostPrivateUserList hidden">
        <div v-for="user in userList" :key="user.id">
          <div class="listOfUsers">
            <button @click="addUsersToList(user.follower)" :class="['userListButton', { 'selected': selectedUsers.includes(user.follower) }]"> {{ user.follower }}</button>
          </div>
        </div>
      </div>
      <button v-if="activeFilter" @click="userPostData" id="submit">Submit</button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { getWebSocketService } from "@/websocket.js";
import { checkIfFilled } from "@/utils.js"; // @ = shorthand alias for the src directory
import { EventBus } from "@/eventBus";

const title = ref('');
const post = ref('');
const postStatus = ref('');
const loggedInUser = ref('')
const userList = ref([])
const selectedUsers = ref([])
const activeFilter = ref('') // change postStatusButton color when clicked

onMounted(() => {
  loggedInUser.value = sessionStorage.getItem("username");
});

watch(() => EventBus.getFollowers, (followers) =>{
  userList.value = followers
});

function addUsersToList(name){
  if (!selectedUsers.value.includes(name)){
    selectedUsers.value.push(name)
  } else {
    selectedUsers.value = selectedUsers.value.filter(user => user !== name);
  }
}

function handlePostStatus(type){
  activeFilter.value = type
  let list = document.querySelector('.almostPrivateUserList');
  if (type === 'public'  || type === 'private'){
     postStatus.value = type;
     list.classList.add('hidden');
  }else if (type === 'almostPrivate'){
     postStatus.value = type;
     list.classList.remove('hidden');
     getWebSocketService().sendMessage('getFollowers', {username: loggedInUser.value})
  }
}

const fileInput = ref(null)
const triggerFileInput = () => {
    if (fileInput.value) {
        fileInput.value.click();
    } else {
        console.error("fileInput is not defined or click is not a function");
    }
};

const imageSrc = ref(null)
let base64String;
let imageName;
const showImage = (event) => {
  const file = event.target.files[0]
  if (file) {
    imageSrc.value = URL.createObjectURL(file);

    const reader = new FileReader();
    reader.onload = () => {
      base64String = reader.result.split(',')[1];
      imageName = file.name;
    }

    reader.readAsDataURL(file);
  }
}

function userPostData() {
    const data = {
      title: title.value,
      content: post.value,
      parent: "NULL",
      groupID: "NULL",
      postIsPublic: postStatus.value,
      allowedUsers: selectedUsers.value ,
      imageData: base64String,
      imageName: imageName
    };

    getWebSocketService().sendMessage("UserPost", data);
}
</script>

<style>

.postStatusButton{
  background-color: transparent;
  border: transparent;
  cursor: pointer;
}

.postStatusButton.active{
  color: red;
}

.almostPrivateUserList{
  border: 3px solid black;
  height: 12rem;
  width: 10rem;
	border-radius: 5px;
}

.listOfUsers{
  overflow-y: auto;
  text-align: center;
  margin-top: 10px;
}

.userListButton{
  font-size: 17px;
  border-color: transparent;
  background-color: transparent;
  cursor: pointer;
}

.userListButton.selected {
  color: #ff0000; 
}

</style>