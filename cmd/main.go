package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/lithium555/SortMP3/draft"
	"github.com/lithium555/SortMP3/postgres"
)

func main() {

	// VARIANT 1:
	variant1()

	//
	//// VARIANT 2:
	//variant2()

	//if err := run(); err != nil {
	//	log.Fatalf("run(). Error: '%v'\n", err)
	//}
}

func variant1() {
	//https://flaviocopes.com/go-list-files/

	var files []string

	root := "./readData"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
		fmt.Println()
	}
}

func variant2() {
	var (
		root  string
		files []string
		err   error
	)

	root = "./readData/"

	// filepath.Walk
	files, err = draft.FilePathWalkDir(root)
	if err != nil {
		panic(err)
	}
	// ioutil.ReadDir
	files, err = draft.IOReadDir(root)
	if err != nil {
		panic(err)
	}
	//os.File.Readdir
	files, err = draft.OSReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}
func run() error {
	db, err := postgres.GetPostgresConnection()
	if err != nil {
		log.Printf("postgre.GetConnection() = '%v'\n", err)
	}

	getPostgres := postgres.Database{PostgresConn: db}
	fmt.Printf("getPostgres = '%v'\n", getPostgres)

	//if err := sort.DropAllTables(getPostgres); err != nil{
	//	return err
	//}

	//if err := sort.CreateTables(getPostgres); err != nil{
	//	return err
	//}

	if err := getPostgres.InsertGENRE("melodic death"); err != nil {
		log.Printf("InsertIntoTableGENRE(). Error: '%v'\n", err)
		return err
	}

	if err := getPostgres.InsertAUTHOR("Dark Tranquillity"); err != nil {
		log.Printf("InsertAUTHOR(), Errror: '%v'\n", err)
		return err
	}
	err = getPostgres.InsertALBUM(1, "The Gallery", 1995, "https://www.google.com/search?q=the+gallery+album&rlz=1C5CHFA_enUA852UA852&sxsrf=ACYBGNQDxJziTc5-WlMhYm4BhluCOmQOkQ:1569847250069&source=lnms&tbm=isch&sa=X&ved=0ahUKEwjxosTdyPjkAhXOlIsKHehlAdoQ_AUIEigB&biw=2560&bih=1248#imgdii=Yb50dLKiAQmnYM:&imgrc=FoCQKyna2fB3gM:")
	if err != nil {
		log.Printf("InsertALBUM(), Error: '%v'\n", err)
		return err
	}

	err = getPostgres.InsertSONG("Punish My Heaven", 45, 1, 5, 1)
	if err != nil {
		log.Printf("InsertSONG(), Error: '%v'\n", err)
		return err
	}
	//_, err = getPostgres.SelectSONG()
	//if err != nil{
	//	return err
	//}

	//_, err = getPostgres.SelectGenre()
	//if err != nil{
	//	return err
	//}

	_, err = getPostgres.SelectAlbum()
	if err != nil {
		return err
	}
	return nil
}
