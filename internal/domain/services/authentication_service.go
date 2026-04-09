package services

import (
	"context"
	"fmt"

	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userRepository interfaces.UserRepository
}

func NewAuthenticationService(userRepo interfaces.UserRepository) *AuthenticationService {
	return &AuthenticationService{
		userRepository: userRepo,
	}
}

func (s *AuthenticationService) Authenticate(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("usuario o contraseña incorrectos")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error al generar hash: %w", err)
	}
	return string(hash), nil
}
