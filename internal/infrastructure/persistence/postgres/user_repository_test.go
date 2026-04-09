package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
)

func TestUserRepository_GetByEmail_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`)).
		WithArgs("test@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}).AddRow("user-1", "Test User", "test@test.com", "hash", createdAt))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT r.id, r.name FROM roles r
			 INNER JOIN user_roles ur ON r.id = ur.role_id
			 WHERE ur.user_id = $1`)).
		WithArgs("user-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "admin"))

	user, err := repo.GetByEmail(context.Background(), "test@test.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected a user, got nil")
	}

	if user.Email != "test@test.com" {
		t.Fatalf("unexpected email: got %q", user.Email)
	}

	if len(user.Roles) != 1 || user.Roles[0].Name != "admin" {
		t.Fatalf("unexpected roles: %+v", user.Roles)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`)).
		WithArgs("missing@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}))

	user, err := repo.GetByEmail(context.Background(), "missing@test.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	user := &entities.User{
		ID:           "user-1",
		Name:         "Test User",
		Email:        "test@test.com",
		PasswordHash: "hash",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`)).
		WithArgs(user.ID, user.Name, user.Email, user.PasswordHash).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
