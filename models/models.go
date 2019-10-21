package models

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
	Cover     string
}

// Song represents song fields in table `SONG`
type Song struct {
	SongID      uint64
	NameOfSong  string
	AlbumID     int
	GenreID     int
	AuthorID    int
	TrackNumber int
}
