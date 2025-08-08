package contracts

import (
	"context"
	"online-registration/internal/interview/domain/dto"

	"gitlab.b2broker.tech/pbsr/pbsr/backend/go/libs/core/pkg/v1/avanpost_proxy"
)

type IAvanpostClient interface {
	CreateUser(ctx context.Context, dto *dto.CreateUserRequestDTO) (*avanpost_proxy.User, error)
	AddUserToGroup(ctx context.Context, userUuid, groupUuid string) error
}
