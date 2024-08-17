<template>
    <div ref="commentWrapper" class="commentWrapper" :id="'commentWrapper_id_'+comment.id">
        <div class="comment">
            <div class="info commentInfo">
                <div class="content commentContent">
                    <p>{{ comment.content }}</p>
                    <img v-if="comment.image_data !== null" :src="'data:'+getMimeType(comment.image_data)+';base64,'+comment.image_data" alt="post_image" class="post-image">
                </div>
                <div class="opAndDate">
                    <div class="creator">OP: {{ comment.creator }}</div>
                    <div class="creation_date">Posted {{ getTimeSince(comment.creation_date) }} ago</div>
                </div>
            </div>
            <div class="interactive commentInteractive">
                <div class="likes">
                    <button class="upvote" id="upvote" :class="{ highlightLike: comment.reaction === 1 }" @click="upVoteEvent(comment)">
                        <i class="fas fa-thumbs-up"></i>
                    </button>
                    <p class="totalVotes votes">{{ comment.upvotes - comment.downvotes }}</p>
                    <button class="downvote" id="downvote" :class="{ highlightLike: comment.reaction === 2 }" @click="downVoteEvent(comment)">
                        <i class="fas fa-thumbs-down"></i>
                    </button>
                </div>
            </div>
        </div>
        <button class="showReply nice-button" ref="showReply" @click="showReplyForm">Reply</button>
        <div class="replySection" ref="replySection">
            <div class="replyContainer" id="replyContainer">
                <div class="replyForm commentReplyForm" id="replyForm">
                    <textarea v-model="commentText" name="commentText" :id="'TextareaNr'+comment.id" class="commentText" placeholder="Comment" required></textarea>
                </div>
                <img v-if="imageSrc" :src="imageSrc" class="preview-img"/>
                <div class="commenting-buttons">
                    <button id="replySubmit" class="nice-button" @click="submitCommentComment">COMMENT</button>
                    <button class="nice-button" @click="triggerFileInput">PICTURE</button>
                    <input type="file" ref="fileInput" @change="showImage" accept=".jpeg, .jpg, .png, .gif, .webp" style="display: none;"/>
                </div>
            </div>
        </div>
        <Comment v-if="commentComments.length > 0" v-for="commentComment in commentComments" :key="commentComment.id" :comment="commentComment" :originalPost="originalPost" :comments="comments" :padding="newPadding"/>
    </div>
</template>

<script setup>
    import { ref, defineProps, onMounted, watchEffect } from 'vue';
    import { downVoteEvent, upVoteEvent } from '@/reactions';
    import { getTimeSince, checkIfFilled, getMimeType } from '@/utils';
    import { getWebSocketService } from '@/websocket';
    import { EventBus } from '@/eventBus';

    const props = defineProps({
        comment: {
            type: Object,
            required: true
        },
        originalPost: {
            type: Number,
            required: true
        },
        comments: {
            type: Array,
            required: true
        },
        padding: {
            type: Number,
            required: true
        }
    })

    const newPadding = props.padding + 1
    const showReply = ref(null)
    const replySection = ref(null)
    const commentComments = ref([])
    const commentText = ref("")
    const commentWrapper = ref(null)

    onMounted(() => {
        const arr = []
        for (let i = 0; i < props.comments.length; i++) {
            if (props.comments[i].parent_post_id === props.comment.id) {
                arr.push(props.comments[i])
            }
        }
        
        commentComments.value = arr
        console.log(props.originalPost)

        // I can make the reddit styling with this
        if (props.padding > 0) {
            commentWrapper.value.style.borderLeft = "2px solid black";
            commentWrapper.value.style.paddingLeft = "40px";
        }
    })

    watchEffect(() => {
        if (EventBus.UserComment) {
            if (EventBus.UserComment[0].parent_post_id == props.comment.id) {
                commentComments.value.push(EventBus.UserComment[0]);
            }
        }
    })

    // Event to show the reply form below the comment.
    function showReplyForm(event) {
        event.target.style.boxShadow = "none"
        replySection.value.classList.toggle("visible");
		showReply.value.classList.toggle("invisible");
        event.target.blur()
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

    // Send comment data to backend and clear input.
    function submitCommentComment() {
        const data = {
            content: commentText.value,
            creator: sessionStorage.getItem("username"),
            parent: props.comment.id,
        };

        if(checkIfFilled(data)){
            return
        }
        
        data["imageData"] = base64String;
        data["imageName"] = imageName;
        imageSrc.value = null

        getWebSocketService().sendMessage("UserComment", data);
        commentText.value = "";
    }
</script>