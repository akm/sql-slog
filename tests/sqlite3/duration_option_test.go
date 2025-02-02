package main_test

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/tests/testhelper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	dsn := "./sqlite3_test.db"
	defer os.Remove(dsn)

	ctx := context.TODO()

	type testCase struct {
		durationKey  string
		durationType sqlslog.DurationType
		assertion    func(t *testing.T, logs *testhelper.LogsAssertion)
	}

	findCompleteLog := func(t *testing.T, logs *testhelper.LogsAssertion) map[string]interface{} {
		for _, line := range logs.JsonLines(t) {
			if msg, ok := line["msg"].(string); ok && strings.Contains(msg, " Complete") {
				return line
			}
		}
		return nil
	}
	assertFloat64Duration := func(key string) func(t *testing.T, log map[string]interface{}) {
		return func(t *testing.T, log map[string]interface{}) {
			require.NotNil(t, log)
			require.NotNil(t, log[key])
			assert.IsType(t, float64(0), log[key]) // Value must be a float64 because of JSON unmarshaling
			assert.GreaterOrEqual(t, log[key], float64(0))
		}
	}

	testCases := []testCase{
		{
			durationKey:  "d",
			durationType: sqlslog.DurationNanoSeconds,
			assertion: func(t *testing.T, logs *testhelper.LogsAssertion) {
				assertFloat64Duration("d")(t, findCompleteLog(t, logs))
			},
		},
		{
			durationKey:  "duration-μs",
			durationType: sqlslog.DurationMicroSeconds,
			assertion: func(t *testing.T, logs *testhelper.LogsAssertion) {
				assertFloat64Duration("duration-μs")(t, findCompleteLog(t, logs))
			},
		},
		{
			durationKey:  "duration-msec",
			durationType: sqlslog.DurationMilliSeconds,
			assertion: func(t *testing.T, logs *testhelper.LogsAssertion) {
				assertFloat64Duration("duration-msec")(t, findCompleteLog(t, logs))
			},
		},
		{
			durationKey:  "duration-of-go",
			durationType: sqlslog.DurationGoDuration,
			assertion: func(t *testing.T, logs *testhelper.LogsAssertion) {
				assertFloat64Duration("duration-of-go")(t, findCompleteLog(t, logs))
			},
		},
		{
			durationKey:  "duration-string",
			durationType: sqlslog.DurationString,
			assertion: func(t *testing.T, logs *testhelper.LogsAssertion) {
				log := findCompleteLog(t, logs)
				require.NotNil(t, log)
				assert.Regexp(t, `^[\d.]+(s|ms|µs|ns)$`, log["duration-string"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.durationKey, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			logs := testhelper.NewLogAssertion(buf)
			logger := slog.New(sqlslog.NewJSONHandler(buf, &slog.HandlerOptions{Level: sqlslog.LevelVerbose}))
			db, err := sqlslog.Open(ctx, "sqlite3", dsn,
				append(
					testhelper.StepEventMsgOptions,
					sqlslog.Logger(logger),
					sqlslog.DurationKey(tc.durationKey),
					sqlslog.Duration(tc.durationType),
				)...,
			)
			require.NoError(t, err)
			defer db.Close()

			tc.assertion(t, logs)
		})
	}
}
