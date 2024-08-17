package api

import (
	"fmt"
	datahandler "social-network/backend/internal/data"
	"social-network/backend/internal/sessions"
)

func CreateEvent(data map[string]interface{}) (string, int, error) {
	_, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	result, err := datahandler.Database.Exec(
		datahandler.CreateEvent, data["title"], data["content"], data["date"], userID, data["groupID"],
	)
	if err != nil {
		fmt.Println("error inserting event to database: ", err)
		return "", 0, err
	}

	lastInsert, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error retrieving event id after posting: ", err)
		return "", 0, err
	}

	return "Event successfully created!", int(lastInsert), nil
}

func AddEventAttendeeWhenCreatingEvent(data map[string]interface{}, eventID int) error {
	username, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	_, err := datahandler.Database.Exec(
		datahandler.AddEventAttendee, eventID, userID, username, data["attendance"],
	)
	if err != nil {
		fmt.Println("error inserting event member to database: ", err)
		return err
	}

	return nil
}

func AddDeleteUpdateEventAttendee(data map[string]interface{}, dataType string) error {
	username, userID := sessions.GetUserNameAndID(data["sessionID"].(string))

	var query string

	if dataType == "addEventAttendee" {
		query = datahandler.AddEventAttendee
		_, err := datahandler.Database.Exec(
			query, data["eventID"], userID, username, data["attendance"],
		)
		if err != nil {
			fmt.Println("error inserting attendee to database: ", err)
			return err
		}
	}

	if dataType == "deleteEventAttendee" {
		query = datahandler.DeleteEventAttendee
		_, err := datahandler.Database.Exec(
			query, data["eventID"], userID,
		)
		if err != nil {
			fmt.Println("error deleting attendee from database: ", err)
			return err
		}
	}

	if dataType == "updateEventAttendee" {
		query = datahandler.UpdateEventAttendee
		_, err := datahandler.Database.Exec(
			query, data["attendance"], data["eventID"], userID,
		)
		if err != nil {
			fmt.Println("error updating attendee in database: ", err)
			return err
		}
	}

	return nil
}
