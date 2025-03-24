package postgres

import (
	"gravitum-test-app/config"
	"gravitum-test-app/internal/repository"
	"gravitum-test-app/internal/repository/postgres/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(cfg *config.Config, db *pgxpool.Pool) *repository.Repository {

	return &repository.Repository{
		User: user.NewRepository(cfg, db),
	}
}
