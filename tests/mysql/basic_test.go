package mysqltest

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestBasic(t *testing.T) {
	dbName := "app1"

	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		t.Fatal(err)
	}
	defer exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run()

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/"+dbName)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		t.Logf("Ping failed: %v", err)
		time.Sleep(2 * time.Second)
	}

	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test (id INT PRIMARY KEY, name VARCHAR(255))"); err != nil {
		t.Fatal(err)
	}
}
