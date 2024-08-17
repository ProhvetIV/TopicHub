<template>
  <main>
    <LeftSidebar @filter-posts="filterPosts" />
    <PageContent :posts="filteredPosts" />
    <div id="sidebarContainer">
      <RightSidebar :users="users"/>

    </div>
  </main>
</template>

<script setup>
import LeftSidebar from "../components/LeftSidebar.vue";
import PageContent from "../components/PageContent.vue";
import RightSidebar from "../components/RightSidebar.vue";

import { ref, computed, onMounted, watchEffect, watch } from "vue";
import { getWebSocketService } from "@/websocket.js";
import { EventBus } from "@/eventBus";

const posts = ref([]);
const users = ref([]);
const followedUsers = ref([]);
const filter = ref("allPosts");
const username = ref(null);

onMounted(() => {
  getData();
});

function getData(){
  getWebSocketService().sendMessage("getPosts");
  getWebSocketService().sendMessage("getUsers");
  username.value = sessionStorage.getItem("username");
  getWebSocketService().sendMessage("getFollowing", {
    username: username.value,
  });
}

watch(() => EventBus.getPosts, (post) =>{
  posts.value = post;
})

watch(() => EventBus.getUsers, (user) =>{
  users.value = user;   
})

watch(() => EventBus.getFollowing, (following) => {
  followedUsers.value = following;
})


watch(() => EventBus.gotGroupPost, (newPost) => {
  posts.value.push(newPost);
})

// updates filter ref
const filterPosts = (criteria) => {
  filter.value = criteria;
  getData();
};

// sorting
const filteredPosts = computed(() => {
  if (posts.value === null) {
    return []
  }

  if (filter.value === "myPosts") {
    return posts.value.filter((post) => post.creator === username.value);
  } else if (filter.value === "likedPosts") {
    return posts.value.filter((post) => post.reaction === 1 || post.reaction === 2);
  } else if (filter.value === "followedUserPosts") {
    let userPost = [];
    posts.value.forEach((post) =>{
      if (post.postIsPublic === "private" || post.postIsPublic === "public"){
        if (followedUsers.value.some((user)=> user.beingFollowed === post.creator) ){
            userPost.push(post);
         }
        }else if (post.postIsPublic === "almostPrivate"){
          const jsonString = JSON.parse(atob(post.allowedUsers));
           const allowedUsers = ref(jsonString);
          if( allowedUsers.value.some((user) => user === username.value) && followedUsers.value.some((user)=> user.beingFollowed === post.creator)){
            userPost.push(post);
          }
        }
  });
    return userPost;
  } else if (filter.value === "allPosts"){
    let allPosts = [];
      if (posts.value === null){
        return;
      }
      posts.value.forEach((post) =>{
        if (post.postIsPublic === "public"){
          allPosts.push(post);
        }else if (post.postIsPublic === "private"){
          if (post.creator === username.value || followedUsers.value.some((user)=> user.beingFollowed === post.creator) ){
            allPosts.push(post);
          }
        } else if (post.postIsPublic === "almostPrivate"){
          const jsonString = JSON.parse(atob(post.allowedUsers));
           const allowedUsers = ref(jsonString);
          if(post.creator === username.value || allowedUsers.value.some((user) => user === username.value) && followedUsers.value.some((user)=> user.beingFollowed === post.creator)){
            allPosts.push(post);
          }
        } 
      });
    return allPosts;
  }
});

</script>
