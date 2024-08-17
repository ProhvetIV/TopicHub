package api

import (
	"fmt"
	"slices"
	"social-network/backend/internal/sessions"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func GetAndSortUsers(dataType string, data interface{}, clients map[string]*websocket.Conn) []map[string]interface{} {
	dataActual := data.(map[string]interface{})
	username, _ := sessions.GetUserNameAndID(dataActual["sessionID"].(string))
	allUsers := QueryDatabase(dataType, data)
	allMessages := QueryDatabase("getChatMessage", data)

	// Make the user-unseen map
	// Loop through all messages, get unique users, count messages that have not been seen
	// and look for the most recent message
	usersUnseen := make(map[string]map[string]string)

	usernames := make(map[string]bool) // map to check if message belongs to group chat
	for _, user := range allUsers {
		username, ok := user["username"].(string)
		if ok {
			usernames[username] = true
		}
	}

	for _, message := range allMessages {
		var userFromMessage string

		// Get the username from the message.
		// Get the last message
		if message["senderUsername"] != username {
			userFromMessage = message["senderUsername"].(string)
		} else {
			userFromMessage = message["recieverUsername"].(string)
		}
		userMessage := message["creation_date"].(time.Time).Format("2006-01-02 15:04:05")

		fmt.Println("recieverUsername: ", message["recieverUsername"].(string), "the map: ", usernames, "trueorfalse: ", usernames[message["recieverUsername"].(string)])
		// If user is empty string, skip to the next user
		if userFromMessage == "" {
			continue
		} else if !usernames[message["recieverUsername"].(string)] {
			continue
		}

		// Check if the user is already in the map. If not, then add unseen and lastMessage
		if _, ok := usersUnseen[userFromMessage]; !ok {
			usersUnseen[userFromMessage] = make(map[string]string)
			usersUnseen[userFromMessage]["username"] = userFromMessage
			usersUnseen[userFromMessage]["unseen"] = "0"
			usersUnseen[userFromMessage]["lastMessage"] = userMessage
		}

		// Check if the last message is more recent than the one in the map.
		// If the message has not yet been seen, increment the unseen count.
		if userMessage > usersUnseen[userFromMessage]["lastMessage"] {
			usersUnseen[userFromMessage]["lastMessage"] = userMessage
		}
		if message["seen"] == false && message["senderUsername"] != username {
			atoid, _ := strconv.Atoi(usersUnseen[userFromMessage]["unseen"])
			atoid++
			usersUnseen[userFromMessage]["unseen"] = strconv.Itoa(atoid)
		}
	}

	// Convert the map to []usersUnseenStruct
	// Sort the usersUnseenStructs by lastMessage
	type usersUnseenStruct struct {
		username    string
		unseen      string
		lastMessage string
	}
	usersUnseenStructs := make([]usersUnseenStruct, 0)
	for _, user := range usersUnseen {
		usersUnseenStructs = append(usersUnseenStructs, usersUnseenStruct{user["username"], user["unseen"], user["lastMessage"]})
	}
	sort.Slice(usersUnseenStructs, func(i, j int) bool {
		return usersUnseenStructs[i].lastMessage > usersUnseenStructs[j].lastMessage
	})

	// Get online users
	onlineUsers := make([]string, 0)
	for key := range clients {
		onlineUsers = append(onlineUsers, key)
	}

	// Convert the []usersUnseenStruct to []map[string]interface{}
	newAllUsers := make([]map[string]interface{}, 0)
	for _, user := range usersUnseenStructs {
		online := false
		if slices.Contains(onlineUsers, user.username) {
			online = true
		}
		unseenInt, _ := strconv.Atoi(user.unseen)
		newAllUsers = append(newAllUsers, map[string]interface{}{"username": user.username, "unseen": unseenInt, "online": online})
	}

	// Sort all users by username
	justUsers := make([]string, 0)
	for _, i := range allUsers {
		justUsers = append(justUsers, i["username"].(string))
	}
	sort.Strings(justUsers)

	// Add users to the map that have not been messaging with the user []map[string]interface{}
	for _, user := range justUsers {
		if user == username {
			continue
		}

		if _, ok := usersUnseen[user]; ok {
			continue
		}

		online := false
		if slices.Contains(onlineUsers, user) {
			online = true
		}

		newAllUsers = append(newAllUsers, map[string]interface{}{"username": user, "unseen": 0, "online": online})
	}

	return newAllUsers
}
