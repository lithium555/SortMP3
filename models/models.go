package models

// Genre represents genre fields in table `GENRE`
type Genre struct {
	GenreID   string
	GenreName string
}

// Author represents author fields in table `AUTHOR`
type Author struct {
	AuthorID   string
	AuthorName string
}

// Album represents album fields in table `ALBUM`
type Album struct {
	AlbumID   string
	AuthorID  int
	AlbumName string
	AlbumYear int
	Cover     string
}

// Song represents song fields in table `SONG`
type Song struct {
	SongID      string
	NameOfSong  string
	AlbumID     int
	GenreID     int
	AuthorID    int
	TrackNumber int
}
