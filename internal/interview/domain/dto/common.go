package dto

import "time"

type CreateEventRequestDTO struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}
