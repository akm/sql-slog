package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// IDGen is a function that generates an ID string.
type IDGen = opts.IDGen

// IDGenerator returns an Option that sets the ID generator.
// The default is IDGeneratorDefault.
func IDGenerator(idGen IDGen) Option { return opts.IDGenerator(idGen) }

// ConnIDKey sets the key for the connection ID.
// The default is ConnIDKeyDefault.
func ConnIDKey(key string) Option { return opts.ConnIDKey(key) }

// TxIDKey sets the key for the transaction ID.
// The default is TxIDKeyDefault.
func TxIDKey(key string) Option { return opts.TxIDKey(key) }

// StmtIDKey sets the key for the statement ID.
// The default is StmtIDKeyDefault.
func StmtIDKey(key string) Option { return opts.StmtIDKey(key) }
