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
	path         string
	articlePath  string
	requestPath  string
	downloadPath string
	defaultFile  string
}

// NewFileSystem creates file system struct and creates folders
func NewFileSystem(path string, defaultFile string) *FileSystem {

	// Create filesystem struct with paths
	fs := &FileSystem{
		path:         path,
		articlePath:  filepath.Join(path, "persistent", "articles"),
		requestPath:  filepath.Join(path, "persistent", "requests"),
		downloadPath: filepath.Join(path, "cache", "downloads"),
		defaultFile:  defaultFile}

	// Make all folder if they don't exist yet
	err := os.MkdirAll(fs.articlePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(fs.requestPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(fs.downloadPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return fs
}

// CreateArticlePath makes a new folder for an article git repository, returns the path.
// Fails if the folder already exists. Does not initialize the repository.
func (fs FileSystem) CreateArticlePath(article int64) (string, error) {
	// Get path to where the directory will be added
	path, err := fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Check if dir already exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return "", fmt.Errorf("trying to create an article repository folder that already exists with id=%d", article)
	}

	// Create directory
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, err
}

// GetArticlePath returns the path to an article git repository, does not check if it exists
func (fs FileSystem) GetArticlePath(article int64) (string, error) {
	idString := strconv.FormatInt(article, 10)
	path, err := filepath.Abs(filepath.Join(fs.articlePath, idString))
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}

// GetDownloadPath returns the path that a (temporary) zip for downloading can be stored in
func (fs FileSystem) GetDownloadPath(filename string) (string, error) {
	// check that the filename does not contain illegal characters
	if !strings.Contains(IllegalChars, filename) {
		return "", fmt.Errorf("filename for zip contains illegal characters: %s", filename)
	}

	// create the filepath, absolute and cleaned
	path := filepath.Join(fs.downloadPath, filename+".zip")
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}

// GetRequestPath returns the path to a folder for this request, with /old and /new folders
// creates the folder if it doesn't exist yet
func (fs FileSystem) GetRequestPath(article int64, request int64) (string, error) {

	// get the path by generating a unique cache id
	id := fmt.Sprintf("%d-%d", article, request)
	path, err := filepath.Abs(filepath.Join(fs.requestPath, id))
	if err != nil {
		return "", err
	}

	// create nested folders
	err = os.MkdirAll(filepath.Join(path, "old"), os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(path, "new"), os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Clean(path), nil
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

// PlaceDefaultFile places the default article contents, as specified from the config,
// in the specified path. Path should refer to a folder, the filename is added by the storer.
func (fs FileSystem) PlaceDefaultFile(path string) error {

	// Read the default file
	input, err := ioutil.ReadFile(fs.defaultFile)
	if err != nil {
		return err
	}

	// Write contents to the main.qmd file in the repo
	target := filepath.Join(path, "main.qmd")
	return ioutil.WriteFile(target, input, 0644)
}
