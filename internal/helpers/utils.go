package helpers

import "time"

// GetCurrentTimestamp returns the current time formatted as "YYYYMMDD HH:MM"
func GetCurrentTimestamp() string {
	return time.Now().Format("20060102 15:04")
}
