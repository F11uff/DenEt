package service

import (
	"context"
	"errors"
	"time"
	"denet/internal/model"
	"denet/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.User, error)
	GenerateToken(user *model.User, jwtSecret string) (string, error)
}

type authService struct {
	uow repository.UnitOfWork
}

func NewAuthService(uow repository.UnitOfWork) AuthService {
	return &authService{uow: uow}
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	_, err := s.uow.Users().GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserExists
	}
	_, err = s.uow.Users().GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserExists
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Balance:  0,
	}
	err = s.uow.Users().CreateWithPassword(ctx, user, req.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, error) {
	user, err := s.uow.Users().VerifyPassword(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) GenerateToken(user *model.User, jwtSecret string) (string, error) {
	claims := &model.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}