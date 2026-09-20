package service

import (
	"context"
	"errors"
	"strings"

	"github.com/mobi07/secure_bank_auth/internal/auth"
	"github.com/mobi07/secure_bank_auth/internal/domain"
	appErrors "github.com/mobi07/secure_bank_auth/internal/errors"
	"github.com/mobi07/secure_bank_auth/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, fullName, email, password string) (*domain.User, error)
	Login(ctx context.Context, emailId, password string) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
	logger     *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtService *auth.JWTService, logger *zap.Logger) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (s *authService) Register(ctx context.Context, fullName, email, password string) (*domain.User, error) {

	email = strings.ToLower(strings.TrimSpace(email))
	fullName = strings.TrimSpace(fullName)

	// check user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, repository.ErrAlreadyExists
	}

	if !errors.Is(err, repository.ErrNotFound) {
		s.logger.Error("Register: failed to check if user exists", zap.Error(err))
		return nil, err
	}

	// hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Register: failed to hash password", zap.Error(err))
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
		s.logger.Error("Register: failed to create user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", appErrors.ErrInvalidCredentials
		}

		s.logger.Error("Login: failed to get user by email", zap.Error(err))
		return "", err
	}

	if !user.IsActive {
		return "", appErrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", appErrors.ErrInvalidCredentials
	}

	// JWT generation will be implemented next.
	token, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		s.logger.Error("Login: failed to generate access token", zap.Error(err))
		return "", err
	}

	return token, nil
}
