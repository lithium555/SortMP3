package draft

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FilePathWalkDir - The path/filepath package provides a handy way to scan all the files in a directory,
//it will automatically scan each sub-directories in the directory
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// IOReadDir - ioutil.ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename.
func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

// Using os.File.Readdir

// OSReadDir - Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values,
// as would be returned by Lstat, in directory order. Subsequent calls on the same file will yield further FileInfos.
func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
