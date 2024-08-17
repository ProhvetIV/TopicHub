<template>
  <div class="profilePage">
    <div class="userInfo" v-for="data in user">
      <div class="userHeader">
        <h2 class="userInfoUsername">{{ username }}</h2>
        <!--public/private button for logged-in user-->
        <button class="userProfileButtons" @click="accountStatusF()" v-if="loggedInUser === username">{{ isPublic ? "Private" : "Public" }}</button>
        <div  v-if="loggedInUser !== username" class="otherUserButtons">
        <p class="userProfileButtons"> {{ isPublic ? "Private" : "Public" }}</p>
        <button class="userProfileButtons"  @click="followUnfollow()">{{ followUserButton }}</button></div>
      </div>

      <div class="profilePicture">
        <div v-if="data.image_data !== null" class="profile-image-container">
          <img class="profile-image" :src="imageSrc">
        </div>
        <div v-else-if="imageSrc !== null" class="profile-image-container">
          <img class="profile-image" :src="imageSrc">
        </div>
        <div v-else></div>
        <button class="userProfileButtons" @click="triggerFileInput">Upload File</button>
        <input type="file" ref="fileInput" @change="handleFileUpload" accept=".jpeg, .jpg, .png, .gif, .webp" style="display: none;"/>
      </div>

      <div class="userInfoTable">
        <p>
          <span>Age:</span> <span class="value">{{ data.age }}</span>
        </p>
        <p>
          <span>Gender:</span> <span class="value">{{ data.gender }}</span>
        </p>
        <p>
          <span>First name:</span>
          <span class="value">{{ data.firstname }}</span>
        </p>
        <p>
          <span>Last name:</span> <span class="value">{{ data.lastname }}</span>
        </p>
        <p>
          <span>Email:</span> <span class="value">{{ data.email }}</span>
        </p>
        <p>
          <span>Creation date:</span>
          <span class="value">{{ formatDate(data.creation_date) }}</span>
        </p>
      </div>
    <div class="follow">
      <h2 v-if="followValue" class="followValue">{{ followValue }}</h2>
      <div class="follower-buttons" style="display:flex; justify-content: center;">
        <button class="userProfileButtons" @click="editAboutMe = editAboutMe ? false : true" v-if="followValue === 'About me' && loggedInUser === username">✏️</button>
        <button class="userProfileButtons" @click="aboutMe">About me</button>
        <button class="userProfileButtons" @click="followers">{{ followerList.length }} followers</button>
        <button class="userProfileButtons" @click="following">{{ followingList.length }} following</button>
      </div>
      <div class="list">
        <div v-if="followValue === 'Followers'">
          <div v-for="follower in followerList" :key="follower.id">
            <div class="followListNames">
              <p class="followList">{{ follower.follower }}</p>
            </div>
          </div>
        </div>
        <div v-if="followValue === 'Following'">
          <div v-for="following in followingList" :key="following.id">
            <div class="followListNames">
              <p class="followList">{{ following.beingFollowed }}</p>
            </div>
          </div>
        </div>
        <div v-if="followValue === 'About me'">
            <p v-if="!editAboutMe" class="followList">{{ aboutMeText }}</p>
            <textarea class="editTextarea"  v-if="editAboutMe">{{ aboutMeText }}</textarea>
            <br>
            <button class="userProfileButtons" style="float: right;" @click="saveNewAboutMe()" v-if="editAboutMe">Save</button>
        </div>
      </div>
    </div>
    </div>
    <div class="userPosts">
      <h1 style="text-align: center">User posts</h1>
      <div
        class="postContainer profile-postContainer"
        v-for="post in filteredPosts"
        v-bind:id="post.id">
        <div class="info profile-info">
          <div class="title">{{ post.title }}</div>
          <div class="content">
            <p>{{ post.content }}</p>
            <img v-if="post.image_data !== null" :src="'data:'+getMimeType(post.image_data)+';base64,'+post.image_data" alt="post image" class="post-image">
          </div>
          <div class="opAndDate">
            <div class="creator">{{ post.creator }}</div>
            <div class="creation_date">
              {{ getTimeSince(post.creation_date) }}
            </div>
          </div>
        </div>
        <div class="interactive profile-interactive">
          <div class="likes likes-interactive">
            <button id="upvote" class="upvote" @click="upVoteEvent(post)">
              <i class="fas fa-thumbs-up"></i>
            </button>
            <p class="votes totalVotes">{{ post.upvotes - post.downvotes }}</p>
            <button id="downvote" class="downvote" @click="downVoteEvent(post)">
              <i class="fas fa-thumbs-down"></i>
            </button>
            <button id="viewpost" class="viewpostButton button-12" @click="viewpostEvent(post)">VIEW POST</button>
          </div>
        </div>
      </div>
    </div>
    <ViewPost v-if="selectedPost !== null" :post="selectedPost" @close="selectedPost = null" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject, watch, watchEffect } from "vue";
