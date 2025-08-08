package contracts

import (
	"context"
	"online-registration/internal/interview/domain/dto"
)

// IBusinessLogService defines the interface for business logging operations
type IBusinessLogService interface {
	// UserRegistrationAttempted logs when a user registration is attempted
	UserRegistrationAttempted(ctx context.Context, dto *dto.CreateUserRequestDTO) error

	// UserRegistrationCompleted logs when a user registration is completed
	UserRegistrationCompleted(ctx context.Context, dto *dto.CreateUserRequestDTO) error

	// UserRegistrationFailed logs when a user dto fails
	UserRegistrationFailed(ctx context.Context, dto *dto.CreateUserRequestDTO, errorMessage string) error
}
