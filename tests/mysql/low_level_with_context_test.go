package main_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/tests/testhelper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLowLevelWithContext(t *testing.T) {
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
	logs := testhelper.NewLogAssertion(buf)
	logs.Start()

	seqIdGen := testhelper.NewSeqIDGenerator()
	connIDKey := "conn_id"
	stmtIDKey := "stmt_id"
	txIDKey := "tx_id"
	connIDExpected := seqIdGen.Next()

	db, _, err := sqlslog.Open(ctx, "mysql", "root@tcp(localhost:3306)/"+dbName,
		append(
			testhelper.StepEventMsgOptions,
			sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
			sqlslog.LogWriter(buf),
			sqlslog.LogLevel(sqlslog.LevelVerbose),
			sqlslog.IDGenerator(seqIdGen.Generate),
			sqlslog.ConnIDKey(connIDKey),
			sqlslog.StmtIDKey(stmtIDKey),
			sqlslog.TxIDKey(txIDKey),
		)...,
	)
	require.NoError(t, err)
	defer db.Close()

	t.Run("sqlslog.Open log", func(t *testing.T) {
		logs.Assert(t, []map[string]interface{}{
			{"level": "DEBUG", "msg": "sqlslog.Open Start", "driver": "mysql", "dsn": dsn},
			{"level": "DEBUG", "msg": "Driver.OpenConnector Start", "dsn": dsn},
			{"level": "INFO", "msg": "Driver.OpenConnector Complete", "dsn": dsn, connIDKey: connIDExpected},
			{"level": "INFO", "msg": "sqlslog.Open Complete", "driver": "mysql", "dsn": dsn},
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
		query := "CREATE TABLE IF NOT EXISTS test1 (id INT PRIMARY KEY, name VARCHAR(255))"
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
				stmtIDExpected := seqIdGen.Next()
				query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
				args := fmt.Sprintf("[{Name: Ordinal:1 Value:%d} {Name: Ordinal:2 Value:%s}]", i+1, name)
				logs.Start()
				result, err := db.ExecContext(ctx, query, i+1, name)
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, "skip": true, connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": "INSERT INTO test1 (id, name) VALUES (?, ?)", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": "INSERT INTO test1 (id, name) VALUES (?, ?)", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})

				logs.Start()
				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				logs.AssertEmpty(t)
				assert.Equal(t, int64(1), rowsAffected)
			})
		}

		t.Run("select", func(t *testing.T) {
			stmtIDExpected := seqIdGen.Next()
			query := "SELECT id, name FROM test1 WHERE name LIKE ?"
			logs.Start()
			rows, err := db.QueryContext(ctx, query, "ba%")
			assert.NoError(t, err)
			defer func() {
				logs.Start()
				assert.NoError(t, rows.Close())
				logs.AssertEmpty(t) // Rows.Close and Stmt.Close are called from rows.Next when EOF
			}()
			args := "[{Name: Ordinal:1 Value:ba%}]"
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.QueryContext Start", "query": query, "args": args, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.QueryContext Complete", "query": query, "args": args, "skip": true, connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": "SELECT id, name FROM test1 WHERE name LIKE ?", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": "SELECT id, name FROM test1 WHERE name LIKE ?", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "DEBUG", "msg": "Stmt.QueryContext Start", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "INFO", "msg": "Stmt.QueryContext Complete", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
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
					logs.AssertEmpty(t)
				})
				t.Run("ColumnTypes[1]", func(t *testing.T) {
					logs.Start()
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
					logs.AssertEmpty(t)
				})
			})

			logs.Start()
			actualResults := []map[string]interface{}{}
			for rows.Next() {
				logs.Assert(t, []map[string]interface{}{
					{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": false, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
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
				{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": true, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "TRACE", "msg": "Rows.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "DEBUG", "msg": "Rows.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
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
			query := "SELECT id, name FROM test1 WHERE id = ?"
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
				err := stmt.QueryRowContext(ctx, 1).Scan(&foo.ID, &foo.Name)
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
			query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
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
				result, err := stmt.ExecContext(ctx, 4, "qux")
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
			txIDExpected := seqIdGen.Next()
			logs.Start()
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				stmtIDExpected := seqIdGen.Next()
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "qux", 3)
				args := "[{Name: Ordinal:1 Value:qux} {Name: Ordinal:2 Value:3}]"
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, "skip": true, connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
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
					{"level": "DEBUG", "msg": "Tx.Rollback Start", connIDKey: connIDExpected, txIDKey: txIDExpected},
					{"level": "INFO", "msg": "Tx.Rollback Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
				})
			})
		})
		t.Run("commit", func(t *testing.T) {
			txIDExpected := seqIdGen.Next()
			logs.Start()
			tx, err := db.BeginTx(ctx, nil)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				stmtIDExpected := seqIdGen.Next()
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "quux", 3)
				args := "[{Name: Ordinal:1 Value:quux} {Name: Ordinal:2 Value:3}]"
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, "skip": true, connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
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
					{"level": "DEBUG", "msg": "Tx.Commit Start", connIDKey: connIDExpected, txIDKey: txIDExpected},
					{"level": "INFO", "msg": "Tx.Commit Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
				})
			})
		})
	})
}
