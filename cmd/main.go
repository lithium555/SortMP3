package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dhowden/tag"
	_ "github.com/lib/pq"
	"github.com/lithium555/SortMP3/postgres"
	"github.com/lithium555/SortMP3/sort"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run(). Error: '%v'\n", err)
	}
}

func run() error {
	db, err := postgres.GetPostgresConnection()
	if err != nil {
		log.Printf("postgre.GetConnection() = '%v'\n", err)
		return err
	}

	getPostgres := postgres.Database{PostgresConn: db}
	fmt.Printf("getPostgres = '%v'\n", getPostgres)

	if err := sort.CreateTables(getPostgres); err != nil {
		return err
	}

	files, err := sort.Variant1()
	if err != nil {
		fmt.Printf(" variant2(), Error: '%v'\n", err)
		return err
	}

	for _, file := range files {
		fmt.Println(file)
		mp3, err := os.Open(file)
		if err != nil {
			log.Printf("os.open(), Error: '%v'\n", err)
			return err
		}
		md, err := tag.ReadFrom(mp3)
		if err != nil {
			// format of file is not mp3
			log.Printf("tag.ReadFrom(mp3), error: '%v'\n", err)
			continue
		}
		authorID, err := getPostgres.InsertAUTHOR(md.Artist())
		if err != nil {
			log.Printf("InsertAUTHOR(), Errror: '%v'\n", err)
			return err
		}
		fmt.Printf(">>>>>>>>>>>>>>>>authorID = '%v'\n", authorID)

		genreID, err := getPostgres.InsertGENRE(md.Genre())
		if err != nil {
			log.Printf("InsertIntoTableGENRE(). Error: '%v'\n", err)
			return err
		}
		fmt.Printf("+++++++++++++++++++++genreID = '%v'\n", genreID)

		numberOfTrack, _ := md.Track()

		albumID, err := getPostgres.InsertALBUM(authorID, md.Album(), md.Year(), "")
		if err != nil {
			log.Printf("InsertALBUM(), Error: '%v'\n", err)
			return err
		}
		fmt.Printf("========================albumID = '%v'\n", albumID)

		err = getPostgres.InsertSONG(md.Title(), albumID, genreID, authorID, numberOfTrack)
		if err != nil {
			log.Printf("InsertSONG(), Error: '%v'\n", err)
			return err
		}
		fmt.Println()
	}

	//if err := sort.DropAllTables(getPostgres); err != nil{
	//	return err
	//}

	return nil
}
