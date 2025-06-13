module mysql-test

go 1.23.2

require (
	github.com/akm/sql-slog v0.1.3
	github.com/akm/sql-slog/tests/testhelper v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
	github.com/stretchr/testify v1.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/akm/sql-slog => ../..

replace github.com/akm/sql-slog/tests/testhelper => ../testhelper
