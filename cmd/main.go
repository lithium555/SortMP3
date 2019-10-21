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

func run() error {
	postgres, err := postgres.GetPostgresConnection()
	if err != nil {
		log.Printf("postgre.GetConnection() = '%v'\n", err)
		return err
	}
	defer postgres.Close()

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
		if strings.Contains(filepath, ".jpg") {
			fmt.Printf("This is a picture: '%v'\n", filepath)
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

		var authorID int
		id, err := postgres.InsertAUTHOR(md.Artist())
		if err != nil {
			//log.Printf("InsertAUTHOR(), Errror: '%v'\n", err)
			existsAuthorID, err := postgres.GetExistsAuthor(md.Artist())
			if err != nil {
				log.Println("Error in func GetExistsAuthor()")
				return err
			}
			authorID = existsAuthorID.AuthorID
		} else {
			authorID = id
		}
		fmt.Printf(">>>>>>>>>>>>>>>>authorID = '%v'\n", authorID)

		var genreID int
		genID, err := postgres.InsertGENRE(md.Genre())
		if err != nil {
			existGenre, err := postgres.GetExistsGenre(md.Genre())
			if err != nil {
				log.Println("Error in func GetExistsGenre()")
				return err
			}
			genreID = existGenre.GenreID
		} else {
			genreID = genID
		}
		fmt.Printf("+++++++++++++++++++++genreID = '%v'\n", genreID)

		numberOfTrack, _ := md.Track()

		var albumID int
		albID, err := postgres.InsertALBUM(authorID, md.Album(), md.Year(), "")
		if err != nil {
			// Sometimes name of albums are the same, but if we will sekk them by 3 arguments,
			// like in this func GetExistsAlbum()
			album, err := postgres.GetExistsAlbum(authorID, md.Album(), md.Year())
			if err != nil {
				log.Println("Error in func GetExistsAlbum()")
				return err
			}
			albumID = album.AlbumID
		} else {
			albumID = albID
		}
		fmt.Printf("========================albumID = '%v'\n", albumID)

		err = postgres.InsertSONG(md.Title(), albumID, genreID, authorID, numberOfTrack)
		if err != nil {
			log.Printf("InsertSONG(), Error: '%v'\n", err)
			return err
		}
		fmt.Println()
	}

	if err := sort.DropAllTables(postgres); err != nil {
		return err
	}

	return nil
}
