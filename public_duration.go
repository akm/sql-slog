package sqlslog

import (
	sqlslogopts "github.com/akm/sql-slog/opts"
)

type DurationType = sqlslogopts.DurationType

const (
	DurationNanoSeconds  DurationType = sqlslogopts.DurationNanoSeconds
	DurationMicroSeconds DurationType = sqlslogopts.DurationMicroSeconds
	DurationMilliSeconds DurationType = sqlslogopts.DurationMilliSeconds
	DurationGoDuration   DurationType = sqlslogopts.DurationGoDuration
	DurationString       DurationType = sqlslogopts.DurationString
)

const DurationKeyDefault = sqlslogopts.DurationKeyDefault
