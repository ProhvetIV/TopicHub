package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	api "social-network/backend/internal/API"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		fmt.Println("Origin: ", origin)
		if origin == "http://127.0.0.1:5173" || origin == "http://localhost:5173" {
			return true
		}
		return false
	},
}

var clients = make(map[string]*websocket.Conn)
var mutex sync.Mutex

type WebSocketMessage struct {
	DataType string      `json:"dataType"`
	Data     interface{} `json:"data"`
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	datahandler.OpenDb()
	defer datahandler.Database.Close()
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	//username, _ := sessions.GetUserNameAndID(sessionID) //cookies.GetUserNameAndID(r) //get username

	/*mutex.Lock()
	clients[username] = conn
	mutex.Unlock()*/

	// Handle WebSocket messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Unmarshal the received JSON message
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Println("Error decoding JSON:", err)
			break
		}

		var response WebSocketMessage

		if wsMessage.DataType != "UserRegistration" && wsMessage.DataType != "login" {
			fmt.Println("datatype: ", wsMessage.DataType)
			username, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))
			saveClientSocket(username, conn)
		}

		switch wsMessage.DataType {
		case "getPosts", "getGroupPosts":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "getComments":
			response.DataType = wsMessage.DataType
			allComments := api.GetComments(api.QueryDatabase(wsMessage.DataType, wsMessage.Data), wsMessage.Data, make(map[int64]bool, 0))
			response.Data = allComments
		case "postChatMessage":
			response.DataType = wsMessage.DataType
			response.Data = wsMessage.Data
			data := wsMessage.Data.(map[string]interface{})

			fmt.Println("response.Data: ", response.Data)

			api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
			username, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))

			//wsMessage.Data.(map[string]interface{})["groupID"] = wsMessage.Data.(map[string]interface{})["recieverUsername"]

			_, ok := data["groupID"].(float64)
			if ok {
				allGroupMembers := api.QueryDatabase("getGroupMembers", wsMessage.Data) //???
				for i := 0; i < len(allGroupMembers); i++ {
					member := allGroupMembers[i]
					messageGroupMember(messageType, response, member["username"].(string), "messageGroupMember")
					//newNotificationMessage(messageType, int64(9), member["username"].(string), wsMessage.Data.(map[string]interface{})["groupName"].(string), "", 0, "")
				}
			} else {
				messageMessageRecipient(messageType, data["recieverUsername"].(string), data["content"].(string), username)
				//messageMessageRecipient(messageType, wsMessage, username)
			}

		case "getChatMessage", "getGroupChatMessage":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "rmNotifications":
			api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
			continue
		case "UserRegistration":
			/*response = postRequestChecks(api.RegisterUser(w, r), wsMessage.DataType, w, r)*/
			response.Data, err = api.RegisterUser(wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				fmt.Println("response.Data: ", response.Data)
				response.DataType = wsMessage.DataType
				username := wsMessage.Data.(map[string]interface{})["username"].(string)
				saveClientSocket(username, conn)

				_, ok := wsMessage.Data.(map[string]interface{})["imageData"].(string)
				if ok {
					query := `UPDATE User SET image_id = ? WHERE id = ?;`

					var args []interface{}
					data := wsMessage.Data.(map[string]interface{})

					imageData, err := base64.StdEncoding.DecodeString(data["imageData"].(string))
					if err != nil {
						log.Println(err)
						continue
					}

					imageID := addImageToDb(data["imageName"].(string), imageData) // imageName string, imageData []byte???

					userid := response.Data.(map[string]interface{})["userID"].(int)
					fmt.Println("userid: ", userid)
					args = append(args, imageID, userid)
					api.ExecuteNonQuery(query, args...)
				}
			}
		case "getUsers":
			response.DataType = wsMessage.DataType
			response.Data = api.GetAndSortUsers(wsMessage.DataType, wsMessage.Data, clients)
			if response.Data == nil {
				response.DataType = "Can't get users"
			} else {
				response.DataType = wsMessage.DataType
			}
		case "getUser":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "chatStatus":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "accStatus":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
			continue
		case "login":
			// Create a function that deletes expired sessions
			response.Data, err = api.LogIn(wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
				//fmt.Println(response.DataType)
			} else {
				response.DataType = wsMessage.DataType
				username := wsMessage.Data.(map[string]interface{})["username"].(string)
				saveClientSocket(username, conn)
			}
		case "logout":
			deleteClientSocket(wsMessage.Data.(map[string]interface{}), conn)
			response.Data = api.LogOut(wsMessage.Data.(map[string]interface{}))
			response.DataType = wsMessage.DataType

		case "UserPost":
			postID := 0
			response.Data, postID, err = api.PostData(w, r, wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				response.DataType = wsMessage.DataType

				wsMessage.Data.(map[string]interface{})["postID"] = postID
				othersResponse := WebSocketMessage{
					DataType: "addPost",
					Data:     api.QueryDatabase("getOnePost", wsMessage.Data),
				}
				//othersResponse.Data.(map[string]interface{})["postID"] = postID
				messageAll(messageType, othersResponse)
			}
		case "UserComment":
			_, _, err = api.PostData(w, r, wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				response.DataType = wsMessage.DataType
			}
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "UserPostInGroup":
			_, postID, err := api.PostData(w, r, wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				response.DataType = wsMessage.DataType
			}
			newData := api.ConvertLastInsertIntoInterface(wsMessage.Data, postID, "postID")
			response.Data = api.QueryDatabase(wsMessage.DataType, newData)
			//username, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))

			allGroupMembers := api.QueryDatabase("getGroupMembers", wsMessage.Data) //???
			for i := 0; i < len(allGroupMembers); i++ {
				member := allGroupMembers[i]
				messageGroupMember(messageType, response, member["username"].(string), "UserPostInGroup")
			}
			//messageGroupMember(messageType, recipientUsername, content string, senderUsername string)

		case "postReaction":
			response.DataType = wsMessage.DataType
			type ResponseData struct {
				PostID float64 `json:"postID"`
				Change float64 `json:"change"`
			}

			response.Data = ResponseData{
				PostID: wsMessage.Data.(map[string]interface{})["postID"].(float64), // post that's total like count should be changed
				Change: api.HandleReaction(wsMessage.DataType, wsMessage.Data),      // -2, -1, +1 or +2
			}
			messageAll(messageType, response) // sends to all connected cleints
			continue                          // so that the user who reacted wouldn't get double change
		case "createGroup":
			_, groupID, err := api.CreateGroup(w, r, wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				response.DataType = wsMessage.DataType
			}
			newData := api.ConvertLastInsertIntoInterface(wsMessage.Data, groupID, "groupID")
			if err := api.AddGroupMemberWhenCreatingGroup(wsMessage.Data.(map[string]interface{}), groupID); err != nil {
				response.DataType = err.Error()
			}
			response.Data = api.QueryDatabase("getGroup", newData)
		case "getGroup", "getGroups", "getUserGroups", "getGroupMembers":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "addGroupMember", "deleteGroupMember", "acceptGroupMember", "denyGroupMember":
			response.DataType = wsMessage.DataType
			if err := api.AddDeleteGroupMember(wsMessage.Data.(map[string]interface{}), wsMessage.DataType); err != nil {
				response.DataType = err.Error()
			}
			response.Data = "added or deleted group member"
		case "requestToJoinGroup":
			response.DataType = wsMessage.DataType
			if err := api.AddDeleteGroupMember(wsMessage.Data.(map[string]interface{}), wsMessage.DataType); err != nil {
				response.DataType = err.Error()
			}

			username, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))
			id := map[string]int64{
				"request": 5,
				"invite":  6,
			}

			notification := id[wsMessage.Data.(map[string]interface{})["change"].(string)]

			var recipient string
			if notification == 6 {
				recipient = wsMessage.Data.(map[string]interface{})["recipient"].(string)
			} else {
				recipient = api.GetUsernameById(int(wsMessage.Data.(map[string]interface{})["recipient"].(float64)))
			}

			groupID := int(wsMessage.Data.(map[string]interface{})["groupID"].(float64))
			groupName := wsMessage.Data.(map[string]interface{})["groupName"].(string)

			fmt.Printf("notification: %v, recipient: %v, groupID: %v, groupName: %v\n", notification, recipient, groupID, groupName)
			newNotificationMessage(messageType, notification, recipient, username, "group", groupID, groupName)
			response.Data = "requested to join a group"
		case "createEvent":
			_, eventID, err := api.CreateEvent(wsMessage.Data.(map[string]interface{}))
			if err != nil {
				response.DataType = err.Error()
			} else {
				response.DataType = wsMessage.DataType
			}
			newData := api.ConvertLastInsertIntoInterface(wsMessage.Data, eventID, "eventID")
			if err := api.AddEventAttendeeWhenCreatingEvent(wsMessage.Data.(map[string]interface{}), eventID); err != nil {
				response.DataType = err.Error()
			}
			response.Data = api.QueryDatabase("getEvent", newData)

			// send message to all group members:
			// 1. get all group members
			// 2. send message to them all
			allGroupMembers := api.QueryDatabase("getGroupMembers", wsMessage.Data)
			for i := 0; i < len(allGroupMembers); i++ {
				member := allGroupMembers[i]
				newNotificationMessage(messageType, int64(9), member["username"].(string), wsMessage.Data.(map[string]interface{})["groupName"].(string), "", 0, "")
			}
		case "getEvents", "getEventAttendees":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "addEventAttendee", "deleteEventAttendee", "updateEventAttendee":
			response.DataType = wsMessage.DataType
			if err := api.AddDeleteUpdateEventAttendee(wsMessage.Data.(map[string]interface{}), wsMessage.DataType); err != nil {
				response.DataType = err.Error()
			}
			response.Data = "added, deleted or updated Event_attendees table"
		case "groupMemberRequestNotificationReply":
			// Use newNotification to send the notification, but instead of adding the name of he sender,
			// put group name. This is what the user who requested gets after group creator replies to the
			// request.
			response.DataType = wsMessage.DataType
			id := map[string]int64{
				"accept":  3,
				"deny":    4,
				"join":    7,
				"decline": 8,
			}

			change := wsMessage.Data.(map[string]interface{})["change"].(string)
			notification := id[change]

			recipient := wsMessage.Data.(map[string]interface{})["username"].(string)
			groupID := int(wsMessage.Data.(map[string]interface{})["groupID"].(float64))
			groupName := wsMessage.Data.(map[string]interface{})["groupName"].(string)

			fmt.Println("change: ", change)
			fmt.Println("notification: ", notification)
			fmt.Println("groupID: ", groupID)
			fmt.Println("groupName: ", groupName)

			newNotificationMessage(messageType, notification, recipient, groupName, "group", groupID, groupName)
			response.Data = "notification to member requester sent"
		case "postFollower":
			response.DataType = wsMessage.DataType
			type ResponseData struct {
				Follower          string `json:"Follower"`
				BeingFollowed     string `json:"BeingFollowed"`
				ChangeTo          string `json:"Change"`
				Notification_type int64  `json:"Notification_type"` //NotificationText string `json:"NotificationText"`
			}

			username, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))
			newState := api.HandleFollow(wsMessage.DataType, wsMessage.Data)
			id := map[string]int64{
				"following": 0,
				"pending":   1,
				"follow":    2,
			}
			beingFollowed := wsMessage.Data.(map[string]interface{})["beingFollowed"].(string)

			response.Data = ResponseData{
				Follower:          username,      // Follower
				BeingFollowed:     beingFollowed, // the one who's being followed
				ChangeTo:          newState,      // following, pending, follow
				Notification_type: id[newState],
			}

			/*text := map[string]string{
				"following": username + " is now following you",
				"pending":   username + " requests to follow you",
				"follow":    username + " has unfollowed you",
			}*/

			//response.Data.(map[string]interface{})["Notification_type"] = id[newState]

			//newNotificationMessage(messageType, response)
			newNotificationMessage(messageType, id[newState], beingFollowed, username, "", 0, "")
		case "getFollowers":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "getFollowing":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "getNotifications":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "removeNotification":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "acceptFollowRequest":
			handleFollowRequest(messageType, wsMessage, "following")
			response.Data = api.QueryDatabase("removeNotification", wsMessage.Data)
			continue
		case "declineFollowRequest":
			handleFollowRequest(messageType, wsMessage, "follow")
			response.Data = api.QueryDatabase("removeNotification", wsMessage.Data)
			continue
		case "getFollowState":
			response.DataType = wsMessage.DataType
			follower, _ := sessions.GetUserNameAndID(wsMessage.Data.(map[string]interface{})["sessionID"].(string))
			beingFollowed := wsMessage.Data.(map[string]interface{})["beingFollowed"].(string)
			response.Data = api.GetFollowState(follower, beingFollowed)
			//response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "getImage":
			response.DataType = wsMessage.DataType
			response.Data = getPicture(int64(wsMessage.Data.(map[string]interface{})["imageID"].(float64)))
		case "addNewAboutMe":
			response.DataType = wsMessage.DataType
			response.Data = api.QueryDatabase(wsMessage.DataType, wsMessage.Data)
		case "postProfilePic":
			query := `UPDATE User SET image_id = ? WHERE id = ?;`

			var args []interface{}
			data := wsMessage.Data.(map[string]interface{})

			imageData, err := base64.StdEncoding.DecodeString(data["imageData"].(string))
			if err != nil {
				log.Println(err)
				continue
			}

			imageID := addImageToDb(data["imageName"].(string), imageData) // imageName string, imageData []byte???

			_, userid := sessions.GetUserNameAndID(data["sessionID"].(string))
			args = append(args, imageID, userid)
			api.ExecuteNonQuery(query, args...)
			//continue
			response.DataType = "getImage"
			response.Data = getPicture(imageID)
		default:
			response.DataType = wsMessage.DataType
			response.Data = "Invalid request/not implemented"
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			//break
		}

		// Send the response message back to the client
		if err := conn.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message:", err)
			//break
		}
	}
}

