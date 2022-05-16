package repositories

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"path/filepath"
	"time"
)

type GitRepository struct {
	Path string
}

// Commit commits all changes in the specified article
func (r GitRepository) Commit(article string) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// add all files
	_, err = w.Add(".")
	status, err := w.Status()

	// TODO: verify if this added all files
	fmt.Println(status)

	// commit
	_, err = w.Commit("version update", &git.CommitOptions{
		Author: &object.Signature{
			// TODO: add actual user name?
			Name:  "Alexandria Git Manager",
			Email: "",
			When:  time.Now(),
		},
	})
	return nil
}

func (r GitRepository) CheckoutBranch(article string, version string) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// checkout
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(version),
	})

	return nil
}

// GetArticlePath returns the path to an article git repository
func (r GitRepository) GetArticlePath(article string) string {
	return filepath.Join(r.Path, article)
}

// getWorktree returns the go-git worktree of an article git repository
func (r GitRepository) getWorktree(article string) (*git.Worktree, error) {
	// Open  repository.
	dir := r.GetArticlePath(article)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return w, nil
}
