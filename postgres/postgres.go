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

	// TableGenre represents table GENRE from postgres database.
	TableGenre = "GENRE"
	// TableAuthor represents table AUTHOR from postgres database.
	TableAuthor = "AUTHOR"
	// TableAlbum represents table ALBUM from postgres database.
	TableAlbum = "ALBUM"
	// TableSong represents table SONG from postgres database.
	TableSong = "SONG"
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
		INSERT INTO GENRE(name) 
						VALUES ('$1') 
						RETURNING id
	`).Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println("New record ID is:", id)
	return nil
}

// AddGenre represents the record insertion into table `GENRE`.
func (db *Database) AddGenre(genreName string) (int, error) {
	// Maybe this genre exists in our table `GENRE`, lets try to find it.
	existGenreID, err := db.GetExistsGenre(genreName)
	if err == nil {
		return existGenreID, nil
	}
	// If we can`t find Genre, let`s add it into table `GENRE`
	var genreID int
	err = db.GetConnection().QueryRow(`
		INSERT INTO GENRE(genre_name) 
			VALUES ($1) 	
			RETURNING id
		`, genreName).Scan(&genreID)
	if err != nil {
		log.Printf("Can`t insert genre `%v` into table. Error: '%v'\n", genreName, err)
		return 0, err
	}

	return genreID, nil
}

// AddAuthor represents the record insertion into table `AUTHOR`.
func (db *Database) AddAuthor(author string) (int, error) {
	// Maybe Author exist in our table, so let`s try to find his ID in a table
	existsAuthorID, err := db.GetExistsAuthor(author)
	if err == nil {
		return existsAuthorID, nil
	}
	// If Author not exist in our table - lets Insert him to table
	var authorID int
	err = db.PostgresConn.QueryRow(`
					INSERT INTO author (author_name)
					VALUES ($1) RETURNING id
		`, author).Scan(&authorID)
	if err != nil {
		log.Printf("Can`t Insert new Author '%v' in func AddAuthor(); Error: '%v'\n", author, err)
		return 0, err
	}

	return authorID, nil
}

// AddAlbum represents the record insertion into table `ALBUM`
func (db *Database) AddAlbum(authorID int, albumName string, albumYear int, cover string) (int, error) {
	// Let`s try to find out name of this album in table, maybe he is exist.

	// Sometimes name of albums are the same, but if we will seek them by 3 arguments,
	// like in this func GetExistsAlbum()
	existAlbumID, err := db.GetExistsAlbum(authorID, albumName, albumYear)
	if err == nil {
		return existAlbumID, nil
	}
	// if this album doesnt exist in table, lets Insert it into table:
	var coverToInsert interface{}
	if cover != "" {
		coverToInsert = cover
	}
	var albumID int
	err = db.PostgresConn.QueryRow(`
			INSERT INTO ALBUM(author_id, album_name, album_year, cover)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, authorID, albumName, albumYear, coverToInsert).Scan(&albumID)
	if err != nil {
		log.Printf("Can`t insert album '%v' into table in func AddAlbum(). Error: '%v'\n", albumName, err)
		return 0, err
	}

	return albumID, nil
}

// https://dba.stackexchange.com/questions/46410/how-do-i-insert-a-row-which-contains-a-foreign-key

// InsertSONG represents the record insertion into table `SONG`
func (db *Database) InsertSONG(songName string, albumID int, genreID int, authorID int, trackNum int) error {
	_, err := db.PostgresConn.Exec(`
	INSERT INTO SONG(
		name_of_song,
		album_id,
		genre_id,
		author_id,
		track_number)
	VALUES ($1, $2, $3, $4, $5)
	`, songName, albumID, genreID, authorID, trackNum)
	if err != nil {
		log.Printf("InsertSONG(), Error: '%v'\n", err)
		return err
	}
	return nil
}

// Drop represents deleting table from postgres database. tableName - name of table for deleting.
func (db *Database) Drop(tableName string) error {
	//TODO: fix deoping tables.
	// http://jinzhu.me/gorm/database.html#migration
	_, err := db.PostgresConn.Exec("DROP TABLE if exists " + tableName + " ;")

	if err != nil {
		log.Printf("DROP, Error: '%v'\n", err)
		return err
	}
	return nil
}

// GetExistsAuthor will find AuthorID, if this author exists in table `AUTHOR`
func (db *Database) GetExistsAuthor(author string) (int, error) {
	var existAuthor models.Author

	row := db.PostgresConn.QueryRow(`SELECT id FROM author WHERE author_name = $1;`, author)
	if err := row.Scan(&existAuthor.AuthorID); err == sql.ErrNoRows {
		log.Errorf("Func GetExistsAuthor(). Zero rows found")
		return 0, err
	} else if err != nil {
		log.Errorf("Func GetExistsAuthor(). Error in row.Scan(). Error: '%v'\n", err)
		return 0, err
	}

	return existAuthor.AuthorID, nil
}

// GetExistsGenre will find genrteID if this genre exits in table `GENRE`
func (db *Database) GetExistsGenre(genreName string) (int, error) {
	row := db.PostgresConn.QueryRow(`SELECT id FROM GENRE WHERE genre_name = $1;`, genreName)

	var genreExist models.Genre
	if err := row.Scan(&genreExist.GenreID); err == sql.ErrNoRows {
		log.Errorf("Func GetExistsGenre(): Zero rows found.")
		return 0, err
	} else if err != nil {
		log.Errorf("Func  GetExistsGenre(). Error in row.Scan(). Error: '%v'\n", err)
		return 0, err
	}

	return genreExist.GenreID, nil
}

// GetExistsAlbum will find albumID if this album exists in table `ALBUM`
func (db *Database) GetExistsAlbum(authorID int, albumName string, albumYear int) (int, error) {
	row := db.PostgresConn.QueryRow(`
		SELECT id FROM ALBUM 
		WHERE author_id = $1 AND album_name = $2 AND album_year = $3`, authorID, albumName, albumYear)

	var albumExist models.Album

	if err := row.Scan(&albumExist.AlbumID); err == sql.ErrNoRows {
		log.Errorf("Func GetExistsAlbum(). Zero rows found.")
		return 0, err
	} else if err != nil {
		log.Errorf("Func  GetExistsAlbum(). Error in row.Scan(). Error: '%v'\n", err)
		return 0, err
	}

	return albumExist.AlbumID, nil
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
