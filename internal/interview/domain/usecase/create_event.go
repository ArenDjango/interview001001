package usecase

import (
	"context"
	"fmt"
	"online-registration/internal/interview/domain/dto"
	"online-registration/internal/interview/domain/entity"
	"online-registration/internal/interview/domain/repository"

	"github.com/rs/zerolog/log"
)

type CreateEventUseCase struct {
	repository repository.IEventRepository
}

func NewCreateEventUseCase(
	repository repository.IEventRepository,
) *CreateEventUseCase {
	return &CreateEventUseCase{
		repository: repository,
	}
}

func (uc *CreateEventUseCase) CreateEvent(
	ctx context.Context,
	requestDTO *dto.CreateEventRequestDTO,
) (*entity.Event, error) {
	event, err := uc.repository.CreateEvent(
		ctx, requestDTO.Title, requestDTO.Description, requestDTO.StartTime, requestDTO.EndTime,
	)
	if err != nil {
		log.Error().Msgf("CreateEventUseCase.CreateEvent: %v", err)
		return nil, fmt.Errorf("create event: %w", err)
	}
	return event, err
}
