package dto

import "time"

type RegistrationStatus string

const (
	RegistrationStatusAttempting RegistrationStatus = "attempted"
	RegistrationStatusSuccess    RegistrationStatus = "success"
	RegistrationStatusFailed     RegistrationStatus = "failed"
)

func (s RegistrationStatus) String() string {
	return string(s)
}

func (s RegistrationStatus) IsSuccess() bool {
	return s == RegistrationStatusSuccess
}

func (s RegistrationStatus) IsFailed() bool {
	return s == RegistrationStatusFailed
}
func (s RegistrationStatus) IsAttempting() bool {
	return s == RegistrationStatusAttempting
}

type CreateEventRequestDTO struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}
