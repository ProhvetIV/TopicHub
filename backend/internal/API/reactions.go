package api

import (
	"database/sql"
	"fmt"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
)

func HandleReaction(dataType string, data interface{}) float64 {
	messageData := data.(map[string]interface{})
	postID := messageData["postID"].(float64)
	reaction := messageData["reaction"].(float64)
	sessionID := messageData["sessionID"].(string)

	_, userID := sessions.GetUserNameAndID(sessionID)

	query := datahandler.PostReaction // add new reaction or change the reaction entry
	var change, newReaction float64

	previousReaction := haveReacted(userID, postID)

	if previousReaction != reaction {
		newReaction = reaction
	} else {
		newReaction = 0
	}

	if previousReaction == reaction {
		change = -1
	} else if previousReaction == 0 {
		change = 1
	} else {
		change = 2
	}
	if reaction == 2 {
		change *= -1
	}

	var args []interface{}
	args = append(args, userID, postID, newReaction)
	_, err := ExecuteQuery(query, args...) // does the things needed to be done
	if err != nil {
		fmt.Println("error executing query", err)
		//json.NewEncoder(w).Encode(http.StatusInternalServerError)
		return 0
	}

	return change
}

func haveReacted(userID int, postID float64) float64 {
	var existingReaction float64
	err := datahandler.Database.QueryRow("SELECT reaction FROM Postreaction WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingReaction)
	if err != nil {
		if err == sql.ErrNoRows { // no previous vote
			return 0
		} else { // somethings wrong
			fmt.Println(err)
		}
	}

	return existingReaction
}
