package postgres

import (
	"database/sql"
	"log"
	"testing"

	"github.com/ory/dockertest"
)

func CreatePostgresForTesting(t testing.TB) (*sql.DB, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}

	cont, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_PASSWORD=postgres",
	})
	if err != nil {
		t.Fatal(err)
	}

	const port = "5432/tcp"
	addr := `postgres://postgres:postgres@` + cont.GetHostPort(port)

	err = pool.Retry(func() error {
		cli, err := sql.Open("postgres", addr)
		if err != nil {
			return err
		}
		defer cli.Close()
		return cli.Ping()
	})
	if err != nil {
		cont.Close()
		t.Fatal(err)
	}

	connPostgres, err := sql.Open("postgres", addr)
	if err != nil {
		cont.Close()
		t.Fatal(err)
	}

	return connPostgres, func() {
		cont.Close()
	}
}

func TestGetConnection(t *testing.T) {
	conn, kill := CreatePostgresForTesting(t)
	kill()

	log.Printf("conn = '%v'\n", conn)
}