func messageAll(messageType int, response WebSocketMessage) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
	}

	for key, value := range clients {

		if err := value.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message to recipient:", key, err)
		}
	}
}

// func messageMessageRecipient(messageType int, wsMessage WebSocketMessage, senderUsername string) {
func messageMessageRecipient(messageType int, recipientUsername, content string, senderUsername string) {
	var response WebSocketMessage

	//data := wsMessage.Data.(map[string]interface{})
	//recipientUsername, ok := data["recieverUsername"].(string) // get the username from data sent to websocket
	//if !ok {
	//	return
	//}
	response.DataType = "userMessageRecieved"

	type ResponseData struct {
		SenderUsername   string    `json:"senderUsername"`
		Content          string    `json:"content"`
		CreationDate     time.Time `json:"creation_date"`
		RecieverUsername string    `json:"recieverUsername"`
	}

	response.Data = ResponseData{
		SenderUsername: senderUsername,
		Content:        content,
		//Content:          data["content"].(string),
		CreationDate:     time.Now(),
		RecieverUsername: recipientUsername,
	}
	// Check if the recipient is connected
	if recipientConn, ok := clients[recipientUsername]; ok {
		// Marshal the response message to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
		}

		// Send the response message to the recipient WebSocket client
		if err := recipientConn.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message to recipient:", err)
		}
	} else {
		log.Println("Recipient is not connected")
	}
}

