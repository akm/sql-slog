package opts

// IDGen is a function that generates an ID string.
type IDGen = func() string

const defaultIDLength = 16

// IDGeneratorDefault is the default ID generator.
var IDGeneratorDefault = NewChaCha8IDGenerator(defaultIDLength).Generate
