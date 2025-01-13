module mysql-test

go 1.21.13

require (
	github.com/akm/sql-slog v0.0.0-20250111025848-713ab89fa0bf
	github.com/go-sql-driver/mysql v1.8.1
	github.com/stretchr/testify v1.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/akm/sql-slog => ../..
