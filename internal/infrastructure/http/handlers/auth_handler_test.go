package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RodriguezMjs/tasks-tracking/internal/application/dtos"
	"github.com/gin-gonic/gin"
)

type fakeLoginUseCase struct {
	response *dtos.LoginResponse
	err      error
	req      *dtos.LoginRequest
}

func (f *fakeLoginUseCase) Execute(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	f.req = req
	return f.response, f.err
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	fakeResponse := &dtos.LoginResponse{
		Token: "fake-token",
		User: dtos.UserDTO{
			ID:    "user-1",
			Name:  "Test User",
			Email: "test@test.com",
		},
		Roles: []dtos.RoleDTO{{ID: 1, Name: "admin"}},
	}

	useCase := &fakeLoginUseCase{response: fakeResponse}
	handler := NewAuthHandler(useCase)

	requestBody, err := json.Marshal(dtos.LoginRequest{Email: "test@test.com", Password: "test123"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Token != fakeResponse.Token {
		t.Fatalf("unexpected token: got %q, want %q", response.Token, fakeResponse.Token)
	}

	if response.User.Email != fakeResponse.User.Email {
		t.Fatalf("unexpected user email: got %q, want %q", response.User.Email, fakeResponse.User.Email)
	}
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	useCase := &fakeLoginUseCase{}
	handler := NewAuthHandler(useCase)

	requestBody, err := json.Marshal(map[string]string{"email": "invalid-email"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response dtos.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if response.Error != "VALIDATION_ERROR" {
		t.Fatalf("unexpected error code: got %q", response.Error)
	}
}

func TestAuthHandler_Login_AuthenticationFailed(t *testing.T) {
	useCase := &fakeLoginUseCase{err: errors.New("invalid credentials")}
	handler := NewAuthHandler(useCase)

	requestBody, err := json.Marshal(dtos.LoginRequest{Email: "test@test.com", Password: "wrongpass"})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response dtos.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if response.Error != "AUTH_FAILED" {
		t.Fatalf("unexpected error code: got %q", response.Error)
	}
}
