package services

import (
	"context"
	"errors"
	"testing"

	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
)

type fakeUserRepo struct {
	user *entities.User
	err  error
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return f.user, f.err
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id string) (*entities.User, error) {
	return nil, nil
}

func (f *fakeUserRepo) Create(ctx context.Context, user *entities.User) error {
	return nil
}

func TestAuthenticationService_Authenticate_Success(t *testing.T) {
	hash, err := HashPassword("test123")
	if err != nil {
		t.Fatalf("error generating password hash: %v", err)
	}

	user := &entities.User{
		ID:           "user-123",
		Name:         "Test User",
		Email:        "test@test.com",
		PasswordHash: hash,
	}

	service := NewAuthenticationService(&fakeUserRepo{user: user})

	result, err := service.Authenticate(context.Background(), "test@test.com", "test123")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if result == nil || result.Email != "test@test.com" {
		t.Fatalf("unexpected user result: %+v", result)
	}
}

func TestAuthenticationService_Authenticate_UserNotFound(t *testing.T) {
	service := NewAuthenticationService(&fakeUserRepo{user: nil})

	_, err := service.Authenticate(context.Background(), "missing@test.com", "test123")
	if err == nil || err.Error() != "usuario no encontrado" {
		t.Fatalf("expected usuario no encontrado error, got: %v", err)
	}
}

func TestAuthenticationService_Authenticate_InvalidPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("error generating password hash: %v", err)
	}

	user := &entities.User{
		ID:           "user-123",
		Name:         "Test User",
		Email:        "test@test.com",
		PasswordHash: hash,
	}

	service := NewAuthenticationService(&fakeUserRepo{user: user})

	_, err = service.Authenticate(context.Background(), "test@test.com", "wrong-password")
	if err == nil || err.Error() != "usuario o contraseña incorrectos" {
		t.Fatalf("expected invalid credentials error, got: %v", err)
	}
}

func TestAuthenticationService_Authenticate_RepoError(t *testing.T) {
	service := NewAuthenticationService(&fakeUserRepo{err: errors.New("db failure")})

	_, err := service.Authenticate(context.Background(), "test@test.com", "test123")
	if err == nil || err.Error() != "error al buscar usuario: db failure" {
		t.Fatalf("expected repository error, got: %v", err)
	}
}
