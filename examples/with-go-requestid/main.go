package main

import (
	"context"
	"io"
	"net/http"

	"log/slog"

	requestid "github.com/akm/go-requestid"
	sqlslog "github.com/akm/sql-slog"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqlslogMiddleware := sqlslog.New("sqlite3", ":memory:",
		sqlslog.LogLevel(sqlslog.LevelTrace),
		sqlslog.HandlerFunc(func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
			return requestid.WrapSlogHandler(sqlslog.NewTextHandler(w, opts))
		}),
	)

	ctx := context.Background()

	var err error
	db, err = sqlslogMiddleware.Open(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to open database", "error", err)
		return
	}
	defer db.Close()

	slog.SetDefault(sqlslogMiddleware.Logger())

	createTable(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", getTodos)
	mux.HandleFunc("POST /todos", createTodo)
	mux.HandleFunc("GET /todos/{id}", getTodoByID)
	mux.HandleFunc("PUT /todos/{id}", updateTodoByID)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodoByID)

	slog.InfoContext(ctx, "Starting server on :8080")
	slog.ErrorContext(ctx, http.ListenAndServe(":8080", requestid.Wrap(mux)).Error())
}
