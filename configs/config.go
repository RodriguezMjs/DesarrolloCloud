package configs

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config representa toda la configuración de la aplicación
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig contiene la configuración de la base de datos
type DatabaseConfig struct {
	URL      string `mapstructure:"url"`
	Driver   string `mapstructure:"driver"`
	MaxConns int    `mapstructure:"max_conns"`
}

// JWTConfig contiene la configuración de JWT
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

// ServerConfig contiene la configuración del servidor
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// LoadConfig carga la configuración desde el archivo config.yaml y permite sobreescribirla con variables de entorno.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// En caso de no disponer de config.yaml, se usan valores de respaldo.
	viper.SetDefault("database.url", "postgres://postgres:postgres@db:5432/projectbackend?sslmode=disable")
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.max_conns", 10)
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration_hours", 24)
	viper.SetDefault("server.port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No se pudo leer el archivo de configuración: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetJWTSecret devuelve el secreto JWT
func (c *Config) GetJWTSecret() string {
	return c.JWT.Secret
}

// GetJWTExpiration devuelve la duración de expiración del JWT
func (c *Config) GetJWTExpiration() time.Duration {
	return time.Hour * time.Duration(c.JWT.ExpirationHours)
}

// GetServerPort devuelve el puerto del servidor
func (c *Config) GetServerPort() string {
	return c.Server.Port
}
