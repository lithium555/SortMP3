package postgres

const (
	// CreateTableGENRE represents query for creating table GENRE.  +
	CreateTableGENRE = `CREATE TABLE IF NOT EXISTS GENRE (
												"id" SERIAL PRIMARY KEY, 
												"genre_name" TEXT
	);`
	// CreateTableAUTHOR represents query for creating table AUTHOR. +
	CreateTableAUTHOR = `CREATE TABLE IF NOT EXISTS AUTHOR (
												"id" SERIAL PRIMARY KEY, 
												"author_name" TEXT
	);`
	// CreateTableALBUM represents query for creating table ALBUM.
	CreateTableALBUM = `CREATE TABLE IF NOT EXISTS ALBUM (
												"id" SERIAL PRIMARY KEY,
												"author_id" INT, 
												"album_name" TEXT,
												"album_year" INT,
												"cover" TEXT
	);`
	// CreateTableSONG represents query for creating table SONG.  +
	CreateTableSONG = `CREATE TABLE IF NOT EXISTS SONG(  
												"id" SERIAL PRIMARY KEY,
												"name_of_song" TEXT,
												"album_id" INT,
												"genre_id" INT,
												"author_id" INT,
												"track_number" INT
	);`
)

// http://www.postgresqltutorial.com/postgresql-identity-column/
