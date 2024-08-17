import router from "./router/index.js";
import { checkUser, createSession, isSessionExpired, removeSession, showMessages, notification } from "./utils.js";
import { EventBus } from "./eventBus.js";
const dataURL = "ws://localhost:8080/data";

class WebSocketService {
  constructor(router) {
    this.socket = null;
    this.router = router;
    //this.eventListeners = new Map();
    this.connect();
  }

  connect() {
    this.socket = new WebSocket(dataURL);

    this.socket.onopen = () => {
      console.log('WebSocket connection established.');
    };

    this.socket.onmessage = this.handleMessage.bind(this);

    this.socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  handleMessage(event) {
    const wsMessage = JSON.parse(event.data);
    //this.$emit(wsMessage.dataType, wsMessage.data);
    const handler = messageHandlers[wsMessage.dataType];
    if (handler) {
      //console.log(wsMessage);
      handler(wsMessage.data, this.router);
    } else {
        console.warn("Unhandled message type:", wsMessage.dataType);
    }
    /*const wsMessage = JSON.parse(event.data);
    const eventListeners = this.eventListeners.get(wsMessage.dataType) || [];
    eventListeners.forEach(callback => callback(wsMessage.data));*/
  }

  sendMessage(dataType, data) {
    if (dataType === "login" || dataType === "UserRegistration") {
      // skip sessionID
    } else {
      if (isSessionExpired() && dataType !== "getPosts") {
        //alert("Your session has expired. Please log in again.");
        removeSession();
        // generate somekind of function to reset back to login page
        this.router.push("/");
      }
      if (data === undefined) {
        data = {
          sessionID: sessionStorage.getItem("sessionID"),
        };
      }
      
      const sessionID = sessionStorage.getItem("sessionID");
      data["sessionID"] = sessionID;
    }
    if (this.socket.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({ dataType, data });
      console.log("Sending message:", message);
      this.socket.send(message);
    } else {
      console.warn("WebSocket connection not open. Message not sent.");
      //initializeWebSocketService();
      //sendMessage(dataType, data);
    }
  }
}

// invalid functons - replace with vue counterparts
const messageHandlers = {
    "UserRegistration": (data, router) => {
        createSession(data);
        getWebSocketService().sendMessage('getNotifications');
        router.push('/home');
    },
    "login": (data, router) => {
        createSession(data);
        getWebSocketService().sendMessage('getNotifications');
        router.push('/home');
    },
    "logout": (_, router) => {
        removeSession();
        router.push('/');
        location.reload()
    },
    "UserPost": () => {
      router.push('/home');
    },
    "addPost": () => {
      router.push('/home');
    },
    "UserPostInGroup": (data) => {
      EventBus.UserPostInGroup = data;
    },

    // "UserComment": (data) => {
    //     appendComment(data);
    // },
    "getPosts": (data) => {
      EventBus.getPosts = data;
      // console.log("getUser: ", data);
    },
    "getGroupPosts": (data) => {
      EventBus.getGroupPosts = data;
    },
    "getComments": (data) => {
        EventBus.getComments = data;
    },
    "getUsers": (data) => {
      EventBus.getUsers = data;
    },
    "getUser": (data) => {
      // console.log("getUser: ", data);
      EventBus.getUser = data;
    },
    "userMessageRecieved": (data) => {
          console.log("userMessageReceived TEST TEST TEST");
         showMessages(data, false);
         notification(data); 
     /// EventBus.userMessageRecieved = data 
     },
    "postChatMessage": (data) => {
        // According to the previous project, nothing to see here.
        EventBus.postChatMessage = data;
    },
    "getChatMessage": (data) => {
        if (data === null) {
            return;
        }

        const messages = []
        data.forEach((userTextObj) => {
          if (!checkUser(userTextObj.senderUsername, userTextObj.recieverUsername)) {
            return;
          } else {
            messages.push(userTextObj);
          }
        });

        EventBus.getChatMessage = messages;
    },
    "getGroupChatMessage": (data) => {
        //console.log("data: ", data);
        if (data === null) {
            return;
        }

        let currentGroupID = EventBus.GetGroup[0].id.toString();
        const messages = []
        data.forEach((userTextObj) => {
          //if (!checkUser(userTextObj.senderUsername, userTextObj.recieverUsername)) {
          if ((userTextObj.senderUsername === currentGroupID || userTextObj.recieverUsername === currentGroupID)) {
            messages.push(userTextObj);
          } else {
            return;
          }
        });
        EventBus.getGroupChatMessage = messages;
    },
    
    // "addPost": (data) => {
    //     const contentBody = document.getElementById("pageContent");
    //     contentBody.appendChild(makePostReviewElement(data[0]));
    // },
    "postReaction": (data) => {
        EventBus.postReaction = data;
    },
    "UserComment": (data) => {
        EventBus.UserComment = data;
    },
    "getGroups": (data) => {
        EventBus.GetGroups = data;
    },
    "getGroup": (data, router) => {
        EventBus.GetGroup = data;
        router.push({ name: 'Group', query: { groupData: JSON.stringify(data) } })
    },
    "createGroup": (data, router) => {
        router.push({ name: 'Group', query: { groupData: JSON.stringify(data) } })
    },
    "getUserGroups": (data) => {
        EventBus.GetUserGroups = data;
    },
    "getGroupMembers": (data) => {
        EventBus.GetGroupMembers = data;
    },
    "deleteGroupMember": () => {
        EventBus.UpdateUserGroups = Math.random();
    },
    "acceptGroupMember": () => {
        EventBus.UpdateMembers = Math.random();
    },
    "denyGroupMember": () => {
        EventBus.UpdateMembers = Math.random();
    },
    "requestToJoinGroup": () => {
        EventBus.UpdateMembers = Math.random();
        EventBus.UpdateUserGroups = Math.random();
    },
    "groupMemberRequestNotificationReply": () => {
      EventBus.UpdateUserGroups = Math.random();
    },
    "getEvents": (data) => {
        EventBus.GetEvents = data
    },
    "createEvent": (data) => {
        EventBus.CreateEvent = data;
    },
    "getEventAttendees": (data) => {
        EventBus.EventAttendees = data;
    },
    "deleteEventAttendee": () => {
        EventBus.UpdatedAttendance = Math.random();
    },
    "addEventAttendee": () => {
        EventBus.UpdatedAttendance = Math.random();
    },
    "updateEventAttendee": () => {
        EventBus.UpdatedAttendance = Math.random();
    },
    "getNotifications": (data) => {
      EventBus.notifications = data;
      EventBus.UpdateUserGroups = Math.random()
      // console.log("getNotifications", data);
    },
    "removeNotification": (data) => {
      // console.log("remove notifications: ", data);
    },
    "notification": (data) => {
      getWebSocketService().sendMessage('getNotifications');
      //EventBus.notification = data;
      //console.log("notification", data);
    },
    "postFollower": (data) => {
      EventBus.followState = data.Change;
    },
    "chatStatus": (data) => {
      EventBus.chatStatus = data;
    },
    "getFollowState": (data) => {
      // console.log("websocket getFollowState data:", data)
      EventBus.followState = data;
    },
    "getImage": (data) => {
      EventBus.image = data;
    },
    /*"sql: no rows in result set": () => { // not registered error to be implemented
      displayBadLoginText();
    },*/
    "updateEventAttendee": () => {
        EventBus.UpdatedAttendance = Math.random();
    },
    "getFollowers": (data) => {
      if (data === null){
        EventBus.getFollowers = [];
      }else{
        EventBus.getFollowers = data;
      }
    },
    "getFollowing": (data) => {
      if (data === null){
        EventBus.getFollowing = [];
      }else{
        EventBus.getFollowing = data;
      }
    }, 
    "messageGroupMember": (data) => {
      EventBus.messageGroupMember = data;
    }, 
    "newGroupPost": (data) => {
      EventBus.gotGroupPost = data;
    },     
    "wrong username or password": (data) => {
      EventBus.badLogin = data;
    },
    "no rows in result set": (data) => {
      EventBus.badLogin = data;
    },
    // "sql: no rows in result set": () => {
    //     displayBadLoginText();
    // },
};

let webSocketServiceInstance = null;

export const initializeWebSocketService = () => {
  if (!webSocketServiceInstance || webSocketServiceInstance.socket.readyState !== WebSocket.OPEN) {
    webSocketServiceInstance = new WebSocketService(router);
    console.log("starting socket", webSocketServiceInstance);
  }
  return webSocketServiceInstance;
};

export const getWebSocketService = () => {
  //console.log("websocket:", webSocketServiceInstance);
  return initializeWebSocketService();

  // webSocketServiceInstance;
};
/*
export const getWebSocketService = () => {
  console.log("websocket:", webSocketServiceInstance);
  if (!webSocketServiceInstance || webSocketServiceInstance.socket.readyState !== WebSocket.OPEN) {
    initializeWebSocketService();
  }
  return webSocketServiceInstance;
};*/
