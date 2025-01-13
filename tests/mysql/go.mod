module mysql-test

go 1.23.2

require github.com/go-sql-driver/mysql v1.8.1

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/akm/sql-slog v0.0.0-20250111025848-713ab89fa0bf // indirect
)

replace github.com/akm/sql-slog => ../..
