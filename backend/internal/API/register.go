package api

import (
	"fmt"
	"log"
	datahandler "social-network/backend/internal/data"
	sessions "social-network/backend/internal/sessions"

	"golang.org/x/crypto/bcrypt"
)

// This function registers a user into the database
func RegisterUser(data map[string]interface{}) (map[string]interface{}, error) {
	var credentialsTaken bool
	var response = make(map[string]interface{})

	err := datahandler.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE LOWER(email) = ? OR username = ?)", data["email"], data["username"]).Scan(&credentialsTaken)
	if err != nil {
		fmt.Println(err)
		response["message"] = "Error checking credentials"
		return response, err
		//return "Error checking credentials", err
	}

	if credentialsTaken {
		response["message"] = "username/email already taken"
		return response, err
		//return "username/email already taken", err
	}

	password, ok := data["password"].(string)
	if !ok {
		fmt.Println("password datatype problems")
		response["message"] = "password datatype problems"
		return response, err
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	dbResult, err := datahandler.Database.Exec(
		"INSERT INTO User(email, username, password, firstname, lastname, gender, age, isPublic, aboutMe) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?)",
		data["email"], data["username"], pw, data["firstName"], data["lastName"], data["gender"], data["age"], false, data["aboutMe"],
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userID, err := dbResult.LastInsertId()
	if err != nil {
		// Handle error
		log.Println("Error retrieving last insert ID:", err)
		response["message"] = "Error retrieving last insert ID:"
		return nil, err
	}

	// As a return value, create the session and send a message of "successfully registered"
	return sessions.MakeSession(int(userID), data["username"].(string), response)
}
