package service

import (
	"context"
	"errors"
	"online-quiz/internal/domain"
	"online-quiz/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, password string, role domain.Role) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtSvc   JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtSvc JWTService) AuthService {
	return &authService{userRepo, jwtSvc}
}

func (s *authService) Register(ctx context.Context, username, password string, role domain.Role) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwtSvc.GenerateToken(user)
}
