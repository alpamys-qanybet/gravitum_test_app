package user

import (
	"context"
	"errors"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/model"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// table users:
// id
// name
// surname
// inserted_at
// updated_at

type UserRepository struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewRepository(cfg *config.Config, db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		cfg: cfg,
		db:  db,
	}
}

func (r *UserRepository) CheckIfExists(ctx context.Context, id uint) (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.cfg.Db.Timeout)*time.Second)
	defer cancel()

	var exists bool

	err := r.db.QueryRow(timeoutCtx, `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE id = $1
		)
	`, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) GetList(ctx context.Context) ([]*model.User, error) {

	result := []*model.User{}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.cfg.Db.Timeout)*time.Second)
	defer cancel()

	rows, err := r.db.Query(timeoutCtx, `
		SELECT
			id,
			name,
			surname,
			inserted_at,
			updated_at
		FROM users
		ORDER BY inserted_at ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var item model.User

			err = rows.Scan(
				&item.Id,
				&item.Name,
				&item.Surname,
				&item.InsertedAt,
				&item.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}

			result = append(result, &item)
		}
	}
	return result, rows.Err()
}

func (r *UserRepository) Get(ctx context.Context, id uint) (*model.User, error) {
	var result model.User

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.cfg.Db.Timeout)*time.Second)
	defer cancel()

	err := r.db.QueryRow(timeoutCtx, `
		SELECT
			id,
			name,
			surname,
			inserted_at,
			updated_at
		FROM users
		WHERE id = $1;
	`, id).Scan(
		&result.Id,
		&result.Name,
		&result.Surname,
		&result.InsertedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrSqlNoRows
		}
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) Create(
	ctx context.Context,
	name string,
	surname *string,
) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.cfg.Db.Timeout)*time.Second)
	defer cancel()

	_, err := r.db.Exec(timeoutCtx, `
		INSERT INTO users (
			name,
			surname
		)
		VALUES ($1, $2);
	`,
		name,
		surname,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	id uint,
	name string,
	surname *string,
) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(r.cfg.Db.Timeout)*time.Second)
	defer cancel()

	_, err := r.db.Exec(timeoutCtx, `
			UPDATE users
			SET name = $2,
				surname = $3,
				updated_at = $4
			WHERE id = $1;
		`,
		id,
		name,
		surname,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
