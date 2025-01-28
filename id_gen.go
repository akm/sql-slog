package sqlslog

// IDGenerator returns an Option that sets the ID generator.
// The default is IDGeneratorDefault.
func IDGenerator(idGen IDGen) Option { return func(o *Options) { o.IdGen = idGen } }

// ConnIDKey sets the key for the connection ID.
// The default is ConnIDKeyDefault.
func ConnIDKey(key string) Option { return func(o *Options) { o.ConnIDKey = key } }

// TxIDKey sets the key for the transaction ID.
// The default is TxIDKeyDefault.
func TxIDKey(key string) Option { return func(o *Options) { o.TxIDKey = key } }

// StmtIDKey sets the key for the statement ID.
// The default is StmtIDKeyDefault.
func StmtIDKey(key string) Option { return func(o *Options) { o.StmtIDKey = key } }
