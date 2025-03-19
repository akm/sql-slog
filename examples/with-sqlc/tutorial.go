package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"reflect"

	sqlslog "github.com/akm/sql-slog"
	_ "modernc.org/sqlite"

	"tutorial.sqlc.dev/app/tutorial"
)

//go:embed schema.sql
var ddl string

func run() error {
	ctx := context.Background()

	var db *sql.DB
	var err error
	if os.Getenv("SKIP_SQLSLOG") != "" {
		db, err = sql.Open("sqlite", ":memory:")
	} else {
		db, _, err = sqlslog.Open(ctx, "sqlite", ":memory:")
	}
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
