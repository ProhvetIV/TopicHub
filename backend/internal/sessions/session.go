package sessions

import (
	"fmt"
	"log"
	datahandler "social-network/backend/internal/data"
	"time"

	"github.com/google/uuid"
)

// This function will look at all the sessions and see if they are expired.
// If they are expired, they will be deleted.
func CheckForTimedOutSessions(sessionID string) {
	rows, err := datahandler.Database.Query("SELECT session FROM Session")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var session string
		var expires time.Time
		err := rows.Scan(&session, &expires)
		if err != nil {
			fmt.Println("error checking session expiration: ", err)
		}

		if expires.Before(time.Now()) {
			DeleteUserSession(session)
		}
	}

	DeleteUserSession(sessionID)
}

// Delete the session from db and the session cookie.
func DeleteUserSession(sessionID string) {
	_, err := datahandler.Database.Exec("DELETE FROM Session WHERE session = ?", sessionID)

	// if session no longer exists, no need to delete
	if err != nil {
		fmt.Println("Already logged out!")
	}
}

func MakeSession(userID int, username string, response map[string]interface{}) (map[string]interface{}, error) {
	// Create session if password is correct
	sessionToken, _ := uuid.NewRandom()
	expires := time.Now().Add(100 * time.Minute)

	// Try to update the session in db
	var session string
	if err := datahandler.Database.QueryRow("SELECT session FROM Session WHERE user_id = ?", userID).Scan(&session); err != nil {
		datahandler.Database.QueryRow("SELECT username FROM User WHERE id = ? ", userID).Scan(&username)
		// If user not in DB, for example google and git users
		_, err = datahandler.Database.Exec("INSERT INTO Session (user_id, session, expires) VALUES (?,?,?)", userID, sessionToken.String(), expires)
		if err != nil {
			log.Println("Error inserting session into database:", err)
			response["message"] = "Error inserting session into database"
			return response, err
		}

		// If user in DB, update the DB
	} else {
		_, err = datahandler.Database.Exec("UPDATE Session SET session = ? WHERE user_id = ?", sessionToken.String(), userID)
		if err != nil {
			log.Println("Error updating session in database:", err)
			response["message"] = "Error updating session in database"
			return response, err
		}
	}

	response["message"] = "successfully logged in"
	response["username"] = username
	response["sessionID"] = sessionToken.String()
	response["expirationTime"] = expires.Format(time.RFC3339)
	response["userID"] = userID

	return response, nil
}
