package sort

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/lithium555/SortMP3/postgres"
)

// CreateTables will create 4 tables for the project: GENRE, AUTHOR, ALBUM, SONG
func CreateTables(getPostgres postgres.Database) error {
	for _, query := range []string{
		postgres.CreateTableGENRE,
		postgres.CreateTableAUTHOR,
		postgres.CreateTableALBUM,
		postgres.CreateTableSONG,
	} {
		if err := getPostgres.CreateTable(query); err != nil {
			log.Printf("Can`t  '%v', error: '%v'\n", query, err)
			return err
		}
	}

	return nil
}

// Variant1 represents a walker of all mp3 files.
func Variant1(root string) ([]string, error) {
	//https://flaviocopes.com/go-list-files/

	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // Если это не директория а файл, то добавляем путь в слайс
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func selectAll(getPostgres postgres.Database) error {
	_, err := getPostgres.SelectSONG()
	if err != nil {
		return err
	}

	_, err = getPostgres.SelectGENRE()
	if err != nil {
		return err
	}

	_, err = getPostgres.SelectALBUM()
	if err != nil {
		return err
	}

	return nil
}

// PrintID3components represents reading all ID3 components ofr mp3 file.
// More information here: https://en.wikipedia.org/wiki/ID3
func PrintID3components(file string) error {
	log.Println(file)
	mp3, err := os.Open(file)
	if err != nil {
		log.Printf("os.open(), Error: '%v'\n", err)
	}
	md, err := tag.ReadFrom(mp3)
	if err != nil {
		log.Printf("tag.ReadFrom(mp3), error: '%v'\n", err)
		return err
	}
	log.Printf("md.Title() = '%v'\n", md.Title())
	log.Printf("md.Album() = '%v'\n", md.Album())
	log.Printf("md.Artist() = '%v'\n", md.Artist())
	log.Printf("md.AlbumArtist() = '%v'\n", md.AlbumArtist())
	log.Printf("md.Composer() = '%v'\n", md.Composer())
	log.Printf("md.Genre() = '%v'\n", md.Genre())
	log.Printf("md.Year() = '%v'\n", md.Year())
	numberOfTrack, totalTracks := md.Track()
	log.Printf("numberOfTrack = '%v', totalTracks = '%v'\n", numberOfTrack, totalTracks)
	d1, d2 := md.Disc()
	log.Printf("d1 = '%v', d2 = '%v'\n", d1, d2)
	log.Printf("md.Lyrics() = '%v'\n", md.Lyrics())
	log.Printf("md.Comment() ='%v'\n", md.Comment())
	log.Printf("md.Raw() = '%v'\n", md.Raw())
	log.Println()

	return nil
}
