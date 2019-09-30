package sort

import (
	"log"

	"github.com/lithium555/SortMP3/postgres"
)

// CreateTables will create 4 tables for the project: GENRE, AUTHOR, ALBUM, SONG
func CreateTables(getPostgres postgres.Database) error {
	if err := getPostgres.CreateTable(postgres.CreateTableGENRE); err != nil {
		log.Printf("Can`t  create table `GENRE`, error: '%v'\n", err)
		return err
	}

	if err := getPostgres.CreateTable(postgres.CreateTableAUTHOR); err != nil {
		log.Printf("Can`t  create table `AUTHOR`, error: '%v'\n", err)
		return err
	}

	if err := getPostgres.CreateTable(postgres.CreateTableALBUM); err != nil {
		log.Printf("Can`t  create table `ALBUM`, error: '%v'\n", err)
		return err
	}

	if err := getPostgres.CreateTable(postgres.CreateTableSONG); err != nil {
		log.Printf("Can`t  create table `SONG`, error: '%v'\n", err)
		return err
	}
	return nil
}

// DropAllTables will drop all tables: GENRE, AUTHOR, ALBUM, SONG
func DropAllTables(getPostgres postgres.Database) error {
	if err := getPostgres.Drop(postgres.DropGenre); err != nil {
		return err
	}
	if err := getPostgres.Drop(postgres.DropAuthor); err != nil {
		return err
	}
	if err := getPostgres.Drop(postgres.DropAlbum); err != nil {
		return err
	}
	if err := getPostgres.Drop(postgres.DropSong); err != nil {
		return err
	}
	return nil
}
