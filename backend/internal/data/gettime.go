package datahandler

import (
	"strconv"
	"time"
)

func GetTimeSince(date time.Time) string {
	elapsedTime := time.Since(date)
	time := ""
	days := int(elapsedTime.Hours() / 24)
	hours := int(elapsedTime.Hours()) - (days * 24)
	if days != 0 {
		time += strconv.Itoa(days) + "d"
	}
	time += strconv.Itoa(hours) + "h"
	return time
}
