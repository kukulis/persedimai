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

// GenerateDateRange generates a slice of date strings between startDate and endDate (inclusive)
// Date format: "2025-12-27"
func GenerateDateRange(startDate, endDate string) ([]string, error) {
	if startDate == "" || endDate == "" {
		return []string{time.Now().Format(time.DateOnly)}, nil
	}

	start, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	end, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(time.DateOnly))
	}

	return dates, nil
}
