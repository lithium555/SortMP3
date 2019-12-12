package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/lithium555/SortMP3/models"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	albumYear = 1995
	albumName = "The Gallery"
	cover     = ""
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
	t.Run("trying to use CreatePostgresForTesting(t)", func(t *testing.T) {
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
	})

	t.Run("table album does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotRes, gotErr := db.AddAlbum(2, albumName, albumYear, cover)
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, 0, gotRes)
	})

	t.Run("add duplicate album", func(t *testing.T) {
		db := ensureTables(t)
		defer db.Close()

		authorID, err := db.AddAuthor("Soufly")
		assert.Nil(t, err)

		albumID, err := db.AddAlbum(authorID, albumName, albumYear, cover)
		assert.Equal(t, 1, albumID)
		assert.Nil(t, err)

		gotID, gotErr := db.AddAlbum(authorID, albumName, albumYear, cover)
		require.Equal(t, 1, gotID) // should return id of elements which exists
		require.Nil(t, gotErr)

		expectLen := 1
		res, err := db.SelectALBUM()
		require.Equal(t, expectLen, len(res))
		require.Nil(t, err)
	})
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

// ensureTables - ensures that tables exist and has the same structure as the test expects.
func ensureTables(t *testing.T) Database {
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

func TestDatabase_GetExistsAuthor(t *testing.T) {
	t.Run("find record which doesnt exist in database", func(t *testing.T) {
		db := ensureTables(t)
		defer db.Close()

		author := "System of a down" //"Dark Tranquillity"
		var authorID int
		err := db.PostgresConn.QueryRow(`
						INSERT INTO author (author_name)
						VALUES ($1) RETURNING id
			`, author).Scan(&authorID)
		require.Nil(t, err)

		gotVal, gotErr := db.GetExistsAuthor("Iwrestledabearonce")
		resErr := convertError(gotErr)
		require.Zero(t, gotVal)
		require.Equal(t, NotFoundErr, resErr) // no record in db

	})
	t.Run("table author doesnt exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotVal, gotErr := db.GetExistsAuthor("Iwrestledabearonce")
		require.Equal(t, 0, gotVal)
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
	})

	t.Run("success, author exists", func(t *testing.T) {
		db := ensureTables(t)
		defer db.Close()

		author := "Dark Tranquillity"
		var authorID int
		err := db.PostgresConn.QueryRow(`
						INSERT INTO author (author_name)
						VALUES ($1) RETURNING id
			`, author).Scan(&authorID)
		require.Nil(t, err)

		gotVal, gotErr := db.GetExistsAuthor(author)
		require.Equal(t, 1, gotVal)
		require.Nil(t, gotErr)
	})
}

