package repository

import (
	"context"
	"online-registration/internal/interview/domain/entity"
	"time"
)

type IEventRepository interface {
	CreateEvent(ctx context.Context, title string, description string, startTime, endTime time.Time) (*entity.Event, error)
}