import { getWebSocketService } from "@/websocket.js";
import { EventBus } from "@/eventBus";
import { formatDate, getMimeType } from "@/utils";
import ViewPost from './ViewPost.vue'
import LoginForm from "./LoginForm.vue";

const posts = ref([]);
const user = ref([]);
const username = ref(null);
let loggedInUser = ref("");
let followValue = ref("");
const followerList = ref([]);
const followingList = ref([]);
const aboutMeText = ref()
let editAboutMe = ref(false)
let isPublic = ref(false);
let followUserButton = ref("follow"); // follow/pending/unfollow button
const imageName = ref(null);
let imageSrc = ref(null);
const fileInput = ref(null);

const triggerFileInput = () => {
  console.log(fileInput.value); // Debugging line
  if (fileInput.value) {
    fileInput.value[0].click();
  } else {
    console.error("fileInput is not defined or click is not a function");
  }
};

// ABOUT ME FUNCTIONS
function aboutMe(){
  getWebSocketService().sendMessage("getUser", { username: username.value});
  followValue.value = 'About me'
  setTimeout(() => {
    aboutMeText.value = user.value[0].aboutMe
  }, 10);
}

function saveNewAboutMe(){
  aboutMeText.value = document.querySelector(".editTextarea").value
  getWebSocketService().sendMessage("addNewAboutMe", {aboutMe: aboutMeText.value });
  editAboutMe.value = false
}
//----

const handleFileUpload = (event) => {
  const file = event.target.files[0];
  if (file) {
    imageSrc.value = URL.createObjectURL(file);
    // Perform actions with the file, e.g., upload it to a server or display it
    /*console.log('Selected file:', file);
    console.log('Selected file parts:', file.name, file.webkitRelativePath);
    getWebSocketService().sendMessage("postProfilePic", {image: file});*/
    const reader = new FileReader();
    reader.onload = () => {
      const base64String = reader.result.split(',')[1]; // Get Base64 string without prefix
      console.log("file.name: ", file.name);
      getWebSocketService().sendMessage("postProfilePic", {
        imageData: base64String,
        imageName: file.name
      }); // sendFileToBackend(base64String, file.name);
    };
    reader.readAsDataURL(file);
  }
};

const usernameProfile = inject("usernameProfile");
username.value = usernameProfile.value;

const fetchProfileData = () => {
  followValue.value = ""
  getWebSocketService().sendMessage("getUser", { username: username.value, isPublic: isPublic.value });
  getWebSocketService().sendMessage("getFollowers", { username: username.value });
  getWebSocketService().sendMessage("getFollowing", { username: username.value });
  getWebSocketService().sendMessage("getFollowState", { beingFollowed: username.value });
  getWebSocketService().sendMessage("getPosts");
};

watch(() => usernameProfile.value, (newUsername) => {
  if (newUsername) {
    username.value = newUsername;
    fetchProfileData();
  }
});

onMounted(() => {
  loggedInUser.value = sessionStorage.getItem("username");
  fetchProfileData();
});

watch(() => EventBus.getPosts, (newPosts) => {
  if (newPosts) {
    console.log("hey! I got new posts!");
    posts.value = newPosts;
  }
});

watch(() => EventBus.getUser, (User) => {
  if (User) {
    user.value = User;
    console.log(user.value[0].image_data)
    isPublic.value = User[0].isPublic;
    followValue.value = "About me"
    aboutMeText.value = user.value[0].aboutMe
    // if (EventBus.image) { 
    //   fetchImage()   
    // } else if (user.value[0].image_id) {
    //   getWebSocketService().sendMessage("getImage", {
    //     imageID: user.value[0].image_id,
    //   });
    // }

    if (user.value[0].image_data !== null) {
      imageSrc.value = 'data:'+getMimeType(user.value[0].image_data)+';base64,'+user.value[0].image_data
    }
  }
});

