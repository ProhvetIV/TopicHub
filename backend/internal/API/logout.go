package api

import (
	"social-network/backend/internal/sessions"
)

func LogOut(data map[string]interface{}) string {
	sessionID := data["sessionID"].(string)

	// Delete the session from db.
	sessions.DeleteUserSession(sessionID)

	return "session deleted from database"
}
