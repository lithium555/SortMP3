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
		expectedID := id

		gotID, gotErr := db.AddAuthor(authors[1])
		assert.Nil(t, gotErr)
		assert.Equal(t, expectedID, gotID) // function should return id if element exist

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
				require.Equal(t, 1, gotID) // if element exist, func should return exists ID
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
		require.Equal(t, gotID, len(genres))
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
	t.Run("insert 2 different songs, where the same author, the same genre, the same album", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		_, err = db.AddAuthor("Animals as Leaders")
		assert.Nil(t, err)
		autorID, err := db.AddAuthor("Here Comes The Kraken")
		assert.Nil(t, err)

		_, err = db.AddGenre("Math-metal")
		assert.Nil(t, err)
		_, err = db.AddGenre("progressive")
		assert.Nil(t, err)
		genreID, genreErr := db.AddGenre("melodic death metal")
		assert.Nil(t, genreErr)

		albumID, albumErr := db.AddAlbum(autorID, "Here Comes The Kraken", 2009, "")
		assert.Nil(t, albumErr)

		insErr := db.InsertSONG("It's Comming", albumID, genreID, autorID, 1)
		assert.Nil(t, insErr)

		gotErr := db.InsertSONG("Don't Fail Me Darko", albumID, genreID, autorID, 2)
		require.Nil(t, gotErr)
	})
	t.Run("Trying to write Duplicate Value to the DB", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		_, err = db.AddAuthor("Animals as Leaders")
		assert.Nil(t, err)
		authorID, err := db.AddAuthor("Here Comes The Kraken")
		assert.Nil(t, err)

		_, err = db.AddGenre("Math-metal")
		assert.Nil(t, err)
		_, err = db.AddGenre("progressive")
		assert.Nil(t, err)
		genreID, genreErr := db.AddGenre("melodic death metal")
		assert.Nil(t, genreErr)

		albumID, albumErr := db.AddAlbum(authorID, "Here Comes The Kraken", 2009, "")
		assert.Nil(t, albumErr)

		insErr := db.InsertSONG("Name of fire", albumID, genreID, authorID, 1)
		assert.Nil(t, insErr)

		gotErr := db.InsertSONG("Name of fire", albumID, genreID, authorID, 3)
		require.Equal(t, DuplicateValueErr, gotErr)
	})
	t.Run("authors with the same album", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("Definition Sane")
		assert.Nil(t, err)
		authorID2, err := db.AddAuthor("Metallica")
		assert.Nil(t, err)

		albumID1, albumErr := db.AddAlbum(authorID1, "Saint Anger", 2003, "")
		require.Nil(t, albumErr)

		albumID2, albumErr := db.AddAlbum(authorID2, "Saint Anger", 2003, "")
		require.Nil(t, albumErr)

		// authorID1 != authorID2, so we should have 2 records with different albumID:
		require.NotEqual(t, albumID1, albumID2)
	})
	t.Run("different authors, the same song, the same album, the same trackNum", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("Definition Sane")
		assert.Nil(t, err)
		authorID2, err := db.AddAuthor("Metallica")
		assert.Nil(t, err)

		genreID1, err := db.AddGenre("Grind metal")
		assert.Nil(t, err)
		genreID2, err := db.AddGenre("Trash metal")
		assert.Nil(t, err)

		albumID1, albumErr := db.AddAlbum(authorID1, "Saint Anger", 2003, "")
		require.Nil(t, albumErr)

		nameOfSong := "Saint Anger"
		insErr := db.InsertSONG(nameOfSong, albumID1, genreID1, authorID1, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG(nameOfSong, albumID1, genreID2, authorID2, 1)
		require.Nil(t, insErr)
	})

	t.Run("different authors, different genre, the same song, the same album, the same trackNum, different genre", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("Definition Sane")
		assert.Nil(t, err)
		authorID2, err := db.AddAuthor("Metallica")
		assert.Nil(t, err)

		genreID1, err := db.AddGenre("Grind metal")
		assert.Nil(t, err)
		genreID2, err := db.AddGenre("Trash metal")
		assert.Nil(t, err)

		albumID, albumErr := db.AddAlbum(authorID1, "Saint Anger", 2003, "")
		require.Nil(t, albumErr)

		nameOfSong := "No Build to last"
		insErr := db.InsertSONG(nameOfSong, albumID, genreID1, authorID1, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG(nameOfSong, albumID, genreID2, authorID2, 1)
		require.Nil(t, insErr)
	})
	t.Run("2 albums of one author", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID, err := db.AddAuthor("KONVENT")
		assert.Nil(t, err)

		albumID1, albumErr1 := db.AddAlbum(authorID, "Demo", 2017, "")
		assert.Nil(t, albumErr1)
		albumID2, albumErr2 := db.AddAlbum(authorID, "Puritan Masochism", 20, "")
		assert.Nil(t, albumErr2)

		// 2 albums of one author, so we should get 2 different albumID
		require.NotEqual(t, albumID1, albumID2)
	})
	t.Run("the same name of song in different albums in one author", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("KONVENT")
		assert.Nil(t, err)

		genreID1, err := db.AddGenre("Death/Doom Metal")
		assert.Nil(t, err)

		albumID1, albumErr1 := db.AddAlbum(authorID1, "Demo", 2017, "")
		assert.Nil(t, albumErr1)
		albumID2, albumErr2 := db.AddAlbum(authorID1, "Puritan Masochism", 2020, "")
		assert.Nil(t, albumErr2)

		insErr := db.InsertSONG("Domination", albumID1, genreID1, authorID1, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Domination", albumID2, genreID1, authorID1, 10)
		require.Nil(t, insErr)
	})
	// Test with the same song '"Squares"' in different albums of one author.
	t.Run("2 authors with few albums, insert few different songs", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("KONVENT")
		assert.Nil(t, err)
		genreID1, err := db.AddGenre("Death/Doom Metal")
		assert.Nil(t, err)

		albumID1, albumErr1 := db.AddAlbum(authorID1, "Demo", 2017, "")
		assert.Nil(t, albumErr1)
		albumID2, albumErr2 := db.AddAlbum(authorID1, "Puritan Masochism", 2020, "")
		assert.Nil(t, albumErr2)

		insErr := db.InsertSONG("Chernobyl Child", albumID1, genreID1, authorID1, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Squares", albumID1, genreID1, authorID1, 2)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("No End", albumID1, genreID1, authorID1, 4)
		require.Nil(t, insErr)

		insErr = db.InsertSONG("Squares", albumID2, genreID1, authorID1, 2)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Bridge", albumID2, genreID1, authorID1, 5)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Idle Hands", albumID2, genreID1, authorID1, 7)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Ropes, Pt. 1", albumID2, genreID1, authorID1, 8)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Ropes, Pt. 2", albumID2, genreID1, authorID1, 9)
		require.Nil(t, insErr)

		// Second author
		authorID2, err := db.AddAuthor("Entombed")
		assert.Nil(t, err)
		genreID2, err := db.AddGenre("Death-n-Roll")
		assert.Nil(t, err)

		albumID3, albumErr3 := db.AddAlbum(authorID2, "Left Hand Path", 1990, "")
		assert.Nil(t, albumErr3)
		albumID4, albumErr4 := db.AddAlbum(authorID2, "Clandestine", 1991, "")
		assert.Nil(t, albumErr4)
		albumID5, albumErr5 := db.AddAlbum(authorID2, "Wolverine Blues", 1993, "")
		assert.Nil(t, albumErr5)
		albumID6, albumErr6 := db.AddAlbum(authorID2, "DCLXVI: To Ride Shoot Straight and Speak the Truth", 1997, "")
		assert.Nil(t, albumErr6)
		albumID7, albumErr7 := db.AddAlbum(authorID2, "Same Difference", 1998, "")
		assert.Nil(t, albumErr7)

		insErr = db.InsertSONG("Left Hand Path", albumID3, genreID2, authorID2, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Revel in Flesh", albumID3, genreID2, authorID2, 3)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("When Life Has Ceased", albumID3, genreID2, authorID2, 4)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Morbid Devourment", albumID3, genreID2, authorID2, 8)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Abnormally Deceased", albumID3, genreID2, authorID2, 9)
		require.Nil(t, insErr)

		insErr = db.InsertSONG("Stranger Aeons", albumID4, genreID2, authorID2, 5)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Chaos Breed", albumID4, genreID2, authorID2, 6)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Crawl", albumID4, genreID2, authorID2, 7)
		require.Nil(t, insErr)

		insErr = db.InsertSONG("Eyemaster", albumID5, genreID2, authorID2, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Rotten Soil", albumID5, genreID2, authorID2, 2)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Wolverine Blues", albumID5, genreID2, authorID2, 3)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Demon", albumID5, genreID2, authorID2, 4)
		require.Nil(t, insErr)

		insErr = db.InsertSONG("To Ride, Shoot Straight and Speak the Truth", albumID6, genreID2, authorID2, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Like This with the Devil", albumID6, genreID2, authorID2, 2)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Lights Out", albumID6, genreID2, authorID2, 3)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("DCLXVI", albumID6, genreID2, authorID2, 7)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Put Me Out", albumID6, genreID2, authorID2, 10)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Just as Sad", albumID6, genreID2, authorID2, 11)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Boats", albumID6, genreID2, authorID2, 12)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Uffe's Horrorshow", albumID6, genreID2, authorID2, 13)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Wreckage", albumID6, genreID2, authorID2, 14)
		require.Nil(t, insErr)

		insErr = db.InsertSONG("Addiction King", albumID7, genreID2, authorID2, 1)
		require.Nil(t, insErr)
	})
	t.Run("the same song 'Squares' in different albums of different authors", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID1, err := db.AddAuthor("KONVENT")
		assert.Nil(t, err)
		genreID1, err := db.AddGenre("Death/Doom Metal")
		assert.Nil(t, err)

		albumID1, albumErr1 := db.AddAlbum(authorID1, "Demo", 2017, "")
		assert.Nil(t, albumErr1)

		insErr := db.InsertSONG("Chernobyl Child", albumID1, genreID1, authorID1, 1)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Squares", albumID1, genreID1, authorID1, 2)
		require.Nil(t, insErr)

		// Second author
		authorID2, err := db.AddAuthor("Entombed")
		assert.Nil(t, err)
		genreID2, err := db.AddGenre("Death-n-Roll")
		assert.Nil(t, err)

		albumID2, albumErr3 := db.AddAlbum(authorID2, "Left Hand Path", 1990, "")
		assert.Nil(t, albumErr3)
		_, albumErr4 := db.AddAlbum(authorID2, "Clandestine", 1991, "")
		assert.Nil(t, albumErr4)

		insErr = db.InsertSONG("Squares", albumID2, genreID1, authorID1, 2)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Bridge", albumID2, genreID2, authorID1, 5)
		require.Nil(t, insErr)
	})
	// TODO: Write test where tracknum the same for one album
	// TODO: this test go without errors. Question: will be trackNum UNIQUE or not?
	t.Run("the same song in different authors", func(t *testing.T) {
		db, err := GetPostgresConnection()
		assert.Nil(t, err)

		dropErr := DropTablesAfterTest(db)
		assert.Nil(t, dropErr)

		err = CreateTablesForTest(db)
		assert.Nil(t, err)

		authorID, err := db.AddAuthor("Entombed")
		assert.Nil(t, err)
		genreID, err := db.AddGenre("Death-n-Roll")
		assert.Nil(t, err)

		albumID, albumErr3 := db.AddAlbum(authorID, "Left Hand Path", 1990, "")
		assert.Nil(t, albumErr3)

		insErr := db.InsertSONG("Eyemaster", albumID, genreID, authorID, 5)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Rotten Soil", albumID, genreID, authorID, 5)
		require.Nil(t, insErr)
		insErr = db.InsertSONG("Wolverine Blues", albumID, genreID, authorID, 3)
		require.Nil(t, insErr)
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

		albumID2, err3 := db.AddAlbum(authorID2, albumName, albumYear, cover)
		fmt.Printf("albumID2 = '%v'\n", albumID2)
		assert.Equal(t, 2, albumID2)
		assert.Nil(t, err3)

		expectLen := 2

		gotVal, gotErr := db.SelectALBUM()
		require.Nil(t, gotErr)
		require.Equal(t, expectLen, len(gotVal))
	})
}
