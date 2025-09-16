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

func main() {
	// load config
	cfg := config.loadConfig()

	// connect db
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err.Error())
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar conexão com o banco de dados: ", err.Error())
	}

	fmt.Println("Conectado ao SQL Server com sucesso!")

	// Create table if not exists
	database.CreateTable(db)

	// Inicialize handlers with db conection
	todoHandler := handlers.NewTodoHandler(db)

	// Config routes
	router := mux.NewRouter()
	router.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandler.getTodo).Methods("GET")
	router.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Init server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor To-DO API rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
