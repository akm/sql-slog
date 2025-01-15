module sqlite3-test

go 1.22.10

require (
	github.com/akm/sql-slog v0.0.0-20250111025848-713ab89fa0bf
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/akm/sql-slog => ../..