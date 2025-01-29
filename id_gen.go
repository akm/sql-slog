package sqlslog

import (
	sqlslogopts "github.com/akm/sql-slog/opts"
)

type IDGen = sqlslogopts.IDGen

var IDGeneratorDefault = sqlslogopts.IDGeneratorDefault

const (
	ConnIDKeyDefault = sqlslogopts.ConnIDKeyDefault
	TxIDKeyDefault   = sqlslogopts.TxIDKeyDefault
	StmtIDKeyDefault = sqlslogopts.StmtIDKeyDefault
)

var (
	IDGenerator = sqlslogopts.IDGenerator
	ConnIDKey   = sqlslogopts.ConnIDKey
	TxIDKey     = sqlslogopts.TxIDKey
	StmtIDKey   = sqlslogopts.StmtIDKey
)
