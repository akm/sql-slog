package mysqltest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	dbPort := 3306
	dsn := fmt.Sprintf("root@tcp(localhost:%d)/%s", dbPort, dbName)

	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("DEBUG") == "" {
		defer exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run()
	}

	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := sqlslog.Open(ctx, "mysql", "root@tcp(localhost:3306)/"+dbName, logger)
	require.NoError(t, err)

	t.Run("sqlslog.Open log", func(t *testing.T) {
		actualEntries := parseJsonLines(t, buf.Bytes())
		exptectedEntries := []map[string]interface{}{
			{"level": "DEBUG", "msg": "sqlslog.Open Start", "driver": "mysql", "dsn": dsn},
			{"level": "DEBUG", "msg": "OpenConnector Start", "dsn": dsn},
			{"level": "INFO", "msg": "OpenConnector Complete", "dsn": dsn},
			{"level": "INFO", "msg": "sqlslog.Open Complete", "driver": "mysql", "dsn": dsn},
		}
		assertMapSlice(t, exptectedEntries, actualEntries, "time")
	})

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
			{"level": "DEBUG", "msg": "ResetSession Start"},
			{"level": "INFO", "msg": "ResetSession Complete"},
			{"level": "DEBUG", "msg": "Ping Start"},
			{"level": "INFO", "msg": "Ping Complete"},
		}
		assertMapSlice(t, exptectedEntries, actualEntries, "time")
	})

	t.Run("create table", func(t *testing.T) {
		buf.Reset()
		result, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test1 (id INT PRIMARY KEY, name VARCHAR(255))")
		assert.NoError(t, err)
		t.Logf("buf.String(): %s\n", buf.String())
		actualEntries := parseJsonLines(t, buf.Bytes())
		exptectedEntries := []map[string]interface{}{
			{"level": "DEBUG", "msg": "ResetSession Start"},
			{"level": "INFO", "msg": "ResetSession Complete"},
			{"level": "DEBUG", "msg": "ExecContext Start"},
			{"level": "INFO", "msg": "ExecContext Complete"},
		}
		assertMapSlice(t, exptectedEntries, actualEntries, "time")

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected)
	})

	t.Run("delete", func(t *testing.T) {
		stmt, err := db.Prepare("DELETE FROM test1")
		assert.NoError(t, err)
		defer stmt.Close()
		result, err := stmt.Exec()
		assert.NoError(t, err)
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, rowsAffected, int64(0))
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

			t.Run("rows.Columns", func(t *testing.T) {
				columns, err := rows.Columns()
				assert.NoError(t, err)
				assert.Equal(t, []string{"id", "name"}, columns)
			})
			t.Run("rows", func(t *testing.T) {
				columnTypes, err := rows.ColumnTypes()
				assert.NoError(t, err)
				assert.Len(t, columnTypes, 2)
				t.Run("ColumnTypes[0]", func(t *testing.T) {
					ct := columnTypes[0]
					assert.Equal(t, "id", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(0), lengthValue)
					assert.False(t, lengthOK)
					assert.Equal(t, "INT", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.False(t, nullableValue)
					assert.True(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "int32", scanType.Name())
				})
				t.Run("ColumnTypes[1]", func(t *testing.T) {
					ct := columnTypes[1]
					assert.Equal(t, "name", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(0), lengthValue)
					assert.False(t, lengthOK)
					assert.Equal(t, "VARCHAR", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.True(t, nullableValue)
					assert.True(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "NullString", scanType.Name())
				})
			})

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

func assertMapSlice(t *testing.T, expected, actual []map[string]interface{}, ignoredFields ...string) {
	t.Helper()
	wellFormedActual := []map[string]interface{}{}
	for _, a := range actual {
		wellFormed := map[string]interface{}{}
		for k, v := range a {
			if contains(ignoredFields, k) {
				continue
			}
			wellFormed[k] = v
		}
		wellFormedActual = append(wellFormedActual, wellFormed)
	}
	assert.Equal(t, expected, wellFormedActual)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
