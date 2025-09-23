package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"todo-api/database"
	"todo-api/handlers"

	"github.com/gorilla/mux"
)

// LoadConfig carrega as configurações do ambiente
func LoadConfig() *database.Config {
	return &database.Config{
		Server:   getEnv("DB_SERVER", "localhost"),
		Port:     getEnv("DB_PORT", "1433"),
		User:     getEnv("DB_USER", "sa"),
		Password: getEnv("DB_PASSWORD", "YourPassword123"),
		Database: getEnv("DB_NAME", "TodoDB"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Carregar configurações
	cfg := LoadConfig()

	// Conectar ao banco de dados
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err.Error())
	}
	defer db.Close()

	// Verificar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar conexão com o banco de dados:", err.Error())
	}

	fmt.Println("Conectado ao SQL Server com sucesso!")

	// Criar a tabela se não existir
	database.CreateTable(db)

	// Inicializar handlers com a conexão do banco
	todoHandler := handlers.NewTodoHandler(db)

	// Configurar as rotas
	router := mux.NewRouter()
	router.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Iniciar o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor To-DO API rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
