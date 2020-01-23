package postgres

const (
	// CreateTableGENRE represents query for creating table GENRE.  +
	CreateTableGENRE = `CREATE TABLE IF NOT EXISTS genre (
												id SERIAL PRIMARY KEY, 
												genre_name TEXT UNIQUE
	);`
	// CreateTableAUTHOR represents query for creating table AUTHOR. +
	CreateTableAUTHOR = `CREATE TABLE IF NOT EXISTS author (
												id SERIAL PRIMARY KEY, 
												author_name TEXT UNIQUE
	);`
	// CreateTableALBUM represents query for creating table ALBUM.
	CreateTableALBUM = `CREATE TABLE IF NOT EXISTS album (
												id SERIAL PRIMARY KEY,
												author_id INT, 
												album_name TEXT,
												album_year INT,
												cover TEXT, 
												UNIQUE(author_id, album_name),
												FOREIGN KEY (author_id) REFERENCES author(id)
	);`
	// CreateTableSONG represents query for creating table SONG.  +
	CreateTableSONG = `CREATE TABLE IF NOT EXISTS song(  
												id SERIAL PRIMARY KEY,
												name_of_song TEXT,
												album_id INT,
												genre_id INT,
												author_id INT,
												track_number INT,
												UNIQUE(author_id, album_id, name_of_song),
												UNIQUE(album_id, track_number),
												FOREIGN KEY (album_id) REFERENCES album(id),
												FOREIGN KEY (genre_id) REFERENCES genre(id),
												FOREIGN KEY (author_id) REFERENCES author(id)
	);`
)

// http://www.postgresqltutorial.com/postgresql-identity-column/
//
// LINE 32:
// The order in the index may be important. In you case you can sort by the name of the song only.
// If you change to author_id, album_id, name, then you will be able to sort by the author, then album (of this author),
// then track name (in this album).
