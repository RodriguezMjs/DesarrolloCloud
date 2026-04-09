package configs

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// NewDatabaseConfig crea una nueva configuración de base de datos desde la configuración general
func NewDatabaseConfig(config *Config) *DatabaseConfig {
	return &DatabaseConfig{
		URL:      config.Database.URL,
		Driver:   config.Database.Driver,
		MaxConns: config.Database.MaxConns,
	}
}

// Connect establece la conexión a la base de datos
func (c *DatabaseConfig) Connect() (*sql.DB, error) {
	db, err := sql.Open(c.Driver, c.URL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a base de datos: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(c.MaxConns)
	db.SetMaxIdleConns(c.MaxConns / 2)

	// Ping a db
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error verificando conexion a base de datos: %w", err)
	}

	log.Println("Conectado a PostgreSQL!")
	return db, nil
}

// MustConnect establece la conexión y termina el programa si falla
func (c *DatabaseConfig) MustConnect() *sql.DB {
	db, err := c.Connect()
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	return db
}
