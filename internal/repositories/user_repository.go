package repositories

import (
	"context"
	"errors"

	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			created_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			email,
			password_hash
		)
		VALUES (
			$1,
			$2
		)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	return err
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.User, error) {

	query := `
		SELECT
			id,
			email,
			password_hash,
			created_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
