package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
)

func PostData(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, int, error) {
	_, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	var parentPostID interface{}
	if data["parent"] == "NULL" {
		parentPostID = nil
	} else {
		parentPostID = data["parent"]
	}

	groupID := data["groupID"]
	if groupID == nil || groupID == "NULL" {
		groupID = nil
	} else {
		groupID = int(data["groupID"].(float64))
	}

	var imageID interface{}
	if data["imageName"] != nil {
		imageID = oiAddImageToDb(data["imageName"].(string), data["imageData"].(string))
	} else {
		imageID = nil
	}

	allowedUsers := data["allowedUsers"]
	jsonString, err := json.Marshal(allowedUsers)
	if err != nil {
		log.Fatalf("Error marshalling allowedUsers: %v", err)
	}

	insertString := "INSERT INTO POST(user, title, content, parent_post_id, group_id, image_id, postIsPublic, allowedUsers) VALUES (?,?,?,?,?,?,?,?)"
	result, errData := datahandler.Database.Exec(
		insertString, userID, data["title"], data["content"], parentPostID, groupID, imageID, data["postIsPublic"], jsonString,
	)
	if errData != nil {
		fmt.Println(errData)
		http.Error(w, "Error adding post", http.StatusInternalServerError)
		return "", 0, errData
	}

	lastInsertID, _ := result.LastInsertId()

	return "successfully posted", int(lastInsertID), nil
}

func GetComments(comments []map[string]interface{}, data interface{}, processed map[int64]bool) []map[string]interface{} {
	newComments := make([]map[string]interface{}, 0)

	for _, comment := range comments {
		commentID := comment["id"].(int64)
		if processed[commentID] {
			continue
		}
		processed[commentID] = true

		queryData := data.(map[string]interface{})
		queryData["postID"] = float64(commentID)
		var newData interface{} = queryData

		queryResult := QueryDatabase("getComments", newData)
		if len(queryResult) >= 1 {
			newComments = append(newComments, queryResult...)
		}
	}

	if len(newComments) == 0 {
		return comments
	}

	comments = append(comments, newComments...)
	return GetComments(comments, data, processed)
}

func oiAddImageToDb(imageName string, imageData string) int64 {
	query := `INSERT INTO Images (image_name, image_data) VALUES (?, ?)`

	var args []interface{}
	args = append(args, imageName, imageData)

	result, err := datahandler.Database.Exec(query, args...)
	if err != nil {
		fmt.Printf("error executing query: %v\n", err)
		return 0
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("error getting last insert ID: %v\n", err)
		return 0
	}

	return lastInsertID
}
