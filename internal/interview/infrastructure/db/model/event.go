package model

import (
	"online-registration/internal/interview/domain/entity"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"table:events,alias:s"`
	ID            uuid.UUID `bun:"id,pk,notnull"`
	Title         string    `bun:"title,notnull"`
	Description   string    `bun:"description,notnull"`
	StartTime     time.Time `bun:"start_time,notnull"`
	EndTime       time.Time `bun:"end_time,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
}

func (m *Event) ToEntity() *entity.Event {
	return &entity.Event{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		StartTime:   m.StartTime,
		EndTime:     m.EndTime,
		CreatedAt:   m.CreatedAt,
	}
}

func (m *Event) ToModel(entity entity.Event) *Event {
	return &Event{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		StartTime:   entity.StartTime,
		EndTime:     entity.EndTime,
		CreatedAt:   entity.CreatedAt,
	}
}
