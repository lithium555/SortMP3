package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lithium555/SortMP3/models"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sorter"
	password = "master"
	dbname   = "musicDB"

	// DropGenre represents deleting table GENRE from postgres database.
	DropGenre = "DROP TABLE GENRE"
	// DropAuthor represents deleting table AUTHOR from postgres database.
	DropAuthor = "DROP TABLE AUTHOR"
	// DropAlbum represents deleting table ALBUM from postgres database.
	DropAlbum = "DROP TABLE ALBUM"
	// DropSong represents deleting table SONG from postgres database.
	DropSong = "DROP TABLE SONG"
)

// GetPostgresConnection represents connection to PostgresDB
func GetPostgresConnection() (Database, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Errorf("Postgres connection error: '%v'\n", err)
		return Database{}, err
	}

	err = db.Ping()
	if err != nil {
		log.Errorf("Postgres Ping() failed. Error: '%v'\n", err)
		return Database{}, err
	}
	log.Println("Postgres successfully connected!")

	getPostgres := Database{PostgresConn: db}

	return getPostgres, nil
}

// Database represents implementation, how to access to the Postgres database
type Database struct {
	PostgresConn *sql.DB
}

// Close closes connection with the database
func (db *Database) Close() error {
	if err := db.PostgresConn.Close(); err != nil {
		return err
	}
	return nil
}

// GetConnection returns database connection
func (db *Database) GetConnection() *sql.DB {
	return db.PostgresConn
}

// Ping runs a trivial ping command just to get in touch with the server
func (db *Database) Ping() error {
	return db.PostgresConn.Ping()
}

