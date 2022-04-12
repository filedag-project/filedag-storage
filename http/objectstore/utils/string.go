package utils

import "time"

//String return points
func String(v string) *string {
	return &v
}

//String return points
func Time(v time.Time) *time.Time {
	return &v
}
