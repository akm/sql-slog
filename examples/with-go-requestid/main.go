package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func main() {
	ctx := context.Background()

	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to open database", "error", err)
		return
	}
	defer db.Close()

	createTable(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", getTodos)
	mux.HandleFunc("POST /todos", createTodo)
	mux.HandleFunc("GET /todos/{id}", getTodoByID)
	mux.HandleFunc("PUT /todos/{id}", updateTodoByID)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodoByID)

	slog.InfoContext(ctx, "Starting server on :8080")
	slog.ErrorContext(ctx, http.ListenAndServe(":8080", mux).Error())
}

func createTable(ctx context.Context) {
	query := `
	CREATE TABLE todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		status TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create table", "error", err)
		return
	}
	slog.InfoContext(ctx, "Table created successfully")
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "getTodos handler started")
	defer slog.InfoContext(ctx, "getTodos handler ended")

	rows, err := db.Query("SELECT id, title, status FROM todos")
	if err != nil {
		slog.ErrorContext(ctx, "Error querying todos", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			slog.ErrorContext(ctx, "Error scanning todo", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "createTodo handler started")
	defer slog.InfoContext(ctx, "createTodo handler ended")

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.ErrorContext(ctx, "Error decoding todo", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO todos (title, status) VALUES (?, ?)", todo.Title, todo.Status)
	if err != nil {
		slog.ErrorContext(ctx, "Error inserting todo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting last insert ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func getTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "getTodoByID handler started")
	defer slog.InfoContext(ctx, "getTodoByID handler ended")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid ID", "error", err)
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	var todo Todo
	if err := db.QueryRow("SELECT id, title, status FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
		if err == sql.ErrNoRows {
			slog.InfoContext(ctx, "Todo not found")
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			slog.ErrorContext(ctx, "Error querying todo by ID", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "updateTodoByID handler started")
	defer slog.InfoContext(ctx, "updateTodoByID handler ended")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid ID", "error", err)
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.ErrorContext(ctx, "Error decoding todo", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := db.Exec("UPDATE todos SET title = ?, status = ? WHERE id = ?", todo.Title, todo.Status, id); err != nil {
		slog.ErrorContext(ctx, "Error updating todo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func deleteTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "deleteTodoByID handler started")
	defer slog.InfoContext(ctx, "deleteTodoByID handler ended")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid ID", "error", err)
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	if _, err := db.Exec("DELETE FROM todos WHERE id = ?", id); err != nil {
		slog.ErrorContext(ctx, "Error deleting todo", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
