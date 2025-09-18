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
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s; port=%s;database=%s;",
		cfg.Server, cfg.User, cfg.Password, cfg.Port, cfg.Database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar ao banco de dados: %w", err) //Using %w to wrap error
	}

	return db, nil
}

// CreateTable if not exists
func CreateTable(db *sql.DB) {
	query := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='todo' AND xtype='U')
	CREATE TABLE todo (
		id INT IDENTITY(1,1) PRIMARY KEY,
		title NVARCHAR(255) NOT NULL,
		description NVARCHAR(MAX),
		completed BIT DEFAULT 0,
		created_at DATETIME2 DEFAULT GETDATE(),
		updated_at DATETIME2 DEFAULT GETDATE()
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err.Error())
	}

	log.Println("Tabela 'todo' verificada/criada com sucesso!")
}
