package mysqltest

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestBasic(t *testing.T) {
	dbName := "app1"

	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		t.Fatal(err)
	}
	defer exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run()

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/"+dbName)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		t.Logf("Ping failed: %v", err)
		time.Sleep(2 * time.Second)
	}

	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test1 (id INT PRIMARY KEY, name VARCHAR(255))"); err != nil {
		t.Fatal(err)
	}

	t.Run("without tx", func(t *testing.T) {
		testData := []string{"foo", "bar", "baz"}
		for i, name := range testData {
			if _, err := db.ExecContext(ctx, "INSERT INTO test1 (id, name) VALUES (?, ?)", i+1, name); err != nil {
				t.Fatal(err)
			}
		}

		rows, err := db.QueryContext(ctx, "SELECT id, name FROM test1 WHERE name LIKE ? ORDER BY id", "b%")
		if err != nil {
			t.Fatal(err)
		}

		results := []map[string]interface{}{}
		for rows.Next() {
			result := map[string]interface{}{}
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				t.Fatal(err)
			}
			result["id"] = id
			result["name"] = name
			results = append(results, result)
		}

		if len(results) != 2 {
			t.Fatalf("Expected 2 results, got %d", len(results))
		}

		if results[0]["id"] != 2 || results[0]["name"] != "bar" {
			t.Fatalf("Unexpected result: %v", results[0])
		}
		if results[1]["id"] != 3 || results[1]["name"] != "baz" {
			t.Fatalf("Unexpected result: %v", results[1])
		}
	})

	t.Run("with tx", func(t *testing.T) {
		t.Run("rollback", func(t *testing.T) {
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			r, err := tx.ExecContext(ctx, "UPDATE test1 SET name = ? WHERE id = ?", "quux", 3)
			if err != nil {
				t.Fatal(err)
			}
			rowsAffected, err := r.RowsAffected()
			if err != nil {
				t.Fatal(err)
			}
			if rowsAffected != 1 {
				t.Fatalf("Expected 1 row affected, got %d", rowsAffected)
			}
			if err := tx.Rollback(); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("commit", func(t *testing.T) {
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			r, err := tx.ExecContext(ctx, "UPDATE test1 SET name = ? WHERE id = ?", "qux", 3)
			if err != nil {
				t.Fatal(err)
			}
			rowsAffected, err := r.RowsAffected()
			if err != nil {
				t.Fatal(err)
			}
			if rowsAffected != 1 {
				t.Fatalf("Expected 1 row affected, got %d", rowsAffected)
			}
			if err := tx.Commit(); err != nil {
				t.Fatal(err)
			}
		})

	})
}
