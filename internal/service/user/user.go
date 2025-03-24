package user

import (
	"context"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/model"
	"gravitum-test-app/internal/repository"
)

type UserService struct {
	cfg  *config.Config
	repo repository.UserRepository
}

func NewService(
	cfg *config.Config,
	repo repository.UserRepository,
) *UserService {
	return &UserService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *UserService) GetList(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetList(ctx)
}

func (s *UserService) Get(ctx context.Context, id uint) (*model.User, error) {
	exists, err := s.repo.CheckIfExists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, model.ErrNoUserWithSuchId
	}

	return s.repo.Get(ctx, id)
}

func (s *UserService) Create(
	ctx context.Context,
	name string,
	surname *string,
) error {
	return s.repo.Create(ctx, name, surname)
}

func (s *UserService) Update(
	ctx context.Context,
	id uint,
	name string,
	surname *string,
) error {
	exists, err := s.repo.CheckIfExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return model.ErrNoUserWithSuchId
	}

	return s.repo.Update(ctx, id, name, surname)
}
