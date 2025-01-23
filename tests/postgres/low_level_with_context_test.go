package main_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/tests/testhelper"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLowLevelWithContext(t *testing.T) {
	dbName := "app1"
	dbPort := 5432
	dsn := fmt.Sprintf("host=127.0.0.1 port=%d user=root password=password dbname=%s sslmode=disable", dbPort, dbName)

	os.Setenv("POSTGRES_PORT", fmt.Sprintf("%d", dbPort))
	os.Setenv("POSTGRES_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("DEBUG") == "" {
		defer exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run()
	}

	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	logs := testhelper.NewLogAssertion(buf)
	logger := slog.New(sqlslog.NewJSONHandler(buf, &slog.HandlerOptions{Level: sqlslog.LevelVerbose}))

	seqIdGen := testhelper.NewSeqIDGenerator()
	connIDKey := "conn_id"
	stmtIDKey := "stmt_id"

	db, err := sqlslog.Open(ctx, "postgres", dsn,
		sqlslog.Logger(logger),
		sqlslog.IDGenerator(seqIdGen.Generate),
		sqlslog.ConnIDKey(connIDKey),
		sqlslog.StmtIDKey(stmtIDKey),
	)
	require.NoError(t, err)
	defer db.Close()

	connIDExpected := seqIdGen.Next()

	t.Run("sqlslog.Open log", func(t *testing.T) {
		logs.Assert(t, []map[string]interface{}{
			{"level": "DEBUG", "msg": "sqlslog.Open Start", "driver": "postgres", "dsn": dsn},
			{"level": "INFO", "msg": "sqlslog.Open Complete", "driver": "postgres", "dsn": dsn},
		})
	})

	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		t.Logf("Ping failed: %v", err)
		time.Sleep(2 * time.Second)
	}

	t.Run("Ping", func(t *testing.T) {
		logs.Start()
		err := db.PingContext(ctx)
		assert.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "VERBOSE", "msg": "Conn.Ping Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.Ping Complete", connIDKey: connIDExpected},
		})
	})

	t.Run("create table", func(t *testing.T) {
		query := "CREATE TABLE IF NOT EXISTS test1 (id INTEGER PRIMARY KEY, name VARCHAR(255))"
		logs.Start()
		result, err := db.ExecContext(ctx, query)
		assert.NoError(t, err)
		t.Logf("buf.String(): %s\n", buf.String())
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": "[]", connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": "[]", connIDKey: connIDExpected},
		})

		logs.Start()
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected)
		logs.AssertEmpty(t)
	})

	t.Run("delete", func(t *testing.T) {
		stmtIDExpected := seqIdGen.Next()
		query := "DELETE FROM test1"
		logs.Start()
		stmt, err := db.PrepareContext(ctx, query)
		assert.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

		logs.Start()
		result, err := stmt.Exec()
		assert.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": "[]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": "[]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

		logs.Start()
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, rowsAffected, int64(0))
		logs.AssertEmpty(t)

		logs.Start()
		stmt.Close()
		logs.Assert(t, []map[string]interface{}{
			{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

	})

	t.Run("without tx", func(t *testing.T) {
		testData := []string{"foo", "bar", "baz"}
		for i, name := range testData {
			t.Run("insert "+name, func(t *testing.T) {
				query := "INSERT INTO test1 (id, name) VALUES ($1,$2);"
				logs.Start()
				result, err := db.ExecContext(ctx, query, int64(i+1), name)
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": fmt.Sprintf("[{Name: Ordinal:1 Value:%d} {Name: Ordinal:2 Value:%s}]", i+1, name), connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": fmt.Sprintf("[{Name: Ordinal:1 Value:%d} {Name: Ordinal:2 Value:%s}]", i+1, name), connIDKey: connIDExpected},
				})

				logs.Start()
				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				logs.AssertEmpty(t)
				assert.Equal(t, int64(1), rowsAffected)
			})
		}

		t.Run("select", func(t *testing.T) {
			query := "SELECT id, name FROM test1 WHERE name LIKE $1"
			logs.Start()
			rows, err := db.QueryContext(ctx, query, "ba%")
			assert.NoError(t, err)
			defer func() {
				logs.Start()
				assert.NoError(t, rows.Close())
				logs.Assert(t, []map[string]interface{}{
					// Rows.Close and Stmt.Close are called from rows.Next when EOF
				})
			}()
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.QueryContext Start", "query": query, "args": "[{Name: Ordinal:1 Value:ba%}]", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.QueryContext Complete", "query": query, "args": "[{Name: Ordinal:1 Value:ba%}]", connIDKey: connIDExpected},
			})

			t.Run("rows.Columns", func(t *testing.T) {
				logs.Start()
				columns, err := rows.Columns()
				assert.NoError(t, err)
				logs.AssertEmpty(t)
				assert.Equal(t, []string{"id", "name"}, columns)
			})
			t.Run("rows", func(t *testing.T) {
				logs.Start()
				columnTypes, err := rows.ColumnTypes()
				assert.NoError(t, err)
				logs.AssertEmpty(t)
				assert.Len(t, columnTypes, 2)
				t.Run("ColumnTypes[0]", func(t *testing.T) {
					logs.Start()
					ct := columnTypes[0]
					assert.Equal(t, "id", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(0), lengthValue)
					assert.False(t, lengthOK)
					assert.Equal(t, "INT4", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.False(t, nullableValue)
					assert.False(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "int32", scanType.Name())
					logs.AssertEmpty(t)
				})
				t.Run("ColumnTypes[1]", func(t *testing.T) {
					logs.Start()
					ct := columnTypes[1]
					assert.Equal(t, "name", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(255), lengthValue)
					assert.True(t, lengthOK)
					assert.Equal(t, "VARCHAR", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.False(t, nullableValue)
					assert.False(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "string", scanType.Name())
					logs.AssertEmpty(t)
				})
			})

			logs.Start()
			actualResults := []map[string]interface{}{}
			for rows.Next() {
				logs.Assert(t, []map[string]interface{}{
					{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": false, connIDKey: connIDExpected},
				})
				logs.Start()

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

			logs.Assert(t, []map[string]interface{}{
				{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": true, connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Rows.Close Start", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Rows.Close Complete", connIDKey: connIDExpected},
			})

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
			stmtIDExpected := seqIdGen.Next()
			query := "SELECT id, name FROM test1 WHERE id = $1"
			logs.Start()
			stmt, err := db.PrepareContext(ctx, query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			})

			defer func() {
				logs.Start()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
			}()

			t.Run("QueryRowContext", func(t *testing.T) {
				logs.Start()
				foo := test1Record{}
				err := stmt.QueryRowContext(ctx, int64(1)).Scan(&foo.ID, &foo.Name)
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Stmt.QueryContext Start", "args": "[{Name: Ordinal:1 Value:1}]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.QueryContext Complete", "args": "[{Name: Ordinal:1 Value:1}]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": false, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "TRACE", "msg": "Rows.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Rows.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
				assert.Equal(t, test1Record{ID: 1, Name: "foo"}, foo)
			})
		})

		t.Run("prepare + ExecContext", func(t *testing.T) {
			stmtIDExpected := seqIdGen.Next()
			query := "INSERT INTO test1 (id, name) VALUES ($1, $2)"
			logs.Start()
			stmt, err := db.PrepareContext(ctx, query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			})

			defer func() {
				logs.Start()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
			}()

			t.Run("ExecContext", func(t *testing.T) {
				logs.Start()
				result, err := stmt.ExecContext(ctx, int64(4), "qux")
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": "[{Name: Ordinal:1 Value:4} {Name: Ordinal:2 Value:qux}]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": "[{Name: Ordinal:1 Value:4} {Name: Ordinal:2 Value:qux}]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
		})

	})

	t.Run("with tx", func(t *testing.T) {
		t.Run("rollback", func(t *testing.T) {
			logs.Start()
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				query := "UPDATE test1 SET name = $1 WHERE id = $2"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "qux", int64(3))
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": "[{Name: Ordinal:1 Value:qux} {Name: Ordinal:2 Value:3}]", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": "[{Name: Ordinal:1 Value:qux} {Name: Ordinal:2 Value:3}]", connIDKey: connIDExpected},
				})

				rowsAffected, err := r.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
			t.Run("rollback", func(t *testing.T) {
				logs.Start()
				err := tx.Rollback()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Tx.Rollback Start", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Tx.Rollback Complete", connIDKey: connIDExpected},
				})
			})
		})
		t.Run("commit", func(t *testing.T) {
			logs.Start()
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				query := "UPDATE test1 SET name = $1 WHERE id = $2"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "quux", int64(3))
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": "[{Name: Ordinal:1 Value:quux} {Name: Ordinal:2 Value:3}]", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": "[{Name: Ordinal:1 Value:quux} {Name: Ordinal:2 Value:3}]", connIDKey: connIDExpected},
				})

				rowsAffected, err := r.RowsAffected()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
			t.Run("commit", func(t *testing.T) {
				logs.Start()
				err := tx.Commit()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Tx.Commit Start", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Tx.Commit Complete", connIDKey: connIDExpected},
				})
			})
		})
	})
}
