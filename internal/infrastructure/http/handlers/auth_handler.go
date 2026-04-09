package handlers

import (
	"github.com/RodriguezMjs/tasks-tracking/internal/application/dtos"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	loginUseCase interfaces.LoginUseCase
}

func NewAuthHandler(loginUseCase interfaces.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Email y password requeridos",
		})
		return
	}

	response, err := h.loginUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
			Error:   "AUTH_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