func TestFindRecord(t *testing.T) {
	// TODO: 1.2) Обработчики ошибок. Разобраться как отличить типы ошибок: "нет записи",
	//  "конфликт по уникальному полю", "неверный FK".

	t.Run("no such record in table", func(t *testing.T) {
		db := ensureTables(t)

		var albumID int
		err := db.PostgresConn.QueryRow(`
			SELECT id FROM author WHERE author_name = $1;
		`, 555).Scan(&albumID)

		expectedID := 0
		require.Equal(t, expectedID, albumID)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("conflict by unique field", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := db.Drop("person")
		require.Nil(t, dropErr)

		_, insertErr := db.PostgresConn.Exec(`
			CREATE TABLE person (
			   id serial PRIMARY KEY,
			   first_name VARCHAR (50),
			   last_name VARCHAR (50),
			   email VARCHAR (50) UNIQUE
			);`)
		assert.Nil(t, insertErr)

		var personID int
		for i := 0; i < 3; i++ {
			insertErr := db.PostgresConn.QueryRow(`
				INSERT INTO person (first_name, last_name, email) 
				VALUES($1, $2, $3)
				RETURNING id`, "Slava", "Pinchuk", "development1810@gmail.com").Scan(&personID)

			if i > 0 {
				gotErr := convertError(insertErr)
				require.Equal(t, DuplicateValueErr, gotErr)
			}
		}
	})

	//- пытаются создать запись в неверным FK
	t.Run("create record with wrong Foreign Key", func(t *testing.T) {
		db := ensureTables(t)

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

		gotErr := convertError(err)
		require.Equal(t, WrongForeignKeyErr, gotErr)
	})
}

func Test_AddAuthor(t *testing.T) {
	//- пытаются внести дубликат по первичному ключу
	t.Run("insert duplicate by foreign key", func(t *testing.T) {
		db := ensureTables(t)
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
			if i == 0 {
				require.Nil(t, err)
			} else {
				require.Equal(t, DuplicateValueErr, convertError(err))
			}
		}
	})

	t.Run("add one record", func(t *testing.T) {
		db := ensureTables(t)

		gotID, gotErr := db.AddAuthor("Entombed")
		require.Nil(t, gotErr)
		require.NotNil(t, gotID)

		expectLen := 1
		expectAuthor := "Entombed"

		res, err := db.FindAuthors()
		require.Nil(t, err)
		require.Equal(t, expectLen, len(res))
		require.Equal(t, expectAuthor, res[0].AuthorName)
	})

	t.Run("add 2 records", func(t *testing.T) {
		db := ensureTables(t)

		metalBands := []string{"Sepultura", "Suicide Silence"}

		for _, v := range metalBands {
			id, err := db.AddAuthor(v)
			require.Nil(t, err)
			require.NotNil(t, id)
		}
		expectLen := 2

		authors, err := db.FindAuthors()
		require.Nil(t, err)
		require.Equal(t, metalBands[0], authors[0].AuthorName)
		require.Equal(t, metalBands[1], authors[1].AuthorName)
		require.Equal(t, expectLen, len(authors))
	})

	t.Run("table author does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotRes, gotErr := db.FindAuthors()
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, 0, len(gotRes))
	})

	t.Run("add duplicate author", func(t *testing.T) {
		db := ensureTables(t)

		authors := []string{"Entombed", "Entombed"}
		expectLength := 1

		id, err := db.AddAuthor(authors[0])
		assert.NotNil(t, id)
		assert.Nil(t, err)

		gotID, gotErr := db.AddAuthor(authors[1])
		assert.Nil(t, gotErr)
		assert.Equal(t, 1, gotID) // function should return id if element exist

		gotVal, gotErr := db.SelectAUTHOR()
		require.Nil(t, gotErr)
		require.Equal(t, expectLength, len(gotVal))
	})
}

func Test_AddGenre(t *testing.T) {
	t.Run("successful test", func(t *testing.T) {
		db := ensureTables(t)
		testGenres := []string{"Jazz", "Blues", "Metal", "PostRock"}

		for _, genre := range testGenres {
			gotID, gotErr := db.AddGenre(genre)
			require.Nil(t, gotErr)
			require.NotNil(t, gotID)
		}

		genres, err := db.FindGenres()
		require.Nil(t, err)
		expectLength := 4

		for k, genre := range genres {
			require.Equal(t, testGenres[k], genre.GenreName)
		}
		require.Equal(t, expectLength, len(genres))
	})

	t.Run("duplicate name of genre", func(t *testing.T) {
		db := ensureTables(t)
		testGenres := []string{"Jazz", "Jazz"}

		for k, genre := range testGenres {
			gotID, gotErr := db.AddGenre(genre)
			require.Nil(t, gotErr)
			require.NotNil(t, gotID)

			if k == 1 {
				require.Equal(t, 1, gotID) // if elemet exist, func should return exists ID
			}
		}

		res, err := db.SelectGENRE()
		require.Nil(t, err)
		require.Equal(t, 1, len(res))
	})

	// TODO: what to do with empty values - maybe write parser before insert?
	t.Run("empty genre string", func(t *testing.T) {
		db := ensureTables(t)
		gotID, gotErr := db.AddGenre("")
		//require.NotNil(t, gotErr)
		require.Nil(t, gotErr)
		require.NotNil(t, gotID)

		genres, err := db.FindGenres()
		require.Nil(t, err)
		require.Equal(t, 1, len(genres))
		//require.Equal(t, 0, len(genres))
	})

	t.Run("add one genreName 2 times, expect one record in database", func(t *testing.T) {
		db := ensureTables(t)
		genreName := "Classic"

		var genreID int
		err := db.GetConnection().QueryRow(`
		INSERT INTO GENRE(genre_name) 
			VALUES ($1) 	
			RETURNING id
		`, genreName).Scan(&genreID)
		assert.Nil(t, err)

		gotID, gotErr := db.AddGenre(genreName)
		require.NotNil(t, gotID)
		require.Equal(t, genreID, gotID) // return existID
		require.Nil(t, gotErr)

		expectLen := 1
		genres, findErr := db.FindGenres()
		require.Nil(t, findErr)
		require.Equal(t, expectLen, len(genres))
		require.Equal(t, genreName, genres[0].GenreName)
	})
}

