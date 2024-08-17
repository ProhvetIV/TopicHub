<template>
    <div id="overlay" class="overlay">
        <div class="overlayContent" id="overlayContent">
            <div class="viewPostContainer" v-bind:id="post.id">
                <div class="viewPostOpAndDate">
                    <div class="viewPostCreator">OP: {{ post.creator }}</div>
                    <div class="viewPostCreationDate">Posted: {{ getTimeSince(post.creation_date) }} ago</div>
                </div>
                <div class="viewPostInfo">
                    <div class="viewPostTitle">{{ post.title }}</div>
                    <div class="viewPostContent">
                        <p>{{ post.content }}</p>
                        <div v-if="post.image_data !== null" class="viewpost-image-container">
                            <img :src="'data:'+getMimeType(post.image_data)+';base64,'+post.image_data" alt="post_image" class="post-image">
                        </div>
                    </div>
                </div>
                <div class="viewPostVotebox">
                    <div class="viewPostLikes">
                        <button class="upvote" id="upvote" :class="{ highlightLike: post.reaction === 1 }" @click="upVoteEvent(post)">
                            <i class="fas fa-thumbs-up"></i>
                        </button>
                        <p class="totalVotes votes">{{ post.upvotes - post.downvotes }}</p>
                        <button class="downvote" id="downvote" :class="{ highlightDislike: post.reaction === 2 }" @click="downVoteEvent(post)">
                            <i class="fas fa-thumbs-down"></i>
                        </button>
                    </div>
                </div>
            </div>
            <div class="replyContainer">
                <div class="replyForm" id="replyForm">
                    <textarea v-model="commentText" class="commentText" name="commentText" :id="'textareaNr'+post.id" placeholder="Comment" required></textarea>
                </div>
                <img v-if="imageSrc" :src="imageSrc" class="preview-img"/>
                <div class="commenting-buttons">
                    <button id="replySubmit" class="nice-button" @click="submitComment">COMMENT</button>
                    <button class="nice-button" @click="triggerFileInput">PICTURE</button>
                    <input type="file" ref="fileInput" @change="showImage" accept=".jpeg, .jpg, .png, .gif" style="display: none;"/>
                </div>
            </div>
            <div class="commentsContainer" id="commentsContainer">
                <Comment v-for="comment in localComments" :key="comment.id" :comment="comment" :originalPost="post.id" :comments="comments" :padding="padding" />
            </div>
        </div>
        <button id="closeOverlay" class="nice-button" @click="$emit('close')">Close</button>
    </div>
</template>

<script setup>
    import { getTimeSince, checkIfFilled, getMimeType } from "@/utils";
    import { defineProps, onMounted, ref, watch } from "vue";
    import { downVoteEvent, upVoteEvent } from "@/reactions";
    import { getWebSocketService } from "@/websocket";
    import Comment from "./Comment.vue";
    import { EventBus } from "@/eventBus";

    const props = defineProps({
        post: {
            type: Object,
            required: true
        }
    })

    // For comments, reddit style.
    const padding = 0

    onMounted(() => {
        const data = {
            postID: props.post.id
        }
        getWebSocketService().sendMessage("getComments", data);
    })

    const comments = ref([])
    watch(() => EventBus.getComments, (newComments) => {
        if (newComments) {
            comments.value = newComments;
            commentSieve();
        }
    })

    watch(() => EventBus.UserComment, (newComment) => {
        if (newComment[0] && newComment[0].parent_post_id == props.post.id) {
            localComments.value.push(newComment[0]);
        }
    })

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

    const commentText = ref("")
    function submitComment() {
        const data = {
            content: commentText.value,
            creator: sessionStorage.getItem("username"),
            parent: props.post.id,
        };

        if(checkIfFilled(data)){
            return
        }
        
        data["imageData"] = base64String;
        data["imageName"] = imageName;
        imageSrc.value = null

        getWebSocketService().sendMessage("UserComment", data);
        commentText.value = ""
    }

    const localComments = ref([])
    function commentSieve() {
        const arr = []
        for (let i = 0; i < comments.value.length; i++) {
            if (comments.value[i].parent_post_id === props.post.id) {
                arr.push(comments.value[i])
            }
        }
        localComments.value = arr
    }
</script>
