package util

import (
	"fmt"
	"time"
)

func ParseDate(t string) time.Time {
	result, _ := time.Parse(time.DateOnly, t)

	return result
}

func ParseDateTime(t string) time.Time {
	result, _ := time.Parse(time.DateTime, t)

	return result
}

func TryToParseDate(t string, tryFormats []string) (time.Time, error) {
	var err error
	var result time.Time
	for _, format := range tryFormats {
		result, err = time.Parse(format, t)
		if err == nil {
			return result, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse date: %s, error: %s", t, err.Error())
}
