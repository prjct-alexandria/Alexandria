package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	git2 "github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/reset"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	"github.com/ldez/go-git-cmd-wrapper/v2/types"
	"mainServer/server/config"
	"mainServer/utils/clock"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type GitRepository struct {
	Path  string
	Clock clock.Clock
}

// NewGitRepository creates a new GitRepository class.
// This is NOT the function used to create a folder/git repository to store an article in.
// See CreateRepo instead
func NewGitRepository(cfg *config.GitConfig) GitRepository {

	// make folder for git files
	err := os.MkdirAll(cfg.Path, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return GitRepository{Path: cfg.Path, Clock: clock.RealClock{}}
}

// CreateRepo creates a new folder/git repository to store an article in, including main version branch.
// This is NOT the function used to create the Go repository class,
// see NewGitRepository instead,
func (r GitRepository) CreateRepo(article int64, version int64) error {
	path, err := r.GetArticlePath(article)
	if err != nil {
		return err
	}

	// Check if folder with the same name already exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("trying to create a git repository that already exists with id=%d", article)
	}

	// Create directory
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	// Git init, Go-git automatically creates a branch called "master"
	_, err = git.PlainInit(path, false)
	if err != nil {
		return err
	}

	// Rename the branch to the main version ID
	err = r.renameInitialBranch(article, version)
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

// CheckoutCommit checks out the specified commit in the specified article repo
func (r GitRepository) CheckoutCommit(article int64, commit [20]byte) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// checkout
	err = w.Checkout(&git.CheckoutOptions{
		Hash: commit,
	})

	return err
}

// CheckoutBranch checks out the specified version in the specified article repo
func (r GitRepository) CheckoutBranch(article int64, version int64) error {
	w, err := r.getWorktree(article)
	if err != nil {
		return err
	}

	// checkout
	branchName := strconv.FormatInt(version, 10)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})

	return err
}

// CreateBranch creates a new branch based on the source one, named as target. Will automatically check out source branch.
func (r GitRepository) CreateBranch(article int64, source int64, target int64) error {

	// Open repository and get worktree
	dir, err := r.GetArticlePath(article)
	if err != nil {
		return err
	}
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Checkout source branch
	sourceName := strconv.FormatInt(source, 10)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(sourceName),
	})
	if err != nil {
		return err
	}

	// Create new branch reference
	targetName := strconv.FormatInt(target, 10)
	targetRef := plumbing.NewBranchReferenceName(targetName)
	if err != nil {
		return err
	}

	// Store the new branch reference to head
	headRef, err := repo.Head()
	ref := plumbing.NewHashReference(targetRef, headRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		return err
	}
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

// renameInitialBranch manually renames the initial branch on a new repo.
// Do not use on git repos that already have commits!
// The main branch must have a database-usable version ID, not "master"
// So, a file in .git is edited manually
func (r GitRepository) renameInitialBranch(article int64, version int64) error {
	path, err := r.GetArticlePath(article)
	if err != nil {
		return err
	}

	headPath := filepath.Join(path, ".git", "HEAD")
	f, err := os.Create(headPath)
	if err != nil {
		_ = f.Close()
		return err
	}

	_, err = fmt.Fprintf(f, "ref: refs/heads/%d", version)
	if err != nil {
		_ = f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// GetLatestCommit returns the commit ID of the latest commit on the specified article version
func (r GitRepository) GetLatestCommit(article int64, version int64) (string, error) {
	path, err := r.GetArticlePath(article)
	if err != nil {
		return "", err
	}
	versionStr := strconv.FormatInt(version, 10)

	// call the git command rev-parse, which returns a commit hash when given a branch name
	output, err := git2.RevParse(revparse.Args(versionStr), runGitIn(path))
	if err != nil {
		return "", errors.New(output)
	}

	return output, nil
}

// Merge merges the source branch (version) into the target branch (version) of the specified repository (article)
func (r GitRepository) Merge(article int64, source int64, target int64) error {
	sourceStr := strconv.FormatInt(source, 10)
	targetStr := strconv.FormatInt(target, 10)

	path, err := r.GetArticlePath(article)
	if err != nil {
		return err
	}

	// checkout target
	res, err := git2.Checkout(checkout.Branch(targetStr), runGitIn(path))
	if err != nil {
		return errors.New(res)
	}

	// revert any possible un-committed changes left over from previous failed executions
	// without a clean worktree, aborting a failed merge might not be possible
	// should be redundant, but the guarantee is good to have
	res, err = git2.Reset(reset.Hard, runGitIn(path))
	if err != nil {
		return errors.New(res)
	}

	// merge source into target and commit immediately
	res, err = git2.Merge(merge.Commits(sourceStr), merge.Commit, runGitIn(path))
	if err != nil {
		res2, err := git2.Merge(merge.Abort, runGitIn(path))
		if err != nil {
			return fmt.Errorf("failed aborting merge with message: %s, merge itself failed with error %s", res, res2)
		}
		return errors.New(res)
	}
	return nil
}

// custom option made for use with the go-git-cmd-wrapper library,
// enables execution in specific paths, without using os change dir, which possibly interferes with other operations
func runGitIn(path string) types.Option {
	return git2.CmdExecutor(
		func(ctx context.Context, name string, debug bool, args ...string) (string, error) {
			cmd := exec.CommandContext(ctx, name, args...)
			cmd.Dir = path

			output, err := cmd.CombinedOutput()
			return strings.TrimSuffix(string(output), "\n"), err
		},
	)
}
