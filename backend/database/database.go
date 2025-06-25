package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB es la instancia global de la base de datos
var DB *sql.DB

// Config contiene la configuración de la base de datos
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establece la conexión con PostgreSQL
func Connect(cfg Config) error {
	// Construir DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Abrir conexión
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error abriendo conexión a la base de datos: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("error verificando conexión a la base de datos: %w", err)
	}

	DB = db
	log.Println("Conexión a PostgreSQL establecida exitosamente")
	return nil
}

// Close cierra la conexión a la base de datos
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *sql.DB {
	return DB
}

// InitializeSchema ejecuta las migraciones iniciales
func InitializeSchema() error {
	if DB == nil {
		return fmt.Errorf("base de datos no conectada")
	}

	// Crear extensión UUID si no existe
	_, err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("error creando extensión uuid-ossp: %w", err)
	}

	log.Println("Schema de base de datos inicializado")
	return nil
}

// HealthCheck verifica que la conexión a la base de datos esté activa
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("base de datos no conectada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("error en health check de base de datos: %w", err)
	}

	return nil
}