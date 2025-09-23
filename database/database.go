package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Get db config

type Config struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

func Connect(cfg *Config) (*sql.DB, error) {
	var connString string
	if cfg.User == "" && cfg.Password == "" {
		// Windows Authentication
		connString = fmt.Sprintf("server=%s;port=%s;database=%s;trusted_connection=yes;",
			cfg.Server, cfg.Port, cfg.Database)
	} else {
		// SQL Server Authentication
		connString = fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;",
			cfg.Server, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	}

	log.Printf("Tentando conectar com: server=%s, database=%s", cfg.Server, cfg.Database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	return db, nil
}

// CheckTableExists
func CheckTableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
	SELECT COUNT(*) 
	FROM INFORMATION_SCHEMA.TABLES 
	WHERE TABLE_NAME = @tableName AND TABLE_TYPE = 'BASE TABLE'
	`

	var count int
	err := db.QueryRow(query, sql.Named("tableName", tableName)).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
