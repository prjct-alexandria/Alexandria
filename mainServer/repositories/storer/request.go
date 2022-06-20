package storer

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetRequestComparisonPath returns the path to a with /old and /new folders,
// for viewing what a request did or will do when accepted
func (s Storer) GetRequestComparisonPath(article int64, request int64) (string, error) {

	// get the path by generating a unique cache id
	id := fmt.Sprintf("%d-%d", article, request)
	path, err := filepath.Abs(filepath.Join(s.Path, "requests", id))
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
