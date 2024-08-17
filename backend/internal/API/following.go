package api

import (
	"fmt"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
)

func HandleFollow(dataType string, data interface{}) string {
	messageData := data.(map[string]interface{})
	follower, _ := sessions.GetUserNameAndID(messageData["sessionID"].(string))
	beingFollowed := messageData["beingFollowed"].(string)
	currentState := messageData["state"].(string)

	var change string

	if currentState == "follow" {
		if !isProfilePublic(beingFollowed) {
			change = "following"
		} else {
			change = "pending"
		}
		//change = "following"
	} else {
		change = "follow"
	}

	// save change to db
	query := datahandler.PostFollowStatus
	var args []interface{}
	args = append(args, follower, beingFollowed, change)
	_, err := ExecuteQuery(query, args...)
	if err != nil {
		fmt.Println("error executing query", err)
		return currentState
	}

	return change // return change for follower
}

func GetFollowState(follower, beingFollowed string) string {
	var state string
	err := datahandler.Database.QueryRow("SELECT currentstate FROM Followers WHERE follower = ? AND beingfollowed = ?", follower, beingFollowed).Scan(&state)
	if err != nil {
		fmt.Println(err)
	}
	if state == ""{
		state = "follow"
	}
	return state
}

func isProfilePublic(usrn string) bool {
	var state bool
	err := datahandler.Database.QueryRow("SELECT isPublic FROM User WHERE username = ? ", usrn).Scan(&state)
	if err != nil {
		return false // sumkinda err
	}
	return state
}