func Test_FindGenres(t *testing.T) {
	t.Run("empty table", func(t *testing.T) {
		db := ensureTables(t)

		gotRes, gotErr := db.FindGenres()
		require.Nil(t, gotErr)
		require.Equal(t, 0, len(gotRes))
	})

	t.Run("table genre does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotRes, gotErr := db.FindGenres()
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, 0, len(gotRes))
	})
}

func TestDatabase_InsertSONG(t *testing.T) {
	t.Run("table SONG does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotErr := db.InsertSONG("Punish my Heaven", 55, 23, 5, 1)
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
	})
	t.Run("wrong foreign key", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		gotErr := db.InsertSONG("Punish my Heaven", 55, 23, 5, 1)
		require.Equal(t, WrongForeignKeyErr, convertError(gotErr))
	})
	t.Run("successful insert one song", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID, err := db.AddAuthor("Dark Tranquillity")
		assert.Nil(t, err)

		genreID, err := db.AddGenre("melodic-death metal")
		assert.Nil(t, err)

		numberOfTrack := 23

		albumID, err := db.AddAlbum(authorID, "The Gallery", 1995, "")
		assert.Nil(t, err)

		gotErr := db.InsertSONG("Punish my Heaven", albumID, genreID, authorID, numberOfTrack)
		require.Nil(t, gotErr)
	})
	t.Run("try to insert the same song twice", func(t *testing.T) {
		db := ensureTables(t)

		authorID, err := db.AddAuthor("Dark Tranquillity")
		assert.Nil(t, err)

		genreID, err := db.AddGenre("melodic-death metal")
		assert.Nil(t, err)

		numberOfTrack := 23

		albumID, err := db.AddAlbum(authorID, "The Gallery", 1995, "")
		assert.Nil(t, err)

		gotErr := db.InsertSONG("Punish my Heaven", albumID, genreID, authorID, numberOfTrack)
		require.Nil(t, gotErr)

		gotErr = db.InsertSONG("Punish my Heaven", albumID, genreID, authorID, numberOfTrack)
		require.Equal(t, DuplicateValueErr, convertError(gotErr))

		expectSong := &models.Song{
			SongID:      1,
			NameOfSong:  "Punish my Heaven",
			AlbumID:     1,
			GenreID:     1,
			AuthorID:    1,
			TrackNumber: 23,
		}
		expectLen := 1 // test is fine, we expect one element, when we have duplicate

		gotVal, gotErr := db.SelectSONG()
		require.Equal(t, expectLen, len(gotVal))
		require.Equal(t, expectSong, gotVal[0])
		require.Nil(t, gotErr)
	})
	t.Run("insert 3 songs", func(t *testing.T) {

	})
}

