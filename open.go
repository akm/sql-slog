package sqlslog

import (
	"context"
	"database/sql"

	"github.com/akm/sql-slog/internal/wrap"
)

func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, error) {
	return wrap.Open(ctx, driverName, dsn, opts...)
}
