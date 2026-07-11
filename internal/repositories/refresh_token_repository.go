package repositories

import (
	"context"
	"errors"

	apperrors "github.com/Hettank/habit-tracker/internal/errors"
	"github.com/Hettank/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	token *models.RefreshToken,
) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *RefreshTokenRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*models.RefreshToken, error) {
	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token models.RefreshToken

	err := r.db.QueryRow(
		ctx,
		query,
		tokenHash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) DeleteById(
	ctx context.Context,
	id int64,
) error {
	query := `
		DELETE
		FROM refresh_tokens
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}

func (r *RefreshTokenRepository) DeleteByUserId(
	ctx context.Context,
	userID int64,
) error {
	query := `
		DELETE
		FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
	)

	return err
}

func (r *RefreshTokenRepository) RevokeByTokenHash(
	ctx context.Context,
	tokenHash string,
) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
			AND revoked_at IS NULL;
	`

	result, err := r.db.Exec(
		ctx,
		query,
		tokenHash,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserId(
	ctx context.Context,
	userId int64,
) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1
			AND revoked_at IS NULL
	`

	result, err := r.db.Exec(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrRefreshTokenNotFound
	}

	return nil
}
