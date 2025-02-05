package sqlslog

import "math/rand/v2"

// IDGen is a function that generates an ID string.
type IDGen = func() string

// IDGenerator returns an Option that sets the ID generator.
// The default is IDGeneratorDefault.
func IDGenerator(idGen IDGen) Option {
	return func(o *options) {
		o.DriverOptions.IDGen = idGen
		o.DriverOptions.ConnOptions.IDGen = idGen
	}
}

const (
	ConnIDKeyDefault = "conn_id"
	TxIDKeyDefault   = "tx_id"
	StmtIDKeyDefault = "stmt_id"
)

// ConnIDKey sets the key for the connection ID.
// The default is ConnIDKeyDefault.
func ConnIDKey(key string) Option {
	return func(o *options) {
		o.DriverOptions.ConnIDKey = key
	}
}

// TxIDKey sets the key for the transaction ID.
// The default is TxIDKeyDefault.
func TxIDKey(key string) Option {
	return func(o *options) { o.DriverOptions.ConnOptions.TxIDKey = key }
}

// StmtIDKey sets the key for the statement ID.
// The default is StmtIDKeyDefault.
func StmtIDKey(key string) Option {
	return func(o *options) { o.DriverOptions.ConnOptions.StmtIDKey = key }
}

// Returns a random ID generator that generates a string of length characters
// using randInt to generate random integers such as Int function from math/rand/v2 package.
func RandIntIDGenerator(
	randInt func() int,
	letters []byte,
	length int,
) IDGen {
	lenLetters := len(letters)
	return func() string {
		b := make([]byte, length)
		for i := range b {
			b[i] = letters[randInt()%lenLetters]
		}
		return string(b)
	}
}

// Returns a random ID generator that generates a string of length characters
// using randRead to generate random bytes such as Read function from crypto/rand package.
func RandReadIDGenerator(
	randRead func(b []byte) (n int, err error),
	letters []byte,
	length int,
) func() (string, error) {
	lenLetters := len(letters)
	return func() (string, error) {
		b := make([]byte, length)
		if _, err := randRead(b); err != nil {
			return "", err
		}
		for i := range b {
			b[i] = letters[int(b[i])%lenLetters]
		}
		return string(b), nil
	}
}

// IDGenErrorSuppressor returns an ID generator that suppresses errors.
// If an error occurs, the recover function is called with the error and the result is returned.
func IDGenErrorSuppressor(idGen func() (string, error), recover func(error) string) IDGen {
	return func() string {
		id, err := idGen()
		if err != nil {
			return recover(err)
		}
		return id
	}
}

const defaultIDLength = 16

var defaultIDLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")

// IDGeneratorDefault is the default ID generator.
var IDGeneratorDefault = RandIntIDGenerator(rand.Int, defaultIDLetters, defaultIDLength)
