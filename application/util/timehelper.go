package util

import "time"

func ParseDate(t string) time.Time {
	result, _ := time.Parse(time.DateOnly, t)

	return result
}

func ParseDateTime(t string) time.Time {
	result, _ := time.Parse(time.DateTime, t)

	return result
}
