package repository

import (
	"context"
	"gravitum-test-app/internal/model"
	"gravitum-test-app/internal/repository/postgres/user"
)

type UserRepository interface {
	CheckIfExists(ctx context.Context, id uint) (bool, error)
	Get(ctx context.Context, id uint) (*model.User, error)
	GetList(ctx context.Context) ([]*model.User, error)
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

type Repository struct {
	User UserRepository
}

var _ UserRepository = (*user.UserRepository)(nil)
