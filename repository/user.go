package repository

import (
	"fmt"
	"errors"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/nekto-sns/server/model"
	"github.com/nekto-sns/server/domain"
)

type userRepository struct{
	db *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) (*userRepository) {
	return &userRepository{
		db: pool,
	}
}

func (r *userRepository) ByName(ctx context.Context, name string) (*model.User, error) {
	var u model.User

	query := "SELECT id, scratch_name, display_name, bio, created_at FROM users WHERE scratch_name = $1"
	err := r.db.QueryRow(ctx, query, name).Scan(&u.ID, &u.ScratchName, &u.DisplayName, &u.Bio, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Row not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("DB query failed (%v): %w", err, domain.ErrInternal)
	}

	return &u, nil
}

func (r *userRepository) ByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	query := "SELECT id, scratch_name, display_name, bio, created_at FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.ScratchName, &u.DisplayName, &u.Bio, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Row not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("DB query failed (%v): %w", err, domain.ErrInternal)
	}

	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, scratchName string) (*model.User, error) {
	var u model.User

	query := `
		INSERT INTO users (scratch_name)
		VALUES ($1)
		RETURNING id, scratch_name, display_name, bio, created_at`

	err := r.db.QueryRow(ctx, query, scratchName).Scan(
		&u.ID, &u.ScratchName, &u.DisplayName, &u.Bio, &u.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, fmt.Errorf("DB insert failed (%v): %w", err, domain.ErrConflict)
		}
		return nil, fmt.Errorf("DB insert failed (%v): %w", err, domain.ErrInternal)
	}

	return &u, nil
}

func (r *userRepository) IDByScratchName(ctx context.Context, scratchName string) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx,`SELECT id FROM users WHERE scratch_name = $1`, scratchName).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("Row not found: %w", domain.ErrNotFound)
		}
		return "", fmt.Errorf("DB query failed (%v): %w", err, domain.ErrInternal)
	}

	return userID, nil
}
