<template>
   <div id="leftSidebar" class="leftSidebar-homepage">
        <button @click="handleSort('allPosts')" :class="{'active' : activeFilter === 'allPosts'}" class="leftSideBarButtons">ALL POSTS</button>
        <br>
        <button @click="handleSort('myPosts')" :class="{'active' : activeFilter === 'myPosts'}" class="leftSideBarButtons">MY POSTS</button>
        <br>
        <button @click="handleSort('likedPosts')" :class="{'active' : activeFilter === 'likedPosts'}" class="leftSideBarButtons">LIKED POSTS</button>
        <br>
        <button @click="handleSort('followedUserPosts')" :class="{'active' : activeFilter === 'followedUserPosts'}" class="leftSideBarButtons">FOLLOWED USER POSTS</button>
        <br>
    </div>
</template>

<script setup>

    import {ref} from 'vue';

    const emit = defineEmits(['filter-posts']);
    const activeFilter = ref('');

    function handleSort(filtertype){
        activeFilter.value = filtertype
        emit('filter-posts', filtertype)
    }
</script>

<style>
.leftSideBarButtons {
    color: black;
}

.leftSideBarButtons.active {
    color: grey;

}
</style>