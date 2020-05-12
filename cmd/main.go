package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

// TODO: Insert 1000 songs per one insert, but not insert one by one. You can open a transaction and do all inserts/selects on it - would be much faster.

func run() error {
	postgres, err := postgres.GetPostgresConnection()
	if err != nil {
		log.Printf("postgre.GetConnection() = '%v'\n", err)
		return err
	}
	defer postgres.Close()

	err = postgres.DropTables(postgres)
	if err != nil {
		return err
	}

	fmt.Printf("postgres = '%v'\n", postgres)

	if err := sort.CreateTables(postgres); err != nil {
		return err
	}

	root := "./readData"
	files, err := sort.Variant1(root)
	if err != nil {
		fmt.Printf(" variant2(), Error: '%v'\n", err)
		return err
	}

	for _, filepath := range files {
		fmt.Printf("Road to file, which we Open: '%v'\n", filepath)
		if !strings.Contains(filepath, ".mp3") {
			fmt.Printf("This is not MP3 file: '%v'\n", filepath)
			continue
		}

		mp3, err := os.Open(filepath)
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

		authorID, err := postgres.AddAuthor(md.Artist())
		if err != nil {
			log.Printf("AddAuthor(), Error: '%v'\n", err)
			return err
		}
		fmt.Printf(">>>>>>>>>>>>>>>>authorID = '%v'\n", authorID)

		genreID, err := postgres.AddGenre(md.Genre())
		if err != nil {
			log.Printf("AddGenre(), Error: '%v'\n", err)
			return err
		}
		fmt.Printf("+++++++++++++++++++++genreID = '%v'\n", genreID)

		numberOfTrack, _ := md.Track()

		albumID, err := postgres.AddAlbum(authorID, md.Album(), md.Year(), "")
		if err != nil {
			return err
		}
		fmt.Printf("========================albumID = '%v'\n", albumID)
		fmt.Printf("==================numberOfTrack = '%v'\n", numberOfTrack)

		err = postgres.InsertSONG(md.Title(), albumID, genreID, authorID, numberOfTrack)
		if err != nil {
			log.Printf("InsertSONG(), Error: '%v'\n", err)
			return err
		}
		fmt.Println()
	}

	if err := postgres.DropAllTables(postgres); err != nil {
		return err
	}

	return nil
}
