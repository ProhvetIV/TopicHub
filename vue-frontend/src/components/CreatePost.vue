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
      </div>
      <button @click="userPostData" id="submit">Submit</button>
    </div>
  </div>
</template>

<script setup>

</script>