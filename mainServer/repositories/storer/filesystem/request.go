package filesystem

import (
	"io/ioutil"
	"path/filepath"
)

// CopyArticleToRequest copies the article in the specified article folder, to a sub folder of the request path,
// which sub folder depends on the afterMerging parameter
func (fs *FileSystem) CopyArticleToRequest(articlePath string, requestPath string, afterMerging bool) error {

	// determine target location
	var targetPath string
	if afterMerging {
		targetPath = filepath.Join(requestPath, "new", "main.qmd")
	} else {
		targetPath = filepath.Join(requestPath, "old", "main.qmd")
	}

	// copy article file to target
	input, err := ioutil.ReadFile(filepath.Join(articlePath, "main.qmd"))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(targetPath, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetRequestComparison returns the before and after main article file of a request
// returns (before, after, err)
func (fs FileSystem) GetRequestComparison(requestPath string) (string, string, error) {

	// read both old and new file from the cache
	oldFile, err := ioutil.ReadFile(filepath.Join(requestPath, "old", "main.qmd"))
	if err != nil {
		return "", "", err
	}
	newFile, err := ioutil.ReadFile(filepath.Join(requestPath, "new", "main.qmd"))
	if err != nil {
		return "", "", err
	}

	return string(oldFile), string(newFile), nil
}
