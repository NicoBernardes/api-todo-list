package main

import (
	"encoding/json"
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
		Server:   getEnv("DB_SERVER", "localhost\\SQLEXPRESS"),
		Port:     getEnv("DB_PORT", "1433"),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
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

	// Inicializar handlers com a conexão do banco
	todoHandler := handlers.NewTodoHandler(db)

	// Configurar as rotas
	router := mux.NewRouter()
	router.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "database": "connected"})
	}).Methods("GET")

	// Iniciar o servidor
	port := getEnv("PORT", "8080")
	fmt.Printf("🚀 Servidor To-DO API rodando na porta %s...\n", port)
	fmt.Println("📊 Endpoints disponíveis:")
	fmt.Println("   GET    http://localhost:" + port + "/todos")
	fmt.Println("   POST   http://localhost:" + port + "/todos")
	fmt.Println("   GET    http://localhost:" + port + "/todos/{id}")
	fmt.Println("   PUT    http://localhost:" + port + "/todos/{id}")
	fmt.Println("   DELETE http://localhost:" + port + "/todos/{id}")
	fmt.Println("   GET    http://localhost:" + port + "/health")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
