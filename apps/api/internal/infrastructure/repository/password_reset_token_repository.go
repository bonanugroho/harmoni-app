package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPasswordResetTokenRepository implements PasswordResetTokenRepository.
type PostgresPasswordResetTokenRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresPasswordResetTokenRepository creates a new repository.
func NewPostgresPasswordResetTokenRepository(pool *pgxpool.Pool) *PostgresPasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{pool: pool}
}

// Create inserts a new password reset token.
func (r *PostgresPasswordResetTokenRepository) Create(ctx context.Context, token *PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, used)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.Used,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	return nil
}

// FindByTokenHash finds an unused, unexpired token by its hash.
func (r *PostgresPasswordResetTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1 AND used = false AND expires_at > NOW()
	`

	token := &PasswordResetToken{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Used,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find password reset token: %w", err)
	}

	return token, nil
}

// MarkUsed marks a token as used.
func (r *PostgresPasswordResetTokenRepository) MarkUsed(ctx context.Context, id string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteByUserID removes all tokens for a user.
func (r *PostgresPasswordResetTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens by user ID: %w", err)
	}

	return nil
}
