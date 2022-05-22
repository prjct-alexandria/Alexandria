package repositories

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"mainServer/utils/clock"
	"os"
	"path/filepath"
	"strconv"
)

type GitRepository struct {
	Path  string
	Clock clock.Clock
}

// NewGitRepository creates a new GitRepository class.
// This is NOT the function used to create a folder/git repository to store an article in.
// See CreateRepo instead
func NewGitRepository(path string) (GitRepository, error) {

	// make folder for git files
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return GitRepository{}, err
	}

	return GitRepository{Path: path, Clock: clock.RealClock{}}, nil
}

// CreateRepo creates a new folder/git repository to store an article in.
// This is NOT the function used to create the Go repository class,
// see NewGitRepository instead,
// Return ID of created repository
func (r GitRepository) CreateRepo(id int64) error {
	path, err := r.GetArticlePath(id)
	if err != nil {
		return err
	}

	// Check if folder with the same name already exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("trying to create a git repository that already exists with id=%d", id)
	}

	// Create directory
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = git.PlainInit(path, false)
	if err != nil {
		return err
	}

	return nil
}

// Commit commits all changes in the specified article
func (r GitRepository) Commit(article int64) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// stage all files
	_, err = w.Add("./")
	if err != nil {
		return err
	}

	// commit
	_, err = w.Commit("version update", &git.CommitOptions{
		Author: &object.Signature{
			// TODO: add actual user name?
			Name:  "Alexandria Git Manager",
			Email: "",
			When:  r.Clock.Now(),
		},
	})
	return nil
}

func (r GitRepository) CheckoutBranch(article int64, version int64) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// checkout
	idString := strconv.FormatInt(article, 10)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(idString),
	})

	return nil
}

// GetArticlePath returns the path to an article git repository
func (r GitRepository) GetArticlePath(article int64) (string, error) {
	idString := strconv.FormatInt(article, 10)
	path, err := filepath.Abs(filepath.Join(r.Path, idString))
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}

// getWorktree returns the go-git worktree of an article git repository
func (r GitRepository) getWorktree(article int64) (*git.Worktree, error) {
	// Open  repository.
	dir, err := r.GetArticlePath(article)
	if err != nil {
		return nil, err
	}

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
