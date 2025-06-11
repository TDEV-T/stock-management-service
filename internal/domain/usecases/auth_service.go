package usecases

import (
	"context"
	"errors"
	"stock-management/internal/domain/models"
	"stock-management/internal/domain/repositories"
	"stock-management/internal/domain/services"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService services.JWTService
}

// UserDTO defines the user data returned to the client.
type UserDTO struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func NewAuthService(userRepo repositories.UserRepository, jwtService services.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*UserDTO, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	user.LastLoginAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, "", err
	}

	tokenString, err := s.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", errors.New("could not generate token: " + err.Error())
	}

	userDto := &UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		LastLoginAt: user.LastLoginAt,
	}

	return userDto, tokenString, nil
}

func (s *AuthService) Register(ctx context.Context, username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	return s.userRepo.Create(ctx, user)
}
