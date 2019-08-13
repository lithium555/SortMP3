package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// https://flaviocopes.com/go-list-files/

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

	//var(
	//	root string
	//	files []string
	//	err error
	//)
	//
	//root = "./readData/Drawing The Endless Shore - Protagonist (2013)"
	//
	//// filepath.Walk
	//files, err = FilePathWalkDir(root)
	//if err != nil {
	//	panic(err)
	//}
	//// ioutil.ReadDir
	//files, err = IOReadDir(root)
	//if err != nil {
	//	panic(err)
	//}
	////os.File.Readdir
	//files, err = OSReadDir(root)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, file := range files {
	//	fmt.Println(file)
	//}
}
