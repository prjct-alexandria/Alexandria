package filesystem

import (
	"os"
	"path/filepath"
	"strconv"
)

type FileSystem struct {
	path string
}

// NewFileSystem creates file system struct and creates folders
func NewFileSystem(path string) *FileSystem {

	// make folder for git repos
	err := os.MkdirAll(filepath.Join(path, "persistent", "articles"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// make folder for request before-and-afters
	err = os.MkdirAll(filepath.Join(path, "persistent", "requests"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// make (cache) folder for download zip's
	err = os.MkdirAll(filepath.Join(path, "cache", "downloads"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &FileSystem{path: path}
}

// GetArticlePath returns the path to an article git repository, does not check if it exists
func (fs FileSystem) GetArticlePath(article int64) (string, error) {
	idString := strconv.FormatInt(article, 10)
	path, err := filepath.Abs(filepath.Join(fs.path, "persistent", idString))
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}
