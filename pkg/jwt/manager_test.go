package jwt

import (
	"testing"
	"time"
)

func TestManager_GenerateAndValidateToken(t *testing.T) {
	secret := "test-secret"
	manager := NewManager(secret, time.Hour)

	token, err := manager.GenerateToken("user-123", "test@test.com")
	if err != nil {
		t.Fatalf("expected token generation success, got: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected token validation success, got: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Fatalf("expected user id user-123, got %s", claims.UserID)
	}

	if claims.Email != "test@test.com" {
		t.Fatalf("expected email test@test.com, got %s", claims.Email)
	}
}

func TestManager_ValidateToken_Invalid(t *testing.T) {
	manager := NewManager("test-secret", time.Hour)

	_, err := manager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected validation error for invalid token")
	}
}

func TestManager_ValidateToken_WrongSecret(t *testing.T) {
	manager1 := NewManager("test-secret", time.Hour)
	token, err := manager1.GenerateToken("user-123", "test@test.com")
	if err != nil {
		t.Fatalf("expected token generation success, got: %v", err)
	}

	manager2 := NewManager("other-secret", time.Hour)
	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected validation error for wrong secret")
	}
}