func messageGroupMember(messageType int, wsMessage WebSocketMessage, recipientUsername, datatype string) {
	//func messageGroupMember(messageType int, recipientUsername, content string, senderUsername string) {
	var response WebSocketMessage
	//data := wsMessage.Data.(map[string]interface{})
	response.DataType = datatype //
	response.Data = wsMessage.Data
	/*if datatype == "messageGroupMember" {
		type ResponseData struct {
			SenderUsername   string    `json:"senderUsername"`
			Content          string    `json:"content"`
			CreationDate     time.Time `json:"creation_date"`
			RecieverUsername string    `json:"recieverUsername"`
		}

		response.Data = ResponseData{
			SenderUsername:   senderUsername,
			Content:          data["content"].(string),
			CreationDate:     time.Now(),
			RecieverUsername: recipientUsername,
		}
	}*/

	/*response.DataType = "messageGroupMember"

	type ResponseData struct {
		SenderUsername   string    `json:"senderUsername"`
		Content          string    `json:"content"`
		CreationDate     time.Time `json:"creation_date"`
		RecieverUsername string    `json:"recieverUsername"`
	}

	response.Data = ResponseData{
		SenderUsername: senderUsername,
		Content:        content,
		//Content:          data["content"].(string),
		CreationDate:     time.Now(),
		RecieverUsername: recipientUsername,
	}*/
	// Check if the recipient is connected
	if recipientConn, ok := clients[recipientUsername]; ok {
		// Marshal the response message to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
		}

		// Send the response message to the recipient WebSocket client
		if err := recipientConn.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message to recipient:", err)
		}
	} else {
		log.Println("Recipient is not connected")
	}
}

