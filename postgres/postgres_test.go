package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/require"
)

func CreatePostgresForTesting(t testing.TB) (*sql.DB, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}

	cont, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=master", "POSTGRES_DB=musicDB"})
	if err != nil {
		t.Fatal(err)
	}

	const port = "5432/tcp"
	addr := cont.GetHostPort(port)
	fmt.Printf("addr = '%v'\n", addr)
	// Examples of connections:
	// db, err := sql.Open("postgres", "user=test password=test dbname=test sslmode=disable")
	// db, err := sql.Open("postgres", "postgres://username:password@localhost/db_name?sslmode=disable")

	err = pool.Retry(func() error {

		//connStr := "user=sorter password=master host=localhost port=5432 dbname=musicDB sslmode=disable"

		cli, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:master@%s/%s?sslmode=disable", addr, "musicDB"))
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

func TestDatabase_AddAlbum(t *testing.T) {
	t.Skip()
	conn, kill := CreatePostgresForTesting(t)
	kill()

	testDB := Database{
		PostgresConn: conn,
	}

	errCreate := testDB.CreateTable(CreateTableALBUM)
	fmt.Printf("errCreate = '%v'\n", errCreate)
	fmt.Printf("errCreate.Error() ='%v'\n", errCreate.Error())

	require.Nil(t, errCreate)
	//gotID, gotErr := testDB.AddAlbum()

}

// https://postgresql.leopard.in.ua/

//
//сделай себе несколько тестов которые:
//- пытаются найти запись которой нет
//- пытаются внести дубликат по первичному ключу
//- пытаются создать запись в неверным FK
//и попробуй распознать эти ошибки
//
//func Example() {
//	db, err := sql.Open("postgres",
//		"host=localhost dbname=Test sslmode=disable user=postgres password=secret")
//	if err != nil {
//		log.Fatal("cannot connect ...")
//	}
//	defer db.Close()
//	db.Exec(`set search_path='mySchema'`)
//
//	rows, err := db.Query(`select blah,blah2 from myTable`)
//
//}

func CreateTablesForTest(getPostgres Database) error {
	for _, query := range []string{
		CreateTableGENRE,
		CreateTableAUTHOR,
		CreateTableALBUM,
		CreateTableSONG,
	} {
		if err := getPostgres.CreateTable(query); err != nil {
			log.Printf("Can`t  '%v', error: '%v'\n", query, err)
			return err
		}
	}

	return nil
}

func DropTablesAfterTest(getPostgres Database) error {
	for _, query := range []string{
		CreateTableGENRE,
		CreateTableAUTHOR,
		CreateTableALBUM,
		CreateTableSONG,
	} {
		if err := getPostgres.CreateTable(query); err != nil {
			log.Printf("Can`t  '%v', error: '%v'\n", query, err)
			return err
		}
	}
	return nil
}

// docker run --name sort_music -e POSTGRES_PASSWORD=master -e POSTGRES_DB=musicDB -e POSTGRES_USER=sorter -p 5432:5432 -d postgres
func TestFindRecord(t *testing.T) {
	db, err := GetPostgresConnection()
	assert.Nil(t, err)

	dropErr := DropTablesAfterTest(db)
	assert.Nil(t, dropErr)

	defer db.Close()

	createErr := CreateTablesForTest(db)
	assert.Nil(t, createErr)

	author := "Dark Tranquillity"
	var authorID int
	err = db.PostgresConn.QueryRow(`
					INSERT INTO author (author_name)
					VALUES ($1) RETURNING id
		`, author).Scan(&authorID)
	require.Nil(t, err)
	fmt.Println("New record authorID is:", authorID)

}
