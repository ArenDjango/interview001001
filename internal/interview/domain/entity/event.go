package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"table:events,alias:s"`
	ID            uuid.UUID
	Title         string
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	CreatedAt     time.Time
}
