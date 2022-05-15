package repositories

import (
	"github.com/go-git/go-git/v5"
	"path/filepath"
	"fmt"
)

type GitRepository struct {
	path string
}

// LoadVersionPath checks out the required version and returns the filepath to the repository
// Can be used for reading and writing a version's files.
// Version contents must be read/write ASAP after this function,
// some other calls might cause the required  version to not be checked out anymore.
func (r GitRepository) LoadVersionPath(article string, version string) (string, error) {
	r.checkoutBranch(article, version)

	return "",nil
}

// Commit commits all changes in the specified branch
func (r GitRepository) Commit(article string, version string) error {
	// Open  repository.
	dir := r.getArticlePath(article)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// TODO: verify if this adds all
	_, err = w.Add(".")
	status, err := w.Status()
	fmt.Println(status)

	return nil
}

func (r GitRepository) checkoutBranch(article string, version string) error {
	// Open  repository.
	dir := r.getArticlePath(article)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// TODO: verify if this adds all
	_, err = w.Add(".")
	status, err :=
	return nil
}

// getVersionPath returns the path to an article git repository
func (r GitRepository) getArticlePath(article string) string {
	return filepath.Join(r.path, article)
}
