package sqlslog

import (
	"github.com/akm/sql-slog/public"
)

type IDGen = public.IDGen

var IDGeneratorDefault = public.IDGeneratorDefault

const (
	ConnIDKeyDefault = public.ConnIDKeyDefault
	TxIDKeyDefault   = public.TxIDKeyDefault
	StmtIDKeyDefault = public.StmtIDKeyDefault
)
