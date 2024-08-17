import { getWebSocketService } from "./websocket";

export function upVoteEvent(post) {
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

export function downVoteEvent(post) {
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

export function handleReactionChangePrePost(postID, reaction) {
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

    getWebSocketService().sendMessage("postReaction", { postID: postID, reaction: reaction });
}

export function changeTotalLikes(data) {
    //const thePost = document.getElementById(data.postID);
    const thePosts = document.querySelectorAll('[id="' + data.postID + '"]');
    console.log(thePosts);

    thePosts.forEach((thePost) => {
        const totalLikes = thePost.querySelector(".totalVotes");
        totalLikes.textContent = parseInt(totalLikes.textContent) + data.change;
    });
}