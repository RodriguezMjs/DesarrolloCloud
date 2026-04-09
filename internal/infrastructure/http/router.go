package http

import (
	"github.com/RodriguezMjs/tasks-tracking/internal/infrastructure/http/handlers"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, loginUseCase interfaces.LoginUseCase) {
	authHandler := handlers.NewAuthHandler(loginUseCase)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}
