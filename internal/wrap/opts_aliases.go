package wrap

import "github.com/akm/sql-slog/internal/opts"

type (
	Option      = opts.Option
	Options     = opts.Options
	StepOptions = opts.StepOptions
)

var (
	NewOptions       = opts.NewOptions
	NewTextHandler   = opts.NewTextHandler
	DurationAttrFunc = opts.DurationAttrFunc

	Logger = opts.Logger
)
