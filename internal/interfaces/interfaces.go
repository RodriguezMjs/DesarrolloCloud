package interfaces

import (
	"context"

	"github.com/RodriguezMjs/tasks-tracking/internal/application/dtos"
	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
	jwtpkg "github.com/RodriguezMjs/tasks-tracking/pkg/jwt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, email, password string) (*entities.User, error)
}

type JWTManager interface {
	GenerateToken(userID, email string) (string, error)
	ValidateToken(tokenString string) (*jwtpkg.Claims, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error)
}
