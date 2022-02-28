package service

import "time"

func executionTime(start time.Time) string {
	duration := time.Now().Sub(start)
	return duration.String()
}
