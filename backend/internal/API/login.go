package api

import (
	"errors"
	"fmt"
	datahandler "social-network/backend/internal/data"
	sessions "social-network/backend/internal/sessions"
)

func LogIn(data map[string]interface{}) (map[string]interface{}, error) {
	var response = make(map[string]interface{})

	// Check if user exists
	var userID int
	var dbpassword string
	err := datahandler.Database.QueryRow("SELECT id, password FROM User WHERE username = ? or email = ?", data["username"], data["username"]).Scan(&userID, &dbpassword)
	if err != nil {
		return nil, err
	}

	password, ok := data["password"].(string)
	if !ok {
		fmt.Println("password datatype problems")
	}

	// Check if password is correct
	if !datahandler.CheckPasswordHash(password, dbpassword) {
		response["message"] = "wrong username or password"
		return response, errors.New("wrong username or password")
	}

	return sessions.MakeSession(userID, data["username"].(string), response)
}
