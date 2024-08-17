package api

import (
	"database/sql"
	"fmt"
	"reflect"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
	"strconv"
)

func getQuery(dataType string, data interface{}) (string, []interface{}) {
	// Make data into map and get sessionID.
	// As a sidequest, check for timed out sessions.
	dataActual := data.(map[string]interface{})
	sessionID := dataActual["sessionID"].(string)
	//sessions.CheckForTimedOutSessions(sessionID)

	// Start getting query
	username, userID := sessions.GetUserNameAndID(sessionID)

	var query string
	var args []interface{}
	switch dataType {
	case "getPosts":
		query = datahandler.FullPost + `WHERE p.parent_post_id IS NULL AND p.group_id IS NULL`
		args = append(args, userID)
		// add extra variables needed like user's ID (the one logged in - unrelated to post otherwise).
	case "getGroupPosts":
		groupID := int(dataActual["groupID"].(float64))
		query = datahandler.FullPost + `WHERE p.parent_post_id IS NULL AND p.group_id = ?`
		args = append(args, userID, groupID)
	case "getComments":
		dataInFloat, ok := dataActual["postID"].(float64)
		if !ok {
			fmt.Println("postID datatype problems")
		}
		parentID := int(dataInFloat)
		query = datahandler.FullPost + `WHERE p.parent_post_id = ?`
		args = append(args, userID, parentID)
	case "UserComment":
		content := dataActual["content"].(string)
		creator := dataActual["creator"].(string)
		parentID := int(dataActual["parent"].(float64))
		query = datahandler.FullPost + `WHERE content = ? AND creator = ? AND p.parent_post_id = ?`
		args = append(args, userID, content, creator, parentID)
	case "UserPostInGroup":
		postID := dataActual["postID"]
		query = datahandler.FullPost + `WHERE p.id = ?`
		args = append(args, userID, postID)
	case "getUsers":
		query = "SELECT username FROM User"
		//args = append(args, userID)
	case "chatStatus":
		realUsername := dataActual["username"].(string)
		query = "SELECT isPublic FROM User WHERE username = ?"
		args = append(args, realUsername)
	case "getUser":
		profileUsername := dataActual["username"].(string)
		query = datahandler.FullUser
		args = append(args, profileUsername)
	case "accStatus":
		profileUsername := dataActual["username"].(string)
		isPublic := dataActual["isPublic"].(bool)
		query = "UPDATE User SET isPublic = ? WHERE username = ?"
		args = append(args, isPublic, profileUsername)
	case "getChatMessage":
		query = datahandler.GetChatMessage
		args = append(args, username, username)
	case "getGroupChatMessage":
		query = datahandler.GetChatMessage
		args = append(args, username, strconv.Itoa(int(dataActual["username"].(float64))))
	case "rmNotifications": //chat notification
		toUser := dataActual["user"].(string)
		fromUser := dataActual["fromUser"].(string)
		//query = "UPDATE Messages SET seen = 1 WHERE id = ?"
		query = `UPDATE Messages SET seen = 1 WHERE from_user_username = ? AND to_user_username = ?`
		args = append(args, fromUser, toUser)
	case "removeNotification": // header notification
		id := dataActual["id"].(float64)
		//query = `DELETE FROM Notifications WHERE id = ?`
		query = `UPDATE Notifications SET seen = 1 WHERE id = ?`
		args = append(args, id)
	case "postChatMessage":
		//fmt.Println("message data", data.(map[string]interface{}))
		messageData, ok := data.(map[string]interface{})
		if !ok {
			fmt.Println("message datatype problems", messageData)
		}

		content, ok := messageData["content"].(string)
		if !ok {
			fmt.Println("content datatype problems")
		}

		recieverUsername, ok := messageData["recieverUsername"].(string)
		if !ok {
			groupid := strconv.Itoa(int(messageData["recieverUsername"].(float64)))
			args = append(args, username, groupid, content)
			//fmt.Println("message reciever datatype problems")
		}

		query = datahandler.PostChatMessage
		args = append(args, username, recieverUsername, content)
	case "getOnePost":
		postID := dataActual["postID"].(int)
		query = datahandler.FullPost + `WHERE p.id = ?`
		args = append(args, 0, postID)
	case "getFollowers":
		query = datahandler.GetFollowers
		username := dataActual["username"].(string)
		args = append(args, username)
	case "getFollowing":
		query = datahandler.GetFollowing
		username := dataActual["username"].(string)
		args = append(args, username)
	case "getNotifications":
		query = datahandler.GetNotifications
		args = append(args, username)
	case "getGroup":
		groupIDtype := reflect.TypeOf(dataActual["groupID"])
		switch groupIDtype.Kind() {
		case reflect.Float64:
			args = append(args, int(dataActual["groupID"].(float64)))
		case reflect.Int:
			args = append(args, dataActual["groupID"].(int))
		}
		query = datahandler.GetGroup
	case "getGroups":
		query = datahandler.GetGroups
	case "getUserGroups":
		query = datahandler.GetUserGroups
		args = append(args, userID)
	case "getGroupMembers":
		groupID := int(dataActual["groupID"].(float64))

		query = datahandler.GetGroupMembers
		args = append(args, groupID)
	case "getEvent":
		eventIDtype := reflect.TypeOf(dataActual["eventID"])
		switch eventIDtype.Kind() {
		case reflect.Float64:
			args = append(args, int(dataActual["eventID"].(float64)))
		case reflect.Int:
			args = append(args, dataActual["eventID"].(int))
		}
		query = datahandler.GetEvent
	case "getEvents":
		groupID := int(dataActual["groupID"].(float64))

		query = datahandler.GetEvents
		args = append(args, groupID)
	case "getEventAttendees":
		eventID := int(dataActual["eventID"].(float64))

		query = datahandler.GetEventAttendees
		args = append(args, eventID)
	case "addNewAboutMe":
		aboutMe := dataActual["aboutMe"].(string)
		query = "UPDATE User SET aboutMe = ? WHERE username = ?"
		args = append(args, aboutMe, username)
	default:
		query = "" // return empty string for other cases
	}

	return query, args
}

