package repository

import (
	"context"
	"fmt"
	"time"

	"online-registration/internal/interview/domain/entity"
	"online-registration/internal/interview/infrastructure/db/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EventRepository struct {
	db *bun.DB
}

func NewDBEventRepository(db *bun.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) CreateEvent(
	ctx context.Context,
	title string, description string, startTime time.Time, endTime time.Time,
) (*entity.Event, error) {

	registrationLog := &entity.Event{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	model := &model.Event{}
	model = model.ToModel(*registrationLog)

	_, err := r.
		db.
		NewInsert().
		Model(model).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("CreateEvent %w", err)
	}

	return model.ToEntity(), nil
}

//
//func (r *EventRepository) UpdateRegistrationLogStatus(
//	ctx context.Context,
//	email string,
//	status dto.RegistrationStatus,
//	errorMessage string,
//) error {
//	ctx, span := otellib.StartSpan(ctx, r.tracer,
//		"EventRepository.UpdateRegistrationLogStatus",
//		trace.SpanKindInternal)
//	defer otellib.EndSpan(span)
//
//	for attempt := 1; attempt <= MaxAttempts; attempt++ {
//		query := r.
//			db.
//			NewUpdate().
//			Model((*model.Event)(nil)).
//			Set("status = ?", status.String()).
//			Set("updated_at = ?", time.Now()).
//			Where("email = ?", email)
//
//		if errorMessage != "" {
//			query = query.Set("error_message = ?", errorMessage)
//		}
//
//		result, err := query.Exec(ctx)
//		if err != nil {
//			return fmt.Errorf("UpdateRegistrationLogStatus %w", err)
//		}
//
//		rowsAffected, _ := result.RowsAffected()
//		if rowsAffected > 0 {
//			log.Debug().
//				Str("email", email).
//				Int("attempt", attempt).
//				Msg("Successfully updated registration log status")
//			return nil
//		}
//
//		if attempt < MaxAttempts {
//			log.Debug().
//				Str("email", email).
//				Int("attempt", attempt).
//				Msg("Registration log not found, retrying...")
//
//			select {
//			case <-ctx.Done():
//				return fmt.Errorf("UpdateRegistrationLogStatus cancelled: %w", ctx.Err())
//			case <-time.After(RetryDelay):
//			}
//		}
//	}
//
//	return fmt.Errorf("UpdateRegistrationLogStatus failed after %d "+
//		"attempts: no registration log found for email: %s", MaxAttempts, email)
//}
//
//func (r *EventRepository) GetRegistrationLogByEmail(
//	ctx context.Context,
//	email string,
//) (*entity.Event, error) {
//	ctx, span := otellib.StartSpan(ctx, r.tracer,
//		"EventRepository.GetRegistrationLogByID",
//		trace.SpanKindInternal)
//	defer otellib.EndSpan(span)
//
//	modelRegistrationLog := new(model.Event)
//	err := r.
//		db.
//		NewSelect().
//		Model(modelRegistrationLog).
//		Where("email = ?", email).
//		Scan(ctx)
//
//	if err != nil {
//		return nil, fmt.Errorf("GetRegistrationLogByID %w", err)
//	}
//
//	return modelRegistrationLog.ToEntity(), nil
//}
