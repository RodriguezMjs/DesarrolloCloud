package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/RodriguezMjs/tasks-tracking/internal/application/dtos"
	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
	jwtpkg "github.com/RodriguezMjs/tasks-tracking/pkg/jwt"
)

type fakeAuthService struct {
	user *entities.User
	err  error
}

func (f *fakeAuthService) Authenticate(ctx context.Context, email, password string) (*entities.User, error) {
	return f.user, f.err
}

type fakeJWTManager struct {
	token string
	err   error
}

func (f *fakeJWTManager) GenerateToken(userID, email string) (string, error) {
	return f.token, f.err
}

func (f *fakeJWTManager) ValidateToken(tokenString string) (*jwtpkg.Claims, error) {
	return nil, nil
}

var _ interfaces.JWTManager = (*fakeJWTManager)(nil)

func TestLoginUseCase_Execute_Success(t *testing.T) {
	user := &entities.User{
		ID:    "user-123",
		Name:  "Test User",
		Email: "test@test.com",
		Roles: []entities.Role{{ID: 1, Name: "ADMIN"}},
	}

	login := NewLoginUseCase(&fakeAuthService{user: user}, &fakeJWTManager{token: "jwt-token"})

	response, err := login.Execute(context.Background(), &dtos.LoginRequest{Email: "test@test.com", Password: "test123"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if response.Token != "jwt-token" {
		t.Fatalf("expected token jwt-token, got %s", response.Token)
	}

	if response.User.Email != "test@test.com" || response.User.Name != "Test User" {
		t.Fatalf("unexpected user response: %+v", response.User)
	}

	if len(response.Roles) != 1 || response.Roles[0].Name != "ADMIN" {
		t.Fatalf("unexpected roles response: %+v", response.Roles)
	}
}

func TestLoginUseCase_Execute_AuthError(t *testing.T) {
	login := NewLoginUseCase(&fakeAuthService{err: errors.New("usuario no encontrado")}, &fakeJWTManager{token: "jwt-token"})

	_, err := login.Execute(context.Background(), &dtos.LoginRequest{Email: "test@test.com", Password: "test123"})
	if err == nil || err.Error() != "autenticacion fallida: usuario no encontrado" {
		t.Fatalf("expected authentication failure, got: %v", err)
	}
}

func TestLoginUseCase_Execute_JWTError(t *testing.T) {
	user := &entities.User{
		ID:    "user-123",
		Name:  "Test User",
		Email: "test@test.com",
		Roles: []entities.Role{{ID: 1, Name: "ADMIN"}},
	}

	login := NewLoginUseCase(&fakeAuthService{user: user}, &fakeJWTManager{err: errors.New("failure generating token")})

	_, err := login.Execute(context.Background(), &dtos.LoginRequest{Email: "test@test.com", Password: "test123"})
	if err == nil || err.Error() != "error al generar token: failure generating token" {
		t.Fatalf("expected jwt generation failure, got: %v", err)
	}
}
