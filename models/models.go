package models

import "database/sql"

// Genre represents genre fields in table `GENRE`
type Genre struct {
	GenreID   int
	GenreName string
}

// Author represents author fields in table `AUTHOR`
type Author struct {
	AuthorID   int
	AuthorName string
}

// Album represents album fields in table `ALBUM`
type Album struct {
	AlbumID   int
	AuthorID  int
	AlbumName string
	AlbumYear int
	Cover     sql.NullString // https://github.com/golang/go/wiki/SQLInterface
}

/*
	Why  field Cover has type  'sql.NullString' ?

	var an_int64 sql.NullInt64
	var a_string sql.NullString
	var another_string sql.NullString
	row := db.QueryRow("SELECT id, name, thumbUrl FROM x WHERE y=? LIMIT 1", model.ID)
	err := row.Scan(&an_int64, &a_string, &another_string))
*/

// Song represents song fields in table `SONG`
type Song struct {
	SongID      int
	NameOfSong  string
	AlbumID     int
	GenreID     int
	AuthorID    int
	TrackNumber int
}
