package service

import (
	"context"
	"errors"

	"github.com/mobi07/secure_bank_auth/internal/auth"
	"github.com/mobi07/secure_bank_auth/internal/domain"
	"github.com/mobi07/secure_bank_auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, fullName, email, password string) (*domain.User, error)
	Login(ctx context.Context, emailId, password string) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *auth.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(ctx context.Context, fullName, email, password string) (*domain.User, error) {

	// check user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, repository.ErrAlreadyExists
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Fullname:     fullName,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         string(domain.RoleUser),
		IsActive:     true,
	}

	// insert into db
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", repository.ErrInvalidCredentials
		}
	}

	if !user.IsActive {
		return "", repository.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", repository.ErrInvalidCredentials
	}

	// JWT generation will be implemented next.
	token, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
