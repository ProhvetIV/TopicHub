package sessions

import (
	"fmt"
	datahandler "social-network/backend/internal/data"
)

func GetUserNameAndID(sessionID string) (string, int) {
	var userID int
	if err := datahandler.Database.QueryRow("SELECT user_id FROM Session WHERE session = ?", sessionID).Scan(&userID); err != nil {
		fmt.Println(err)
		return "", 0
	}

	var username string
	if err := datahandler.Database.QueryRow("SELECT username FROM User WHERE id = ?", userID).Scan(&username); err != nil {
		fmt.Println(err)
		return "", 0
	}

	return username, userID
}
