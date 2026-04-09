package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RodriguezMjs/tasks-tracking/configs"
	"github.com/RodriguezMjs/tasks-tracking/internal/application/usecases"
	"github.com/RodriguezMjs/tasks-tracking/internal/domain/services"
	"github.com/RodriguezMjs/tasks-tracking/internal/infrastructure/http"
	"github.com/RodriguezMjs/tasks-tracking/internal/infrastructure/persistence/postgres"
	"github.com/RodriguezMjs/tasks-tracking/internal/interfaces"
	"github.com/RodriguezMjs/tasks-tracking/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	config := loadConfig()
	db := initDatabase(config)
	defer db.Close()

	userRepository := initRepositories(db)
	authService := initServices(userRepository)
	jwtManager := initJWT(config)
	loginUseCase := initUseCases(authService, jwtManager)
	router := setupRouter(loginUseCase)

	startServer(router, config.GetServerPort())
}

func loadConfig() *configs.Config {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}
	return config
}

func initDatabase(config *configs.Config) *sql.DB {
	dbConfig := configs.NewDatabaseConfig(config)
	db := dbConfig.MustConnect()
	return db
}

func initRepositories(db *sql.DB) interfaces.UserRepository {
	return postgres.NewUserRepository(db)
}

func initServices(userRepo interfaces.UserRepository) interfaces.AuthenticationService {
	return services.NewAuthenticationService(userRepo)
}

func initJWT(config *configs.Config) interfaces.JWTManager {
	return jwt.NewManager(config.GetJWTSecret(), config.GetJWTExpiration())
}

func initUseCases(authService interfaces.AuthenticationService, jwtManager interfaces.JWTManager) interfaces.LoginUseCase {
	return usecases.NewLoginUseCase(authService, jwtManager)
}

func setupRouter(loginUseCase interfaces.LoginUseCase) *gin.Engine {
	router := gin.Default()
	http.SetupRoutes(router, loginUseCase)
	return router
}

func startServer(router *gin.Engine, port string) {
	log.Printf("Servidor iniciando en puerto %s...\n", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