watch(() => EventBus.getFollowers, (Followers) => {
  if (Followers) {
    followerList.value = Followers;
    console.log(followerList.value)
  }
});

watch(() => EventBus.getFollowing, (Following) =>{
  if (Following){
    followingList.value = Following;
  }
}) 


watchEffect(() => {
  const followState = EventBus.followState;
  if (followState) {
    followUserButton.value = followState;
    console.log("followState updated:", followState);
  }
});

watch(() => EventBus.postreaction, (PostReaction) => {
  if (PostReaction) {
    changeTotalLikes(PostReaction);
  }
});

// watch(() => EventBus.image, (image) => {
//   if (image) { 
//     fetchImage()   
//     //imageName.value = "http://localhost:8080/images/" + image; // "../../../backend/"
   
//     /*fetch(imageName.value)
//     .then(response => {
//       if (!response.ok) {
//         throw new Error(`HTTP error! Status: ${response.status}`);
//       }
//       return response.blob();
//     })
//     .then(blob => {
//       imageSrc.value = URL.createObjectURL(blob);
//       console.log("got new image: ", imageSrc);
//     })
//     .catch(error => {
//       console.error('Error fetching image:', error);
//       // Handle error (e.g., show placeholder image or error message)
//     });*/
//   }
// });

// function fetchImage() {
//   imageName.value = "http://localhost:8080/images/" + EventBus.image
//   fetch(imageName.value)
//     .then(response => {
//       if (!response.ok) {
//         throw new Error(`HTTP error! Status: ${response.status}`);
//       }
//       return response.blob();
//     })
//     .then(blob => {
//       imageSrc.value = URL.createObjectURL(blob);
//       console.log("got new image: ", imageSrc);
//     })
//     .catch(error => {
//       console.error('Error fetching image:', error);
//       // Handle error (e.g., show placeholder image or error message)
//     });
// }

const filteredPosts = computed(() => {
  loggedInUser.value = sessionStorage.getItem("username");
  if (username.value === loggedInUser.value){
    return posts.value.filter((post) => post.creator === username.value);
  }else{
  //display posts to other users
  let showPosts = []; 
    posts.value.forEach((post) =>{
      if (post.postIsPublic === 'public' && post.creator === username.value){
        showPosts.push(post);
      }else if(post.postIsPublic === 'private' && post.creator === username.value){
        if(followerList.value.some((user) => user.follower === loggedInUser.value)){
          showPosts.push(post);
        } 
      }else if (post.postIsPublic === 'almostPrivate' && post.creator === username.value){
        const jsonString = JSON.parse(atob(post.allowedUsers))
        const allowedUsers = ref(jsonString)
        console.log(allowedUsers.value)
        if (allowedUsers.value.some((user) => user === loggedInUser.value && followerList.value.some((user) => user.follower === loggedInUser.value))){
          showPosts.push(post)
        }
      }
    })
    return showPosts;
 }
});

// function for public/private button
function accountStatusF() {
  isPublic.value = isPublic.value ? false : true;
  getWebSocketService().sendMessage("accStatus", {
    isPublic: isPublic.value,
    username: username.value,
  });
}

// follow/unfollow button
function followUnfollow(){
  console.log("follow/unfollow");
  getWebSocketService().sendMessage("postFollower", {
    //Follower: loggedInUser,
    beingFollowed: username.value,
    state: followUserButton.value,
  }) 
} 

// followers list 
function followers() {
  followValue.value = "Followers";
  getWebSocketService().sendMessage("getFollowers", {
    username: username.value
  });
}

// following list
function following() {
  followValue.value = "Following";
   getWebSocketService().sendMessage("getFollowing", {
    username: username.value
  });
}


