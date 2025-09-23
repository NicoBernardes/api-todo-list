package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Todo representa uma tarefa
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodoHandler gerencia as requisições de TODOs
type TodoHandler struct {
	db *sql.DB
}

// NewTodoHandler cria uma nova instância do TodoHandler
func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

// GetTodos retorna todas as tarefas
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Erro ao buscar tarefas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			http.Error(w, "Erro ao ler tarefa: "+err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo retorna uma tarefa específica
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var todo Todo
	err = h.db.QueryRow(
		"SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = @id",
		sql.Named("id", id),
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar tarefa: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// CreateTodo cria uma nova tarefa
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if todo.Title == "" {
		http.Error(w, "Título é obrigatório", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO todos (title, description, completed) 
	OUTPUT INSERTED.id, INSERTED.created_at, INSERTED.updated_at
	VALUES (@title, @description, @completed)
	`

	var id int
	var createdAt, updatedAt time.Time
	err := h.db.QueryRow(query,
		sql.Named("title", todo.Title),
		sql.Named("description", todo.Description),
		sql.Named("completed", todo.Completed),
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		http.Error(w, "Erro ao criar tarefa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID = id
	todo.CreatedAt = createdAt
	todo.UpdatedAt = updatedAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo atualiza uma tarefa
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	UPDATE todos 
	SET title = @title, description = @description, completed = @completed, updated_at = GETDATE()
	OUTPUT INSERTED.id, INSERTED.title, INSERTED.description, INSERTED.completed, INSERTED.created_at, INSERTED.updated_at
	WHERE id = @id
	`

	var updatedTodo Todo
	err = h.db.QueryRow(query,
		sql.Named("title", todo.Title),
		sql.Named("description", todo.Description),
		sql.Named("completed", todo.Completed),
		sql.Named("id", id),
	).Scan(
		&updatedTodo.ID,
		&updatedTodo.Title,
		&updatedTodo.Description,
		&updatedTodo.Completed,
		&updatedTodo.CreatedAt,
		&updatedTodo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao atualizar tarefa: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

// DeleteTodo remove uma tarefa
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM todos WHERE id = @id", sql.Named("id", id))
	if err != nil {
		http.Error(w, "Erro ao deletar tarefa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erro ao verificar deleção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
