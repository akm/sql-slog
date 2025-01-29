package sqlslog

import (
	"github.com/akm/sql-slog/sqlslogopts"
)

type IDGen = sqlslogopts.IDGen

var IDGeneratorDefault = sqlslogopts.IDGeneratorDefault

const (
	ConnIDKeyDefault = sqlslogopts.ConnIDKeyDefault
	TxIDKeyDefault   = sqlslogopts.TxIDKeyDefault
	StmtIDKeyDefault = sqlslogopts.StmtIDKeyDefault
)