// save to db and message frontend
func newNotificationMessage(messageType int, notificationType int64, recipientUsername, follower string, forWhoWhat string, groupID int, groupName string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := datahandler.PostNotification

	var args []interface{}
	args = append(args, recipientUsername, follower, notificationType, currentTime, groupID, groupName, false)
	fmt.Println("Before insertion to db: ", args)
	lastInsertID, err := api.ExecuteNonQuery(query, args...)
	if err != nil {
		fmt.Println("error executing query", err)
	}

	var response WebSocketMessage
	response.DataType = "notification"
	type ResponseData struct {
		ID                int64  `json:"id"`
		Username          string `json:"user"`
		SenderUsername    string `json:"actor"`
		Notification_type int64  `json:"type"`
		GroupID           int    `json:"group_id"`
		GroupName         string `json:"group_name"`
	}

	if forWhoWhat == "group" {
		fmt.Println("i'm here")
		response.Data = ResponseData{
			ID:                lastInsertID,
			Username:          recipientUsername,
			SenderUsername:    follower,
			Notification_type: notificationType,
			GroupID:           groupID,
			GroupName:         groupName,
		}
	} else {
		response.Data = ResponseData{
			ID:                lastInsertID,
			Username:          recipientUsername,
			SenderUsername:    follower,
			Notification_type: notificationType,
		}
	}
	fmt.Println(response.Data)

	//response.Data.(map[string]interface{})["NotificationID"] = lastInsertID // is int64
	// Check if the recipient is connected
	if recipientConn, ok := clients[recipientUsername]; ok {
		// Marshal the response message to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
		}

		// Send the response message to the recipient WebSocket client
		if err := recipientConn.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message to recipient:", err)
		}
	} else {
		log.Println("Recipient is not connected")
	}
}