function getTimeSince(date) {
  var creationDate = new Date(date);
  var currentDate = new Date();

  var timeDifference = currentDate - creationDate;
  var secondsDifference = Math.floor(timeDifference / 1000);

  // Convert seconds to minutes, hours, and days
  var minutes = Math.floor(secondsDifference / 60); //
  var hours = Math.floor(minutes / 60);
  var days = Math.floor(hours / 24);

  if (days != 0) {
    return days + "d " + (hours % 24) + "h";
  } else if (hours != 0) {
    return (hours % 24) + "h " + (minutes % 60) + "min";
  } else {
    return (minutes % 60) + "min";
  }
}

function upVoteEvent(post) {
  handleReactionChangePrePost(post.id, 1);

  if (post.reaction === 1) {
    post.upvotes--;
    post.reaction = 0;
  } else {
    if (post.reaction === 2) {
      post.downvotes--;
    }
    post.upvotes++;
    post.reaction = 1;
  }
}

function downVoteEvent(post) {
  handleReactionChangePrePost(post.id, 2);

  if (post.reaction === 2) {
    post.downvotes--;
    post.reaction = 0;
  } else {
    if (post.reaction === 1) {
      post.upvotes--;
    }
    post.downvotes++;
    post.reaction = 2;
  }
}

function handleReactionChangePrePost(postID, reaction) {
  const thePosts = document.querySelectorAll('[id="' + postID + '"]');

  thePosts.forEach((thePost) => {
    const like = thePost.querySelector("#upvote");
    const dislike = thePost.querySelector("#downvote");

    if (reaction === 2) {
      dislike.classList.toggle("highlightDislike");
      like.classList.remove("highlightLike");
    } else if (reaction === 1) {
      like.classList.toggle("highlightLike");
      dislike.classList.remove("highlightDislike");
    }
  });

  getWebSocketService().sendMessage("postReaction", {
    postID: postID,
    reaction: reaction,
  });
}

const selectedPost = ref(null);

function viewpostEvent(post) {
        selectedPost.value = post;
    }

function changeTotalLikes(data) {
  //const thePost = document.getElementById(data.postID);
  const thePosts = document.querySelectorAll('[id="' + data.postID + '"]');
  console.log(thePosts);

  thePosts.forEach((thePost) => {
    const totalLikes = thePost.querySelector(".totalVotes");
    totalLikes.textContent = parseInt(totalLikes.textContent) + data.change;
  });
}
</script>

<style scoped>
.userInfoUsername {
  margin-left: 10px;
}

.followValue {
  text-align: center;
  margin: 5px 0;
}

.list {
  width: 90%;
  padding: 10px;
  height: 340px;
}

.user-info {
  display: flex;
}

.userHeader {
  display: flex;
  align-items: center;
}

.userProfileButtons {
  font-family: sans-serif;
  color: white;
  padding: 6px 12px;
  font-size: 16px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s ease;
  margin-left: 10px;
  background-color: rgb(54, 54, 54);
}

.userProfileButtons:hover {
  opacity: 0.9;
}

.userInfoTable {
  width: 98%;
  background-color: #f9f9f9;
  border: 3px solid rgba(85, 179, 174, 1);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  font-family: Arial, sans-serif;
}

.userInfoTable p {
  margin: 10px 10%;
  padding: 9px;
  background-color: #ffffff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  justify-content: space-between;
  align-items: center;
}

.userInfoTable p span {
  font-weight: bold;
  color: #333;
}

.userInfoTable p span.value {
  font-weight: normal;
  color: #504e4e;
}

.followList {
  font-size: 20px;
  margin-top: 5px;
}

.followListNames {
  display: flex;
  align-items: center;
}

.followListButton {
  background-color: transparent;
  border: transparent;
  margin-top: 10px;
  margin-left: 10px;
  margin-bottom: 10px;
  cursor: pointer;
}

.follow{
	border: 3px solid  rgba(85, 179, 174, 1);
  background-color: #f9f9f9;
  margin: 0;
  padding: 0;
	border-radius: 8px;
	width: 98%;
	height: 42%;
  overflow-y: auto;
  display: flex;  
  flex-direction: column;
}

.follower-buttons {
  margin-top: 10px;
}

.otherUserButtons{
  display: flex;
  align-items: center;
}

.profile-image {
    width: 200px;
    height: 400px;
    object-fit: contain;
}

.editTextarea{
  border: 1px solid #ccc;
  border-radius: 5px;
  justify-content: center;
  width: 100%;
  height: 150px;
  resize: none;
  font-size: 17px;
}


</style>
