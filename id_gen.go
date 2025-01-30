package sqlslog

import "github.com/akm/sql-slog/opts"

// IDGen is a function that generates an ID string.
type IDGen = opts.IDGen

// IDGenerator returns an Option that sets the ID generator.
// The default is IDGeneratorDefault.
func IDGenerator(idGen IDGen) Option { return func(o *options) { o.idGen = idGen } }

const defaultIDLength = 16

// IDGeneratorDefault is the default ID generator.
var IDGeneratorDefault = NewChaCha8IDGenerator(defaultIDLength).Generate

const (
	ConnIDKeyDefault = "conn_id"
	TxIDKeyDefault   = "tx_id"
	StmtIDKeyDefault = "stmt_id"
)

// ConnIDKey sets the key for the connection ID.
// The default is ConnIDKeyDefault.
func ConnIDKey(key string) Option { return func(o *options) { o.connIDKey = key } }

// TxIDKey sets the key for the transaction ID.
// The default is TxIDKeyDefault.
func TxIDKey(key string) Option { return func(o *options) { o.txIDKey = key } }

// StmtIDKey sets the key for the statement ID.
// The default is StmtIDKeyDefault.
func StmtIDKey(key string) Option { return func(o *options) { o.stmtIDKey = key } }
