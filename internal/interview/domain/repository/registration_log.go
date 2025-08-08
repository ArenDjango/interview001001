package repository

import (
	"context"
	"online-registration/internal/interview/domain/entity"
	"time"
)

type IEventRepository interface {
	CreateEvent(ctx context.Context, title string, description string, startTime, endTime time.Time) (*entity.Event, error)
	//UpdateRegistrationLogStatus(ctx context.Context, email string,
	//	status dto.RegistrationStatus, errorMessage string) error
	//GetRegistrationLogByEmail(ctx context.Context, email string) (*entity.Event, error)
}
