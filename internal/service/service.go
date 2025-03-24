package service

import (
	"context"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/model"
	"gravitum-test-app/internal/repository"
	"gravitum-test-app/internal/service/user"
)

type UserService interface {
	GetList(ctx context.Context) ([]*model.User, error)
	Get(ctx context.Context, id uint) (*model.User, error)
	Create(
		ctx context.Context,
		name string,
		surname *string,
	) error
	Update(
		ctx context.Context,
		id uint,
		name string,
		surname *string,
	) error
}

type Service struct {
	User UserService
}

func NewService(
	cfg *config.Config,
	repositories *repository.Repository,
) *Service {
	return &Service{
		User: user.NewService(
			cfg,
			repositories.User,
		),
	}
}

var _ UserService = (*user.UserService)(nil)