// CreateTable represents creating table in postgresDB by using `query`
func (db *Database) CreateTable(query string) error {
	_, err := db.GetConnection().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// InsertIntoTableGENRETest represents inserting value into table GENRE
func (db *Database) InsertIntoTableGENRETest() error {
	id := 0
	err := db.GetConnection().QueryRow(`
		INSERT INTO GENRE(id, name) 
						VALUES (DEFAULT, '$1') 
						RETURNING id
	`).Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println("New record ID is:", id)
	return nil
}

// InsertGENRE represents the record insertion into table `GENRE`.
func (db *Database) InsertGENRE(name string) (int, error) {
	var genreID int
	err := db.GetConnection().QueryRow(`
	INSERT INTO GENRE(id, genre_name) 
		VALUES (DEFAULT, $1) 	
		RETURNING id
	`, name).Scan(&genreID)
	if err != nil {
		return 0, err
	}
	return genreID, nil
}

// InsertAUTHOR represents the record insertion into table `AUTHOR`
func (db *Database) InsertAUTHOR(author string) (int, error) {
	var authorID int
	err := db.PostgresConn.QueryRow(`
		INSERT INTO AUTHOR (author_name)
		VALUES ($1) RETURNING id
	`, author).Scan(&authorID)
	if err != nil {
		return 0, err
	}
	return authorID, nil
}

// InsertALBUM represents the record insertion into table `ALBUM`
func (db *Database) InsertALBUM(authorID int, albumName string, albumYear int, cover string) (int, error) {
	var albumID int
	err := db.PostgresConn.QueryRow(`
		INSERT INTO ALBUM(author_id, album_name, album_year, cover)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, authorID, albumName, albumYear, cover).Scan(&albumID)
	if err != nil {
		return 0, err
	}
	return albumID, nil
}

// InsertSONG represents the record insertion into table `SONG`
func (db *Database) InsertSONG(songName string, albumID int, genreID int, authorID int, trackNum int) error {
	_, err := db.PostgresConn.Exec(`
		INSERT INTO SONG(id, name_of_song, album_id, genre_id, author_id, track_number)
		VALUES (DEFAULT, $1, $2, $3, $4, $5)
	`, songName, albumID, genreID, authorID, trackNum)
	if err != nil {
		log.Printf("InsertSONG(), Error: '%v'\n", err)
		return err
	}
	return nil
}

// Drop represents deleting table from postgres database. sqlQuery - query for deleting.
func (db *Database) Drop(sqlQuery string) error {
	_, err := db.PostgresConn.Exec(sqlQuery)
	if err != nil {
		log.Printf("DROP, Error: '%v'\n", err)
		return err
	}
	return nil
}

// GetExistsAuthor will will find Author, which exists.
func (db *Database) GetExistsAuthor(author string) (models.Author, error) {
	// TODO: fix query: SELECT id FROM AUTHOR WHERE author_name = ?;
	rows, err := db.PostgresConn.Query(`SELECT id FROM AUTHOR WHERE author_name = ?;`, author)
	if err != nil {
		log.Println("Func GetExistsAuthor()")
		return models.Author{}, err
	}
	defer rows.Close()

	var existAuthor models.Author
	for rows.Next() {
		err := rows.Scan(existAuthor.AuthorID, existAuthor.AuthorName)
		if err != nil {
			log.Errorf("Func GetExistsAuthor(). Error in rows.Scan(). Error: '%v'\n", err)
			return models.Author{}, err
		}
	}

	return existAuthor, nil
}

// SelectSONG represents sql SELECT query for table SONG.
func (db *Database) SelectSONG() ([]*models.Song, error) {
	rows, err := db.PostgresConn.Query(`
			SELECT id, name_of_song, album_id, genre_id, author_id, track_number 
			FROM SONG`)
	if err != nil {
		log.Printf("Select() for table `SONG` not passed. Error: '%v'\n", err)
		return nil, err
	}
	defer rows.Close()

	var songs []*models.Song
	for rows.Next() {
		song := new(models.Song)
		err := rows.Scan(&song.SongID, &song.NameOfSong, &song.AlbumID, &song.GenreID, &song.AuthorID, &song.TrackNumber)
		if err != nil {
			log.Errorf("Func SelectSONG(). Error in rows.Scan(). Error: '%v'\n", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("rows.Err(), Error: '%v'\n", err)
		return nil, err
	}
	for _, song := range songs {
		fmt.Printf("song.SongID = '%v'\n", song.SongID)
		fmt.Printf("song.NameOfSong = '%v'\n", song.NameOfSong)
		fmt.Printf("song.AlbumID = '%v'\n", song.AlbumID)
		fmt.Printf("song.GenreID = '%v'\n", song.GenreID)
		fmt.Printf("song.AuthorID = '%v'\n", song.AuthorID)
		fmt.Printf("song.TrackNumber = '%v'\n", song.TrackNumber)
		fmt.Println()
		fmt.Println()
	}
	return songs, nil
}

// SelectGENRE represents sql SELECT query for table GENRE.
func (db *Database) SelectGENRE() ([]*models.Genre, error) {
	rows, err := db.PostgresConn.Query(`
		SELECT id, genre_name 
		FROM GENRE`)
	if err != nil {
		log.Printf("Select() for table `GENRE` not passed. Error: '%v'\n", err)
		return nil, err
	}
	defer rows.Close()

	genre := make([]*models.Genre, 0)
	for rows.Next() {
		g := new(models.Genre)
		if err := rows.Scan(&g.GenreID, &g.GenreName); err != nil {
			log.Errorf("Func SelectGenre(). Error in rows.Scan(). Error: '%v'\n", err)
			return nil, err
		}
		genre = append(genre, g)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("rows.Err(), Error: '%v'\n", err)
		return nil, err
	}
	for _, g := range genre {
		fmt.Printf("g.GenreID = '%v'\n", g.GenreID)
		fmt.Printf("g.GenreName = '%v'\n'", g.GenreName)
		fmt.Println()
	}

	return genre, nil
}

// SelectAUTHOR represents sql-query SELECT for table AUTHOR
func (db *Database) SelectAUTHOR() ([]*models.Author, error) {
	rows, err := db.PostgresConn.Query(`
			SELECT id, author_name 
			FROM AUTHOR`)
	if err != nil {
		log.Printf("Select() for table `AUTHOR` not passed. Error: '%v'\n", err)
		return nil, err
	}
	defer rows.Close()

	authors := make([]*models.Author, 0)
	for rows.Next() {
		author := new(models.Author)
		if err := rows.Scan(&author.AuthorID, &author.AuthorName); err != nil {
			log.Errorf("Func SelectAuthor(). Error in rows.Scan(). Error: '%v'\n", err)
			return nil, err
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("rows.Err(), Error: '%v'\n", err)
		return nil, err
	}

	for _, a := range authors {
		fmt.Printf("a.AuthorID = '%v'\n", a.AuthorID)
		fmt.Printf("a.AuthorName = '%v'\n", a.AuthorName)
		fmt.Println()
	}

	return authors, nil
}

// SelectALBUM represents sql-query SELECT for table `ALBUM`
func (db *Database) SelectALBUM() ([]*models.Album, error) {
	rows, err := db.PostgresConn.Query(`
			SELECT id, author_id, album_name, album_year, cover  
			FROM ALBUM`)
	if err != nil {
		log.Printf("Select() for table `ALBUM` not passed. Error: '%v'\n", err)
		return nil, err
	}
	defer rows.Close()

	albums := make([]*models.Album, 0)
	for rows.Next() {
		album := new(models.Album)
		if err := rows.Scan(&album.AlbumID, &album.AuthorID, &album.AlbumName, &album.AlbumYear, &album.Cover); err != nil {
			log.Errorf("Func SelectAlbum(). Error in rows.Scan(). Error: '%v'\n", err)
			return nil, err
		}
		albums = append(albums, album)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("rows.Err(), Error: '%v'\n", err)
		return nil, err
	}
	for _, al := range albums {
		fmt.Printf("al.AlbumID = '%v'\n", al.AlbumID)
		fmt.Printf("al.AuthorID = '%v'\n", al.AuthorID)
		fmt.Printf("al.AlbumName = '%v'\n", al.AlbumName)
		fmt.Printf("al.AlbumYear = '%v'\n", al.AlbumYear)
		fmt.Printf("al.Cover = '%v'\n", al.Cover)
		fmt.Println()
	}
	return albums, nil
}