func Test_SelectSONG(t *testing.T) {
	t.Run("table song does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		gotVal, gotErr := db.SelectSONG()
		require.Equal(t, 0, len(gotVal))
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
	})
	t.Run("success select", func(t *testing.T) {
		db := ensureTables(t)

		authorID, err := db.AddAuthor("Dark Tranquillity")
		assert.Nil(t, err)

		genreID, err := db.AddGenre("melodic-death metal")
		assert.Nil(t, err)

		numberOfTrack := 5

		albumID, err := db.AddAlbum(authorID, "The Gallery", 1995, "")
		assert.Nil(t, err)

		gotErr := db.InsertSONG("The Gallery", albumID, genreID, authorID, numberOfTrack)
		assert.Nil(t, gotErr)

		expectSong := &models.Song{
			SongID:      1,
			NameOfSong:  "The Gallery",
			AlbumID:     1,
			GenreID:     1,
			AuthorID:    1,
			TrackNumber: 5,
		}
		expectLen := 1

		gotVal, gotErr := db.SelectSONG()
		require.Nil(t, gotErr)
		require.Equal(t, expectLen, len(gotVal))
		require.Equal(t, expectSong, gotVal[0])
	})
}

func Test_SelectGENRE(t *testing.T) {
	t.Run("table genre does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		expectLength := 0

		gotVal, gotErr := db.SelectGENRE()
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, expectLength, len(gotVal))
	})
	t.Run("empty table genre", func(t *testing.T) {
		db := ensureTables(t)

		expectLength := 0

		gotVal, gotErr := db.SelectGENRE()
		require.Nil(t, gotErr)
		require.Equal(t, expectLength, len(gotVal))
	})
	t.Run("successful select genre", func(t *testing.T) {
		db := ensureTables(t)

		expectLength := 4

		genres := []string{"classic", "blues", "jazz", "post-rock"}

		for _, genre := range genres {
			id, err := db.AddGenre(genre)
			assert.Nil(t, err)
			assert.NotNil(t, id)
		}

		gotVal, gotErr := db.SelectGENRE()
		require.Nil(t, gotErr)
		require.Equal(t, expectLength, len(gotVal))
	})
}

func Test_SelectAUTHOR(t *testing.T) {
	t.Run("table author does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		expectLength := 0

		gotVal, gotErr := db.SelectAUTHOR()
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, expectLength, len(gotVal))
	})
	t.Run("success select author", func(t *testing.T) {
		db := ensureTables(t)

		expectLength := 4
		authors := []string{"Entombed", "Definition Sane", "Bob Dilan", "Trombone shorty"}

		for _, author := range authors {
			authorID, err := db.AddAuthor(author)
			assert.NotNil(t, authorID)
			assert.Nil(t, err)
		}

		gotVal, gotErr := db.SelectAUTHOR()
		require.Nil(t, gotErr)
		require.Equal(t, expectLength, len(gotVal))
	})
}

func Test_SelectALBUM(t *testing.T) {
	t.Run("table album does not exist", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		expectLength := 0

		gotVal, gotErr := db.SelectALBUM()
		require.Equal(t, TableDoesntExistErr, convertError(gotErr))
		require.Equal(t, expectLength, len(gotVal))
	})
	t.Run("empty table", func(t *testing.T) {
		db := ensureTables(t)

		expectLen := 0

		gotVal, gotErr := db.SelectALBUM()
		require.Equal(t, expectLen, len(gotVal))
		require.Nil(t, gotErr)
	})
	t.Run("successful select album", func(t *testing.T) {
		db := ensureTables(t)

		authorID, err := db.AddAuthor("Soufly")
		assert.Nil(t, err)

		albumID, err := db.AddAlbum(authorID, albumName, albumYear, cover)
		assert.Equal(t, 1, albumID)
		assert.Nil(t, err)

		authorID2, err2 := db.AddAuthor("Behemoth")
		assert.Nil(t, err2)

		albumID2, err2 := db.AddAlbum(authorID2, albumName, albumYear, cover)
		fmt.Printf("albumID2 = '%v'\n", albumID2)
		assert.Equal(t, 2, albumID2)
		assert.Nil(t, err)

		expectLen := 2

		gotVal, gotErr := db.SelectALBUM()
		require.Nil(t, gotErr)
		require.Equal(t, expectLen, len(gotVal))
	})
}
