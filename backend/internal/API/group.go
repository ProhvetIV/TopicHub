package api

import (
	"fmt"
	"net/http"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
)

func CreateGroup(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (string, int, error) {
	_, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	result, err := datahandler.Database.Exec(
		datahandler.CreateGroup, data["title"], data["description"], userID,
	)
	if err != nil {
		fmt.Println("error inserting group to database: ", err)
		http.Error(w, "Error adding post", http.StatusInternalServerError)
		return "", 0, err
	}

	lastInsert, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error retrieving group id after posting: ", err)
		http.Error(w, "Error adding post", http.StatusInternalServerError)
		return "", 0, err
	}

	return "Group successfully created!", int(lastInsert), nil
}

func ConvertLastInsertIntoInterface(data interface{}, lastInsert int, dataName string) interface{} {
	dataActual := data.(map[string]interface{})
	dataActual[dataName] = lastInsert

	var itf interface{}
	itf = dataActual

	return itf
}

func AddGroupMemberWhenCreatingGroup(data map[string]interface{}, groupID int) error {
	username, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	_, err := datahandler.Database.Exec(
		datahandler.AddGroupMember, groupID, userID, username, data["state"],
	)
	if err != nil {
		fmt.Println("error inserting group member to database: ", err)
		return err
	}

	return nil
}

func AddDeleteGroupMember(data map[string]interface{}, dataType string) error {
	username, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	var query string

	if dataType == "requestToJoinGroup" {
		var userid int
		var user string
		if data["change"] == "request" {
			user = username
			userid = userID
		} else {
			user = data["recipient"].(string)
			userid = GetIDByUsername(user)
		}

		query = datahandler.AddGroupMember
		_, err := datahandler.Database.Exec(
			query, data["groupID"], userid, user, data["state"],
		)
		if err != nil {
			fmt.Println("error adding group member to database: ", err)
			return err
		}
	}

	if dataType == "addGroupMember" {
		query = datahandler.AddGroupMember
		_, err := datahandler.Database.Exec(
			query, data["groupID"], userID, username, data["state"],
		)
		if err != nil {
			fmt.Println("error adding group member to database: ", err)
			return err
		}
	}

	if dataType == "deleteGroupMember" {
		query = datahandler.DeleteGroupMember
		_, err := datahandler.Database.Exec(
			query, data["groupID"], userID,
		)
		if err != nil {
			fmt.Println("error deleting group member from database: ", err)
			return err
		}
	}

	if dataType == "denyGroupMember" {
		query = datahandler.DeleteGroupMember
		_, err := datahandler.Database.Exec(
			query, data["groupID"], GetIDByUsername(data["username"].(string)),
		)
		if err != nil {
			fmt.Println("error deleting group member from database: ", err)
			return err
		}
	}

	if dataType == "acceptGroupMember" {
		query = datahandler.AcceptGroupMember
		fmt.Println("State, group, username", data["state"], data["groupID"], data["username"])
		_, err := datahandler.Database.Exec(
			query, data["state"], data["groupID"], GetIDByUsername(data["username"].(string)),
		)
		if err != nil {
			fmt.Println("error accepting group member into group: ", err)
			return err
		}
	}

	return nil
}

func GetUsernameById(id int) string {
	query := `SELECT username FROM User WHERE id = ?`
	var username string

	err := datahandler.Database.QueryRow(query, id).Scan(&username)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return username
}

func GetIDByUsername(username string) int {
	query := `SELECT id FROM User WHERE username = ?`
	var id int

	err := datahandler.Database.QueryRow(query, username).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return id
}
