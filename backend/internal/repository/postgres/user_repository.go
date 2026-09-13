package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobi07/secure_bank_auth/internal/domain"
	"github.com/mobi07/secure_bank_auth/internal/repository"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(database *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			full_name,
			email,
			password_hash,
			role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Fullname, user.Email, user.PasswordHash, user.Role, user.IsActive).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.Fullname, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, emailId string) (*domain.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, emailId).Scan(&user.Id, &user.Fullname, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
