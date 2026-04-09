package usecases

import (
	"context"
	"fmt"

	"github.com/RodriguezMjs/tasks-tracking/internal/application/dtos"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
)

type LoginUseCase struct {
	authService interfaces.AuthenticationService
	jwtManager  interfaces.JWTManager
}

func NewLoginUseCase(authService interfaces.AuthenticationService, jwtManager interfaces.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		authService: authService,
		jwtManager:  jwtManager,
	}
}

func (u *LoginUseCase) Execute(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := u.authService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("autenticacion fallida: %w", err)
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	roles := make([]dtos.RoleDTO, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = dtos.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	response := &dtos.LoginResponse{
		Token: token,
		User: dtos.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Roles: roles,
	}

	return response, nil
}
