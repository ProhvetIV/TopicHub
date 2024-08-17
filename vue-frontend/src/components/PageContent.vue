<template>
    <div id="pageContent">
        <div v-if="group" class="post-in-group">
            <label for="postTitle">SHARE A STORY WITH THE WHOLE GROUP</label>
            <input v-model="title" type="text" name="title" id="postTitle" placeholder="Title" required>
            <input v-model="content" type="text" name="content" id="postContent" placeholder="Content" required>
            <button class="nice-button" @click="triggerFileInput">PICTURE</button>
            <input type="file" ref="fileInput" @change="showImage" accept=".jpeg, .jpg, .png, .gif, .webp" style="display: none;"/>
            <img v-if="imageSrc" :src="imageSrc" class="preview-img"/>
            <button @click="userPostDataInGroup" class="nice-button" id="submit">POST</button>
        </div>
        <div class="postContainer"
            v-for="post in posts" 
            v-bind:id="post.id"
        >
            <div class="info">
                <div class="title">{{ post.title }}</div>
                <div class="content">
                    <p>{{ post.content }}</p>
                    <img v-if="post.image_data !== null" :src="'data:'+getMimeType(post.image_data)+';base64,'+post.image_data" alt="post_image" class="post-image">
                </div>
                <div class="opAndDate">
                    <div class="creator">{{ post.creator }}</div>
                    <div class="creation_date">{{ getTimeSince(post.creation_date) }}</div>
                </div>
            </div>
            <div class="interactive">
                <div class="likes">
                    <button id="upvote" class="upvote" :class="{ highlightLike: post.reaction === 1 }" @click="upVoteEvent(post)">
                        <i class="fas fa-thumbs-up"></i>
                    </button>
                    <p class="votes totalVotes">{{ post.upvotes - post.downvotes }}</p>
                    <button id="downvote" class="downvote" :class="{ highlightDislike: post.reaction === 2 }" @click="downVoteEvent(post)">
                        <i class="fas fa-thumbs-down"></i>
                    </button>
                </div>
                <button id="viewpost" class="viewpostButton button-12" @click="viewpostEvent(post)">VIEW POST</button>
            </div>
        </div>
        <ViewPost v-if="selectedPost !== null" :post="selectedPost" @close="selectedPost = null" />
        <Event v-if="selectedEvent" :event="selectedEvent"/>
        <Members v-if="isViewingMembers === true" :group="group"/>
        <CreateEvent v-if="isCreatingEvent === true" :group="group"/>
    </div>
</template>

<script setup>
    import { EventBus } from '@/eventBus';
    import { ref, defineProps, watchEffect } from 'vue';
    import { downVoteEvent, upVoteEvent, changeTotalLikes } from '@/reactions';
    import { checkIfFilled, getTimeSince, getMimeType } from '@/utils';
    import ViewPost from './ViewPost.vue'
    import Event from './groups/Event.vue'
    import Members from './groups/Members.vue'
    import CreateEvent from './groups/CreateEvent.vue'
    import { getWebSocketService } from '@/websocket';
    
    const props = defineProps({
        posts: {
            type: Array,
            required: true
        },
        group: {
            type: Object,
            required: false
        },
        selectedEvent: {
            type: Object,
            required: false
        },
        isViewingMembers: {
            type: Boolean,
            required: false
        },
        isCreatingEvent: {
            type: Boolean,
            required: false
        },
        events: {
            type: Array,
            required: false
        }
    })

    watchEffect(() => {
        if (EventBus.postreaction) {
            changeTotalLikes(EventBus.postreaction);
        }
    })

    const selectedPost = ref(null);
    function viewpostEvent(post) {
        selectedPost.value = post;
    }

    const title = ref('')
    const content = ref('')

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

    function userPostDataInGroup() {
        const data = {
            title: title.value,
            content: content.value,
            parent: "NULL",
            groupID: props.group.id
        };

        console.log(data)

        if (checkIfFilled(data)) {
            console.log("Please fill the required fields in order to post!")
            return
        }

        data["imageData"] = base64String;
        data["imageName"] = imageName;
        imageSrc.value = null
        
        getWebSocketService().sendMessage("UserPostInGroup", data);
        title.value = "";
        content.value = "";
    }

</script>