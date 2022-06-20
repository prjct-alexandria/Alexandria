package filesystem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// GetDownloadPath returns the path that a (temporary) dow
func (fs FileSystem) GetDownloadPath(filename string) (string, error) {
	// check that the filename does not contain illegal characters
	if !strings.Contains(IllegalChars, filename) {
		return "", fmt.Errorf("filename for zip contains illegal characters: %s", filename)
	}

	// create the filepath, absolute and cleaned
	path := filepath.Join(fs.path, "cache", "downloads", filename+".zip")
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}

// SaveArticleFile saves the file from the multipart file header as the main contents of an article
func (fs FileSystem) SaveArticleFile(c *gin.Context, file *multipart.FileHeader, articlePath string) error {
	filePath := filepath.Join(articlePath, "main.qmd")
	return c.SaveUploadedFile(file, filePath)
}

// ReadArticleFile returns the main contents of the file
func (fs FileSystem) ReadArticleFile(articlePath string) (string, error) {
	fileContent, err := ioutil.ReadFile(filepath.Join(articlePath, "main.qmd"))
	if err != nil {
		return "", err
	}
	return string(fileContent), nil
}