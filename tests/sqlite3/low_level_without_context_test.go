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
	"github.com/akm/sql-slog/opts"
	"github.com/akm/sql-slog/tests/testhelper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLowLevelWithoutContext(t *testing.T) {
	dsn := "./sqlite3_test.db"
	defer os.Remove(dsn)

	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	logs := testhelper.NewLogAssertion(buf)
	logger := slog.New(opts.NewJSONHandler(buf, &slog.HandlerOptions{Level: opts.LevelVerbose}))

	seqIdGen := testhelper.NewSeqIDGenerator()
	connIDKey := "conn_id"
	stmtIDKey := "stmt_id"
	txIDKey := "tx_id"

	db, err := sqlslog.Open(ctx, "sqlite3", dsn,
		append(
			testhelper.StepEventMsgOptions,
			opts.Logger(logger),
			opts.IDGenerator(seqIdGen.Generate),
			opts.ConnIDKey(connIDKey),
			opts.StmtIDKey(stmtIDKey),
			opts.TxIDKey(txIDKey),
		)...,
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
		if err := db.Ping(); err == nil {
			break
		}
		t.Logf("Ping failed: %v", err)
		time.Sleep(2 * time.Second)
	}

	t.Run("Ping", func(t *testing.T) {
		buf.Reset()
		err := db.Ping()
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
		buf.Reset()
		result, err := db.Exec(query)
		assert.NoError(t, err)
		t.Logf("buf.String(): %s\n", buf.String())
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": "[]", connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": "[]", connIDKey: connIDExpected},
		})

		buf.Reset()
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected)
		logs.Assert(t, []map[string]interface{}{})
	})

	t.Run("delete", func(t *testing.T) {
		stmtIDExpected := seqIdGen.Next()
		query := "DELETE FROM test1"
		buf.Reset()
		stmt, err := db.Prepare(query)
		assert.NoError(t, err)

		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

		buf.Reset()
		result, err := stmt.Exec()
		assert.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Stmt.ExecContext Start", "args": "[]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			{"level": "INFO", "msg": "Stmt.ExecContext Complete", "args": "[]", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

		buf.Reset()
		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, rowsAffected, int64(0))
		logs.Assert(t, []map[string]interface{}{})

		buf.Reset()
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
				query := "INSERT INTO test1 (id, name) VALUES (?, ?)"
				args := fmt.Sprintf("[{Name: Ordinal:1 Value:%d} {Name: Ordinal:2 Value:%s}]", i+1, name)
				buf.Reset()
				result, err := db.Exec(query, i+1, name)
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
					{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Conn.ExecContext Start", "query": query, "args": args, connIDKey: connIDExpected},
					{"level": "INFO", "msg": "Conn.ExecContext Complete", "query": query, "args": args, connIDKey: connIDExpected},
				})

				buf.Reset()
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
				err := db.QueryRow(query).Scan(&cnt)
				assert.NoError(t, err)
				assert.Equal(t, 3, cnt)
			})
			t.Run("with condition", func(t *testing.T) {
				query := "SELECT count(*) FROM test1 WHERE name LIKE ?"
				var cnt int
				err := db.QueryRow(query, "ba%").Scan(&cnt)
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
			})
		})

		t.Run("select all", func(t *testing.T) {
			query := "SELECT id, name FROM test1"
			buf.Reset()
			rows, err := db.Query(query)
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
			buf.Reset()
			rows, err := db.Query(query, "ba%")
			assert.NoError(t, err)
			defer func() {
				buf.Reset()
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
				buf.Reset()
				columns, err := rows.Columns()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{})
				assert.Equal(t, []string{"id", "name"}, columns)
			})
			t.Run("rows", func(t *testing.T) {
				buf.Reset()
				columnTypes, err := rows.ColumnTypes()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{})
				assert.Len(t, columnTypes, 2)
				t.Run("ColumnTypes[0]", func(t *testing.T) {
					buf.Reset()
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
					buf.Reset()
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

			buf.Reset()
			actualResults := []map[string]interface{}{}
			for rows.Next() {
				logs.Assert(t, []map[string]interface{}{
					{"level": "TRACE", "msg": "Rows.Next Start", connIDKey: connIDExpected},
					{"level": "DEBUG", "msg": "Rows.Next Complete", "eof": false, connIDKey: connIDExpected},
				})
				buf.Reset()

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
			stmtIDExpected := seqIdGen.Next()
			query := "SELECT id, name FROM test1 WHERE id = ?"
			buf.Reset()
			stmt, err := db.Prepare(query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			})

			defer func() {
				buf.Reset()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
			}()

			t.Run("QueryRowContext", func(t *testing.T) {
				buf.Reset()
				foo := test1Record{}
				err := stmt.QueryRow(int64(1)).Scan(&foo.ID, &foo.Name)
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
			buf.Reset()
			stmt, err := db.Prepare(query)
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
			})

			defer func() {
				buf.Reset()
				assert.NoError(t, stmt.Close())
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
					{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
				})
			}()

			t.Run("ExecContext", func(t *testing.T) {
				buf.Reset()
				result, err := stmt.Exec(4, "qux")
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
			buf.Reset()
			tx, err := db.Begin()
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				buf.Reset()
				r, err := tx.Exec(query, "qux", int64(3))
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
				buf.Reset()
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
			buf.Reset()
			tx, err := db.Begin()
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{
				{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
				{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
				{"level": "DEBUG", "msg": "Conn.BeginTx Start", connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Conn.BeginTx Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
			})

			t.Run("update", func(t *testing.T) {
				query := "UPDATE test1 SET name = ? WHERE id = ?"
				buf.Reset()
				r, err := tx.Exec(query, "quux", int64(3))
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
				buf.Reset()
				err := tx.Commit()
				assert.NoError(t, err)
				logs.Assert(t, []map[string]interface{}{
					{"level": "DEBUG", "msg": "Tx.Commit Start", connIDKey: connIDExpected, txIDKey: txIDExpected},
					{"level": "INFO", "msg": "Tx.Commit Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
				})
			})
		})
	})

	t.Run("Conn", func(t *testing.T) {
		buf.Reset()
		conn, err := db.Conn(ctx)
		require.NoError(t, err)
		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
		})

		defer func() {
			buf.Reset()
			err := conn.Close()
			assert.NoError(t, err)
			logs.Assert(t, []map[string]interface{}{})
		}()

		t.Run("Raw", func(t *testing.T) {
			buf.Reset()
			err := conn.Raw(func(driverConn interface{}) error {
				logs.Assert(t, []map[string]interface{}{})
				assert.Equal(t, "*wrap.connWithContextWrapper", fmt.Sprintf("%T", driverConn))
				if assert.Implements(t, (*driver.Conn)(nil), driverConn) {
					dConn := driverConn.(driver.Conn)

					var tx driver.Tx
					txIDExpected := seqIdGen.Next()

					t.Run("Begin", func(t *testing.T) {
						buf.Reset()
						var err error
						tx, err = dConn.Begin()
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Conn.Begin Start", connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Conn.Begin Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
						})
					})

					t.Run("Prepare and Query", func(t *testing.T) {
						stmtIDExpected := seqIdGen.Next()

						query := "SELECT id, name FROM test1 WHERE id = ?"
						buf.Reset()
						stmt, err := dConn.Prepare(query)
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Conn.Prepare Start", "query": query, connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Conn.Prepare Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
						})

						defer func() {
							buf.Reset()
							stmt.Close()
							logs.Assert(t, []map[string]interface{}{
								{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
								{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
							})
						}()

						t.Run("Query", func(t *testing.T) {
							rows, err := stmt.Query([]driver.Value{int64(1)})
							require.NoError(t, err)
							defer rows.Close()
						})
					})

					t.Run("Prepare and Exec", func(t *testing.T) {
						stmtIDExpected := seqIdGen.Next()

						query := "DELETE FROM test1 WHERE id = ?"
						buf.Reset()
						stmt, err := dConn.Prepare(query)
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Conn.Prepare Start", "query": query, connIDKey: connIDExpected},
							{"level": "INFO", "msg": "Conn.Prepare Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
						})

						defer func() {
							buf.Reset()
							stmt.Close()
							logs.Assert(t, []map[string]interface{}{
								{"level": "DEBUG", "msg": "Stmt.Close Start", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
								{"level": "INFO", "msg": "Stmt.Close Complete", connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
							})
						}()

						t.Run("Exec", func(t *testing.T) {
							result, err := stmt.Exec([]driver.Value{int64(1)})
							require.NoError(t, err)
							rowsAffected, err := result.RowsAffected()
							require.NoError(t, err)
							assert.Equal(t, int64(1), rowsAffected)
						})
					})

					t.Run("Rollback", func(t *testing.T) {
						buf.Reset()
						err := tx.Rollback()
						require.NoError(t, err)
						logs.Assert(t, []map[string]interface{}{
							{"level": "DEBUG", "msg": "Tx.Rollback Start", connIDKey: connIDExpected, txIDKey: txIDExpected},
							{"level": "INFO", "msg": "Tx.Rollback Complete", connIDKey: connIDExpected, txIDKey: txIDExpected},
						})
					})

				}

				return nil
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Stmt.ExecContext", func(t *testing.T) {
		stmtIDExpected := seqIdGen.Next()

		query := "DELETE FROM test1 WHERE id = ?"
		buf.Reset()
		stmt, err := db.Prepare(query)
		require.NoError(t, err)

		defer func() {
			buf.Reset()
			assert.NoError(t, stmt.Close())
			logs.Assert(t, []map[string]interface{}{
				{"level": "DEBUG", "msg": "Stmt.Close Start", stmtIDKey: stmtIDExpected, connIDKey: connIDExpected},
				{"level": "INFO", "msg": "Stmt.Close Complete", stmtIDKey: stmtIDExpected, connIDKey: connIDExpected},
			})
		}()

		logs.Assert(t, []map[string]interface{}{
			{"level": "VERBOSE", "msg": "Conn.ResetSession Start", connIDKey: connIDExpected},
			{"level": "TRACE", "msg": "Conn.ResetSession Complete", connIDKey: connIDExpected},
			{"level": "DEBUG", "msg": "Conn.PrepareContext Start", "query": query, connIDKey: connIDExpected},
			{"level": "INFO", "msg": "Conn.PrepareContext Complete", "query": query, connIDKey: connIDExpected, stmtIDKey: stmtIDExpected},
		})

		t.Run("Stmt.ExecContext", func(t *testing.T) {
			buf.Reset()
			result, err := stmt.ExecContext(ctx, int64(2))
			require.NoError(t, err)
			rowAffected, err := result.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowAffected)
		})
	})
}