func QueryDatabase(dataType string, data interface{}) []map[string]interface{} {
	query, args := getQuery(dataType, data) // basicly figures out whats needed to be done

	if query == "" {
		fmt.Println("empty query")
		//json.NewEncoder(w).Encode(http.StatusInternalServerError) // errorhandling could use some rework here
		return nil
	}

	result, err := ExecuteQuery(query, args...) // does the things needed to be done
	if err != nil {
		fmt.Println("error executing query")
		//json.NewEncoder(w).Encode(http.StatusInternalServerError)
		return nil
	}
	return result
}

func ExecuteQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		// Additional arguments are provided
		rows, err = datahandler.Database.Query(query, args...)
	} else {
		// No additional arguments provided
		rows, err = datahandler.Database.Query(query)
	}
	if err != nil {
		fmt.Println("error here", err)
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create a slice to hold the data
	var data []map[string]interface{}

	// Create a slice of interface{} to represent each row
	values := make([]interface{}, len(columns))
	for i := range columns {
		values[i] = new(interface{})
	}

	// Fetch rows
	for rows.Next() {
		// Scan the values into the interface{} slice
		err := rows.Scan(values...)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Create a map to hold the data for this row
		row := make(map[string]interface{})

		// Convert each interface{} value to the appropriate type and store it in the map
		for i, col := range columns {
			val := *(values[i].(*interface{}))
			row[col] = val
		}

		// Append the map to the data slice
		data = append(data, row)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

// ExecuteNonQuery executes a non-query SQL statement (INSERT, UPDATE, DELETE) and returns the last inserted ID
func ExecuteNonQuery(query string, args ...interface{}) (int64, error) {
	result, err := datahandler.Database.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return lastInsertID, nil
}