func saveClientSocket(username string, conn *websocket.Conn) {
	mutex.Lock()
	clients[username] = conn
	mutex.Unlock()
}

func deleteClientSocket(data map[string]interface{}, conn *websocket.Conn) {

	username, _ := sessions.GetUserNameAndID(data["sessionID"].(string))
	fmt.Println("delete socket of: ", data["sessionID"].(string))

	mutex.Lock()
	delete(clients, username)
	mutex.Unlock()
}

func handleFollowRequest(messageType int, wsMessage WebSocketMessage, state string) {
	var response WebSocketMessage

	data := wsMessage.Data.(map[string]interface{})

	query := datahandler.PostFollowStatus
	var args []interface{}

	follower := data["actor"].(string)
	theFollowed, _ := sessions.GetUserNameAndID(data["sessionID"].(string))
	//theFollowed := data["username"].(string)
	args = append(args, follower, theFollowed, state)
	_, err := api.ExecuteNonQuery(query, args...)
	if err != nil {
		fmt.Println("error executing query", err)
	}

	response.DataType = "postFollower"
	response.Data = struct {
		ChangeTo string `json:"Change"`
	}{
		ChangeTo: state,
	}

	/*type ResponseData struct {
		Follower          string `json:"Follower"`
		BeingFollowed     string `json:"BeingFollowed"`
		ChangeTo          string `json:"Change"`
		Notification_type int64  `json:"Notification_type"` //NotificationText string `json:"NotificationText"`
	}
	response.Data = ResponseData{
				Follower:          username,      // Follower
				BeingFollowed:     beingFollowed, // the one who's being followed
				ChangeTo:          newState,      // following, pending, follow
				Notification_type: id[newState],
			}
	*/
	// change for follower
	if recipientConn, ok := clients[follower]; ok {
		// Marshal the response message to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
		}

		// Send the response message to the recipient WebSocket client
		if err := recipientConn.WriteMessage(messageType, jsonResponse); err != nil {
			log.Println("Error sending message to recipient:", err)
		}
	} else {
		log.Println("Recipient is not connected")
	}
}

func getPicture(imageID int64) string {
	var imageName string
	var imageData []byte

	err := datahandler.Database.QueryRow("SELECT image_name, image_data FROM Images WHERE id = ?", imageID).Scan(&imageName, &imageData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("imageName: ", imageName)
	if fileExists("/internal/data/images/" + imageName) {
		return imageName
	} else {
		if queryAndConvertPicture(imageData, imageName) {
			return imageName
		}
		return "imageName not found"
	}
}

func queryAndConvertPicture(imageData []byte, imageName string) bool {
	// save the image to images folder in backend
	savePath := "../internal/data/images/" + imageName
	err := os.WriteFile(savePath, imageData, 0644)
	if err != nil {
		return false
	}
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func addImageToDb(imageName string, imageData []byte) int64 {
	query := `INSERT INTO Images (image_name, image_data) VALUES (?, ?)`

	var args []interface{}

	args = append(args, imageName, imageData)
	lastInsertID, _ := api.ExecuteNonQuery(query, args...)
	return lastInsertID
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	// Extract image filename from URL path or query parameter
	// Example: /images/{filename}
	filename := r.URL.Path[len("/images/"):]

	// Open the image file
	imagePath := "../internal/data/images/" + filename
	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set Content-Type header based on file type
	contentType := mime.TypeByExtension(filepath.Ext(imagePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)

	// Stream the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to stream image", http.StatusInternalServerError)
		return
	}
}
