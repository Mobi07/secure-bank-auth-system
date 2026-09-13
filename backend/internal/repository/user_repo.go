package repository

import (
	"context"

	"github.com/mobi07/secure_bank_auth/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
