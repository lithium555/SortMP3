BEGIN;

CREATE SCHEMA IF NOT EXISTS sortMP3;

CREATE TABLE IF NOT EXISTS genre
(
                                     id SERIAL PRIMARY KEY,
                                     genre_name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS author
(
                                      id SERIAL PRIMARY KEY,
                                      author_name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS album
(
                                     id SERIAL PRIMARY KEY,
                                     author_id INT,
                                     album_name TEXT,
                                     album_year INT,
                                     cover TEXT,
                                     UNIQUE(author_id, album_name),
                                     FOREIGN KEY (author_id) REFERENCES author(id)
);

CREATE TABLE IF NOT EXISTS song
(
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
);

COMMIT;