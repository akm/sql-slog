package main_test

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/tests/testhelper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLowLevelWithContext(t *testing.T) {
	dsn := "./sqlite3_test.db"
	defer os.Remove(dsn)

	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	logs := testhelper.NewLogAssertion(buf)
	logger := slog.New(sqlslog.NewJSONHandler(buf, &slog.HandlerOptions{Level: sqlslog.LevelVerbose}))

	seqIdGen := testhelper.NewSeqIDGenerator()
	connIDKey := "conn_id"
	stmtIDKey := "stmt_id"

	db, err := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.Logger(logger),
		sqlslog.IDGenerator(seqIdGen.Generate),
		sqlslog.ConnIDKey(connIDKey),
		sqlslog.StmtIDKey(stmtIDKey),
		sqlslog.ConnBegin(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Begin Complete" }),
		sqlslog.ConnClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Close Complete" }),
		sqlslog.ConnPrepare(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Prepare Complete" }),
		sqlslog.ConnResetSession(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ResetSession Complete" }),
		sqlslog.ConnPing(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Ping Complete" }),
		sqlslog.ConnExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ExecContext Complete" }),
		sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.QueryContext Complete" }),
		sqlslog.ConnPrepareContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.PrepareContext Complete" }),
		sqlslog.ConnBeginTx(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.BeginTx Complete" }),
		sqlslog.ConnectorConnect(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Connector.Connect Complete" }),
		sqlslog.DriverOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.Open Complete" }),
		sqlslog.DriverOpenConnector(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.OpenConnector Complete" }),
		sqlslog.SqlslogOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "sqlslog.Open Complete" }),
		sqlslog.RowsClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Close Complete" }),
		sqlslog.RowsNext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Next Complete" }),
		sqlslog.RowsNextResultSet(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.NextResultSet Complete" }),
		sqlslog.StmtClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Close Complete" }),
		sqlslog.StmtExec(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Exec Complete" }),
		sqlslog.StmtQuery(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Query Complete" }),
		sqlslog.StmtExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.ExecContext Complete" }),
		sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.QueryContext Complete" }),
		sqlslog.TxCommit(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Commit Complete" }),
		sqlslog.TxRollback(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Rollback Complete" }),
	)
	require.NoError(t, err)
	defer db.Close()

	connIDExpected := seqIdGen.Next()

	t.Run("sqlslog.Open log", func(t *testing.T) {
		logs.Assert(t, []map[string]interface{}{
			{"level": "DEBUG", "msg": "sqlslog.Open Start", "driver": "sqlite3", "dsn": dsn},
			{"level": "INFO", "msg": "sqlslog.Open Complete", "driver": "sqlite3", "dsn": dsn},
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
		logs.Assert(t, []map[string]interface{}{})
	})

	t.Run("delete", func(t *testing.T) {
		query := "DELETE FROM test1"
		logs.Start()
		stmt, err := db.PrepareContext(ctx, query)
		assert.NoError(t, err)

		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected},
		})

		logs.Start()
		result, err := stmt.Exec()
		assert.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": "[]", connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": "[]", connIDKey: connIDExpected},
		})

		logs.Start()
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, rowsAffected, int64(0))
		logs.Assert(t, []map[string]interface{}{})

		logs.Start()
		stmt.Close()
		logs.Assert(t, []map[string]interface{}{
			{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected},
		})

	})

	t.Run("without tx", func(t *testing.T) {
		testData := []string{"foo", "bar", "baz"}
		for i, name := range testData {
			t.Run("insert "+name, func(t *testing.T) {
				query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
				args := fmt.Sprintf("[{Name: Ordinal:1 Value:%d} {Name: Ordinal:2 Value:%s}]", i+1, name)
				logs.Start()
				result, err := db.ExecContext(ctx, query, i+1, name)
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, connIDKey: connIDExpected},
				})

				logs.Start()
				rowsAffected, err := result.RowsAffected()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{})
				assert.Equal(t, int64(1), rowsAffected)
			})
		}

		t.Run("select count", func(t *testing.T) {
			t.Run("without condition", func(t *testing.T) {
				query := "SELECT count(*) FROM test1"
				var cnt int
				err := db.QueryRowContext(ctx, query).Scan(&cnt)
				assert.NoError(t, err)
				assert.Equal(t, 3, cnt)
			})
			t.Run("with condition", func(t *testing.T) {
				query := "SELECT count(*) FROM test1 WHERE name LIKE ?"
				var cnt int
				err := db.QueryRowContext(ctx, query, "ba%").Scan(&cnt)
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
			})

		})

		t.Run("select all", func(t *testing.T) {
			query := "SELECT id, name FROM test1"
			logs.Start()
			rows, err := db.QueryContext(ctx, query)
			assert.NoError(t, err)
			defer rows.Close()

			actualResults := []map[string]interface{}{}
			for rows.Next() {
				result := map[string]interface{}{}
				var id int
				var name string
				require.NoError(t, rows.Scan(&id, &name))
				result["id"] = id
				result["name"] = name
				actualResults = append(actualResults, result)
			}
			assert.Equal(t, []map[string]interface{}{
				{"id": 1, "name": "foo"},
				{"id": 2, "name": "bar"},
				{"id": 3, "name": "baz"},
			}, actualResults)
		})

		t.Run("select", func(t *testing.T) {
			query := "SELECT id, name FROM test1 WHERE name LIKE ?"
			logs.Start()
			rows, err := db.QueryContext(ctx, query, "ba%")
			assert.NoError(t, err)
			defer func() {
				logs.Start()
				assert.NoError(t, rows.Close())
				logs.Assert(t, []map[string]interface{}{
					// {"level": "TRACE", "msg": "Rows.Close Start"},
					// {"level": "DEBUG", "msg": "Rows.Close Complete"},
				})
			}()
			args := "[{Name: Ordinal:1 Value:ba%}]"
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.QueryContext Start", "query": query, "args": args, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.QueryContext Complete", "query": query, "args": args, connIDKey: connIDExpected},
			})

			t.Run("rows.Columns", func(t *testing.T) {
				logs.Start()
				columns, err := rows.Columns()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{})
				assert.Equal(t, []string{"id", "name"}, columns)
			})
			t.Run("rows", func(t *testing.T) {
				logs.Start()
				columnTypes, err := rows.ColumnTypes()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{})
				assert.Len(t, columnTypes, 2)
				t.Run("ColumnTypes[0]", func(t *testing.T) {
					logs.Start()
					ct := columnTypes[0]
					assert.Equal(t, "id", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(0), lengthValue)
					assert.False(t, lengthOK)
					assert.Equal(t, "INTEGER", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.True(t, nullableValue)
					assert.True(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "NullInt64", scanType.Name())
					logs.Assert(t, []map[string]interface{}{})
				})
				t.Run("ColumnTypes[1]", func(t *testing.T) {
					logs.Start()
					ct := columnTypes[1]
					assert.Equal(t, "name", ct.Name())
					lengthValue, lengthOK := ct.Length()
					assert.Equal(t, int64(0), lengthValue)
					assert.False(t, lengthOK)
					assert.Equal(t, "VARCHAR(255)", ct.DatabaseTypeName())
					dsPrecision, dsScale, dsOK := ct.DecimalSize()
					assert.Equal(t, int64(0), dsPrecision)
					assert.Equal(t, int64(0), dsScale)
					assert.False(t, dsOK)
					nullableValue, nullableOK := ct.Nullable()
					assert.True(t, nullableValue)
					assert.True(t, nullableOK)
					scanType := ct.ScanType()
					assert.Equal(t, "NullString", scanType.Name())
					logs.Assert(t, []map[string]interface{}{})
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
				require.NoError(t, rows.Scan(&id, &name))
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

		t.Run("PrepareContext", func(t *testing.T) {
			query := "SELECT id, name FROM test1 WHERE id = ?"
			logs.Start()
			stmt, err := db.PrepareContext(ctx, query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected},
			})

			defer func() {
				logs.Start()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected},
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
					{"level": "DEBUG", "msg": "Stmt.QueryContext Start", "args": "[{Name: Ordinal:1 Value:1}]", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Stmt.QueryContext Complete", "args": "[{Name: Ordinal:1 Value:1}]", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": false, connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Rows.Close Start", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Rows.Close Complete", connIDKey: connIDExpected},
				})
				assert.Equal(t, test1Record{ID: 1, Name: "foo"}, foo)
			})
		})

		t.Run("PrepareContext invalid query", func(t *testing.T) {
			query := "invalid select query"
			logs.Start()
			_, err := db.PrepareContext(ctx, query)
			assert.Error(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "ERROR", "msg": "Conn.PrepareContext Error", "query": query, "error": "near \"invalid\": syntax error", connIDKey: connIDExpected},
			})

		})

		t.Run("prepare + ExecContext", func(t *testing.T) {
			query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
			logs.Start()
			stmt, err := db.PrepareContext(ctx, query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected},
			})

			defer func() {
				logs.Start()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected},
				})
			}()

			t.Run("ExecContext", func(t *testing.T) {
				t.Run("success", func(t *testing.T) {
					logs.Start()
					result, err := stmt.ExecContext(ctx, 4, "qux")
					assert.NoError(t, err)
					logs.Assert(t, []map[string]interface{}{
						{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
						{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
						{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": "[{Name: Ordinal:1 Value:4} {Name: Ordinal:2 Value:qux}]", connIDKey: connIDExpected},
						{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": "[{Name: Ordinal:1 Value:4} {Name: Ordinal:2 Value:qux}]", connIDKey: connIDExpected},
					})
					rowsAffected, err := result.RowsAffected()
					assert.NoError(t, err)
					assert.Equal(t, int64(1), rowsAffected)
				})
				t.Run("error", func(t *testing.T) {
					logs.Start()
					_, err := stmt.ExecContext(ctx, "abc", "qux")
					assert.Error(t, err)
					args := "[{Name: Ordinal:1 Value:abc} {Name: Ordinal:2 Value:qux}]"
					logs.Assert(t, []map[string]interface{}{
						{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
						{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
						{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": args, connIDKey: connIDExpected},
						{"level": "ERROR", "msg": "Stmt.ExecContext Error", "args": args, "error": "datatype mismatch", connIDKey: connIDExpected},
					})
				})
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
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "qux", int64(3))
				args := "[{Name: Ordinal:1 Value:qux} {Name: Ordinal:2 Value:3}]"
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, connIDKey: connIDExpected},
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
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				logs.Start()
				r, err := tx.ExecContext(ctx, query, "quux", int64(3))
				args := "[{Name: Ordinal:1 Value:quux} {Name: Ordinal:2 Value:3}]"
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, connIDKey: connIDExpected},
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

	t.Run("Conn", func(t *testing.T) {
		logs.Start()
		conn, err := db.Conn(ctx)
		require.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
		})

		defer func() {
			logs.Start()
			err := conn.Close()
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{})
		}()

		t.Run("Raw", func(t *testing.T) {
			logs.Start()
			err := conn.Raw(func(driverConn interface{}) error {
				logs.Assert(t, []map[string]interface{}{})
				assert.Equal(t, "*sqlslog.connWithContextWrapper", fmt.Sprintf("%T", driverConn))
				if assert.Implements(t, (*driver.Conn)(nil), driverConn) {
					dConn := driverConn.(driver.Conn)

					var tx driver.Tx
					t.Run("Begin", func(t *testing.T) {
						logs.Start()
						var err error
						tx, err = dConn.Begin()
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Conn.Begin Start", connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Conn.Begin Complete", connIDKey: connIDExpected},
						})
					})

					t.Run("Prepare", func(t *testing.T) {
						stmtIDExpected := seqIdGen.Next()

						query := "SELECT id, name FROM test1 WHERE id = ?"
						logs.Start()
						stmt, err := dConn.Prepare(query)
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Conn.Prepare Start", "query": query, connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Conn.Prepare Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
						})

						defer func() {
							logs.Start()
							stmt.Close()
							logs.Assert(t, []map[string]interface{}{
								{"level": "DEBUG", "msg": "Stmt.Close Start", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
								{"level": "INFO", "msg": "Stmt.Close Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
							})
						}()

						t.Run("Query", func(t *testing.T) {
							rows, err := stmt.Query([]driver.Value{int64(1)})
							require.NoError(t, err)
							defer rows.Close()
						})
					})

					t.Run("prepare + Exec", func(t *testing.T) {
						stmtIDExpected := seqIdGen.Next()

						query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
						stmt, err := dConn.Prepare(query)
						assert.NoError(t, err)

						defer func() { assert.NoError(t, stmt.Close()) }()

						t.Run("Exec", func(t *testing.T) {
							t.Run("success", func(t *testing.T) {
								logs.Start()
								result, err := stmt.Exec([]driver.Value{4, "qux"})
								assert.NoError(t, err)
								args := "[4 qux]"
								logs.Assert(t, []map[string]interface{}{
									{"level": "DEBUG", "msg": "Stmt.Exec Start", "query": query, "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
									{"level": "INFO", "msg": "Stmt.Exec Complete", "query": query, "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
								})
								rowsAffected, err := result.RowsAffected()
								assert.NoError(t, err)
								assert.Equal(t, int64(1), rowsAffected)
							})
							t.Run("error", func(t *testing.T) {
								logs.Start()
								_, err := stmt.Exec([]driver.Value{"abc", "qux"})
								assert.Error(t, err)
								args := "[abc qux]"
								logs.Assert(t, []map[string]interface{}{
									{"level": "DEBUG", "msg": "Stmt.Exec Start", "query": query, "args": args, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
									{"level": "ERROR", "msg": "Stmt.Exec Error", "query": query, "args": args, "error": "datatype mismatch", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
								})
							})
						})
					})

					t.Run("Rollback", func(t *testing.T) {
						logs.Start()
						err := tx.Rollback()
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Tx.Rollback Start", connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Tx.Rollback Complete", connIDKey: connIDExpected},
						})
					})
				}

				return nil
			})
			assert.NoError(t, err)
		})
	})

}
