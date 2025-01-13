package mysqltest

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"

	sqlslog "github.com/akm/sql-slog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	dbName := "app1"

	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		t.Fatal(err)
	}
	defer exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run()

	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := sqlslog.Open(ctx, "mysql", "root@tcp(localhost:3306)/"+dbName, logger)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		t.Logf("Ping failed: %v", err)
		time.Sleep(2 * time.Second)
	}

	t.Run("Ping", func(t *testing.T) {
		buf.Reset()
		err := db.PingContext(ctx)
		assert.NoError(t, err)
		actualEntries := parseJsonLines(t, buf.Bytes())
		exptectedEntries := []map[string]interface{}{
			{"level": "DEBUG", "msg": "ResetSession Start", "driver": "mysql", "dsn": "root@tcp(localhost:3306)/app1"},
			{"level": "INFO", "msg": "ResetSession Complete", "driver": "mysql", "dsn": "root@tcp(localhost:3306)/app1"},
			{"level": "DEBUG", "msg": "Ping Start", "driver": "mysql", "dsn": "root@tcp(localhost:3306)/app1"},
			{"level": "INFO", "msg": "Ping Complete", "driver": "mysql", "dsn": "root@tcp(localhost:3306)/app1"},
		}
		assert.Len(t, actualEntries, len(exptectedEntries))
		for i, expected := range exptectedEntries {
			for k, v := range expected {
				assert.Equal(t, v, actualEntries[i][k])
			}
		}
	})

	t.Run("create table", func(t *testing.T) {
		result, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test1 (id INT PRIMARY KEY, name VARCHAR(255))")
		assert.NoError(t, err)
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected)
	})

	t.Run("without tx", func(t *testing.T) {
		testData := []string{"foo", "bar", "baz"}
		for i, name := range testData {
			t.Run("insert "+name, func(t *testing.T) {
				result, err := db.ExecContext(ctx, "INSERT INTO test1 (id, name) VALUES (?, ?)", i+1, name)
				assert.NoError(t, err)
				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
		}

		t.Run("select", func(t *testing.T) {
			rows, err := db.QueryContext(ctx, "SELECT id, name FROM test1 WHERE name LIKE ?", "ba%")
			assert.NoError(t, err)
			defer rows.Close()

			actualResults := []map[string]interface{}{}
			for rows.Next() {
				result := map[string]interface{}{}
				var id int
				var name string
				if err := rows.Scan(&id, &name); err != nil {
					t.Fatal(err)
				}
				result["id"] = id
				result["name"] = name
				actualResults = append(actualResults, result)
			}
			expectedResults := []map[string]interface{}{
				{"id": 2, "name": "bar"},
				{"id": 3, "name": "baz"},
			}
			assert.Equal(t, expectedResults, actualResults)
		})

		type test1Record struct {
			ID   int
			Name string
		}

		t.Run("prepare", func(t *testing.T) {
			stmt, err := db.PrepareContext(ctx, "SELECT id, name FROM test1 WHERE id = ?")
			assert.NoError(t, err)
			defer stmt.Close()

			t.Run("QueryRowContext", func(t *testing.T) {
				foo := test1Record{}
				err := stmt.QueryRowContext(ctx, 1).Scan(&foo.ID, &foo.Name)
				assert.NoError(t, err)
				assert.Equal(t, test1Record{ID: 1, Name: "foo"}, foo)
			})
		})
	})

	t.Run("with tx", func(t *testing.T) {
		t.Run("rollback", func(t *testing.T) {
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			t.Run("update", func(t *testing.T) {
				r, err := tx.ExecContext(ctx, "UPDATE test1 SET name = ? WHERE id = ?", "quux", 3)
				assert.NoError(t, err)
				rowsAffected, err := r.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
			t.Run("rollback", func(t *testing.T) {
				err := tx.Rollback()
				assert.NoError(t, err)
			})
		})
		t.Run("commit", func(t *testing.T) {
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			t.Run("update", func(t *testing.T) {
				r, err := tx.ExecContext(ctx, "UPDATE test1 SET name = ? WHERE id = ?", "qux", 3)
				assert.NoError(t, err)
				rowsAffected, err := r.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
			t.Run("commit", func(t *testing.T) {
				err := tx.Commit()
				assert.NoError(t, err)
			})
		})

	})
}

func parseJsonLines(t *testing.T, b []byte) []map[string]interface{} {
	lines := bytes.Split(b, []byte("\n"))
	results := []map[string]interface{}{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		result := map[string]interface{}{}
		if err := json.Unmarshal(line, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		results = append(results, result)
	}
	return results
}
