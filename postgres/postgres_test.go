package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
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
	allTables := []string{TableSong, TableAlbum, TableAuthor, TableGenre}
	for _, table := range allTables {
		if err := getPostgres.Drop(table); err != nil {
			return err
		}
	}
	return nil
}

func workWithTables(t *testing.T) Database {
	db, err := GetPostgresConnection()
	assert.Nil(t, err)

	dropErr := DropTablesAfterTest(db)
	assert.Nil(t, dropErr)

	createErr := CreateTablesForTest(db)
	assert.Nil(t, createErr)

	return db
}

//сделай себе несколько тестов которые:
//- пытаются найти запись которой нет
func TestFindRecord(t *testing.T) {

	var (
		albumYear = 1995
		albumName = "The Gallery"
		cover     = ""
	)
	t.Run("find record which doesnt exist in database", func(t *testing.T) {
		db := workWithTables(t)
		defer db.Close()

		author := "System of a down" //"Dark Tranquillity"
		var authorID int
		err := db.PostgresConn.QueryRow(`
						INSERT INTO author (author_name)
						VALUES ($1) RETURNING id
			`, author).Scan(&authorID)
		require.Nil(t, err)

		//expectError := "sql: no rows in result set"

		gotVal, gotErr := db.GetExistsAuthor("Iwrestledabearonce")
		require.Zero(t, gotVal)
		//require.Equal(t, expectError, gotErr.Error())
		if errors.Is(err, os.ErrExist) {
			//TODO: what should we do with errors from convertError func?
			valErr := convertError(gotErr)
			fmt.Printf("Error from 'erros.go': = '%v'\n", valErr)
		}
		//TODO: what should we do with errors IS NOT FROM convertError() func?
	})

	//- пытаются внести дубликат по первичному ключу
	t.Run("insert duplicate by foreign key", func(t *testing.T) {
		db := workWithTables(t)
		defer db.Close()

		authorID, err := db.AddAuthor("Dark Tranquillity")
		assert.Nil(t, err)

		for i := 0; i < 2; i++ {
			var albumID int
			err = db.PostgresConn.QueryRow(`
			INSERT INTO ALBUM(author_id, album_name, album_year, cover)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, authorID, albumName, albumYear, cover).Scan(&albumID)
			require.Nil(t, err)
		}
	})

	// TODO: 1.2) Обработчики ошибок. Разобраться как отличить типы ошибок: "нет записи",
	//  "конфликт по уникальному полю", "неверный FK".

	t.Run("no such record in table", func(t *testing.T) {
		db := workWithTables(t)

		var albumID int
		err := db.PostgresConn.QueryRow(`
			SELECT id FROM author WHERE author_name = $1;
		`, 555).Scan(&albumID)

		expectError := "sql: no rows in result set"
		require.Equal(t, expectError, err.Error())
	})

	t.Run("conflict by unique field", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := db.Drop("person")
		require.Nil(t, dropErr)

		_, err = db.PostgresConn.Exec(`
			CREATE TABLE person (
			   id serial PRIMARY KEY,
			   first_name VARCHAR (50),
			   last_name VARCHAR (50),
			   email VARCHAR (50) UNIQUE
			);`)
		assert.Nil(t, err)

		//expectedErr := "pq: duplicate key value violates unique constraint"

		var personID int
		for i := 0; i < 3; i++ {
			insertErr := db.PostgresConn.QueryRow(`
				INSERT INTO person (first_name, last_name, email) 
				VALUES($1, $2, $3)
				RETURNING id`, "Slava", "Pinchuk", "development1810@gmail.com").Scan(&personID)

			if i > 0 {
				//fmt.Printf("insertErr.Error() = '%v'\n", insertErr.Error())
				//require.Contains(t, insertErr.Error(), expectedErr)
				//fmt.Printf("Conver = '%v'\n", convertError(insertErr))
				if errors.Is(err, os.ErrExist) {
					//TODO: what should we do with errors from convertError func?
					valErr := convertError(insertErr)
					fmt.Printf("Error from 'erros.go': = '%v'\n", valErr)
				}
			}
		}
	})

	//- пытаются создать запись в неверным FK
	t.Run("create record with wrong Foreign Key", func(t *testing.T) {
		db := workWithTables(t)

		metalBands := []string{"Sepultura", "Suicide Silence"}

		for _, v := range metalBands {
			_, err := db.AddAuthor(v)
			assert.Nil(t, err)
		}

		var albumID int
		err := db.PostgresConn.QueryRow(`
			INSERT INTO ALBUM(author_id, album_name, album_year, cover)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, 555, albumName, albumYear, cover).Scan(&albumID)
		if errors.Is(err, os.ErrExist) {
			fmt.Printf("Conver = '%v'\n", convertError(err))
		}
	})
}

func ErrorHandler(err error) {
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(">>>>>>> pq error-Code:", err.Code)
		fmt.Println(">>>>>>> pq error.Code.Name():", err.Code.Name())
	}
}

func Test_AddAuthor(t *testing.T) {
	t.Run("add one record", func(t *testing.T) {
		db := workWithTables(t)

		gotID, gotErr := db.AddAuthor("Entombed")
		require.Nil(t, gotErr)
		require.NotNil(t, gotID)
	})

	t.Run("add 2 records", func(t *testing.T) {
		db := workWithTables(t)

		metalBands := []string{"Sepultura", "Suicide Silence"}

		for _, v := range metalBands {
			id, err := db.AddAuthor(v)
			require.Nil(t, err)
			require.NotNil(t, id)
		}
	})
}
