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
												author_id INT  UNIQUE, 
												album_name TEXT UNIQUE,
												album_year INT,
												cover TEXT, 
												FOREIGN KEY (author_id) REFERENCES author(id)
	);`
	// CreateTableSONG represents query for creating table SONG.  +
	CreateTableSONG = `CREATE TABLE IF NOT EXISTS song(  
												id SERIAL PRIMARY KEY,
												name_of_song TEXT UNIQUE,
												album_id INT,
												genre_id INT,
												author_id INT  UNIQUE,
												track_number INT, 
												FOREIGN KEY (album_id) REFERENCES album(id),
												FOREIGN KEY (genre_id) REFERENCES genre(id),
												FOREIGN KEY (author_id) REFERENCES author(id)
	);`
)

// http://www.postgresqltutorial.com/postgresql-identity-column/
