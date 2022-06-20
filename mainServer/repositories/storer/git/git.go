package git

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
	"io/ioutil"
	"mainServer/entities"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Repo struct {
	path   string
	runner types.Option
}

// NewRepo creates a new GitRepo class. This references a git repository that represents an article.
// This function does not actually initialize a repository.
// This is NOT the function used to create a folder/git repository to store an article in.
// See CreateRepo instead. It
func NewRepo(path string) Repo {
	return Repo{path: path, runner: runGitIn(path)}
}

// Init initializes
func (r Repo) Init(mainVersion int64) error {

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
	err = r.renameInitialBranch(mainVersion)
	if err != nil {
		return err
	}

	return nil
}

// Commit commits all changes in the specified article
func (r Repo) Commit(article int64) error {
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
func (r Repo) CheckoutCommit(article int64, commit [20]byte) error {
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
func (r Repo) CheckoutBranch(version int64) error {
	w, err := r.getWorktree()

	// checkout
	branchName := strconv.FormatInt(version, 10)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})

	return err
}

// CreateBranch creates a new branch based on the source one, named as target. Will automatically check out source branch.
func (r Repo) CreateBranch(article int64, source int64, target int64) error {

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

// getWorktree returns the go-git worktree of an article git repository
func (r Repo) getWorktree() (*git.Worktree, error) {
	repo, err := git.PlainOpen(r.path)
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return w, nil
}

// GetLatestCommit returns the commit ID of the latest commit on the specified article version
func (r Repo) GetLatestCommit(version int64) (string, error) {
	versionStr := strconv.FormatInt(version, 10)

	// call the git command rev-parse, which returns a commit hash when given a branch name
	output, err := git2.RevParse(revparse.Args(versionStr), r.runner)
	if err != nil {
		return "", errors.New(output)
	}

	return output, nil
}

// Merge merges the source branch (version) into the target branch (version) of the specified repository (article)
func (r Repo) Merge(source int64, target int64) error {
	sourceStr := strconv.FormatInt(source, 10)
	targetStr := strconv.FormatInt(target, 10)

	// checkout target
	res, err := git2.Checkout(checkout.Branch(targetStr), runGitIn(r.path))
	if err != nil {
		return errors.New(res)
	}

	// revert any possible un-committed changes left over from previous failed executions
	// without a clean worktree, aborting a failed merge might not be possible
	// should be redundant, but the guarantee is good to have
	res, err = git2.Reset(reset.Hard, r.runner)
	if err != nil {
		return errors.New(res)
	}

	// merge source into target and commit immediately
	res, err = git2.Merge(merge.Commits(sourceStr), merge.Commit, r.runner)
	if err != nil {
		res2, err := git2.Merge(merge.Abort, r.runner)
		if err != nil {
			return fmt.Errorf("failed aborting merge with message: %s, merge itself failed with error %s", res, res2)
		}
		return errors.New(res)
	}
	return nil
}

// StoreRequestComparison performs a merge without committing.
// Stores the before-and-after in a cache folder
// Requires the commit/history ID's to be specified in the req struct
// Might leave the repo behind with a detached HEAD
// Returns whether there are conflicts
func (r Repo) StoreRequestComparison(req entities.Request) (bool, error) {
	// get paths
	repo, err := r.GetArticlePath(req.ArticleID)
	if err != nil {
		return false, err
	}
	comparison, err := r.GetRequestComparisonPath(req.ArticleID, req.RequestID)
	if err != nil {
		return false, err
	}

	// checkout target commit, (possibly creating a detached head)
	res, err := git2.Checkout(checkout.Branch(req.TargetHistoryID), r.runner)
	if err != nil {
		return false, errors.New(res)
	}

	// copy files to "old" cache
	input, err := ioutil.ReadFile(filepath.Join(repo, "main.qmd"))
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(filepath.Join(comparison, "old", "main.qmd"), input, 0644)
	if err != nil {
		return false, err
	}

	// merge source commit into target, without committing
	mergeRes, err := git2.Merge(merge.Commits(req.SourceHistoryID), merge.NoCommit, merge.NoFf, r.runner)
	conflicts := strings.Contains(mergeRes, "CONFLICT")
	if err != nil && !conflicts { // if err is just that there are conflicts, the execution can continue as normal
		return false, errors.New(mergeRes)
	}

	// copy merged files to "new" cache (possibly with conflicts)
	input, err = ioutil.ReadFile(filepath.Join(repo, "main.qmd"))
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(filepath.Join(comparison, "new", "main.qmd"), input, 0644)
	if err != nil {
		return false, err
	}

	// abort merge / revert
	if mergeRes != "Already up to date." {
		res, err = git2.Merge(merge.Abort, runGitIn(repo))
		if err != nil {
			return false, errors.New(res)
		}
	}
	return conflicts, nil
}

// GetRequestComparison returns the before and after main article file of a request
// requires the history ID's to be up-to-date in the req parameter
func (r Repo) GetRequestComparison(article int64, request int64) (string, string, error) {
	// get paths
	path, err := r.GetRequestComparisonPath(article, request)
	if err != nil {
		return "", "", err
	}

	// read both old and new file from the cache
	oldFile, err := ioutil.ReadFile(filepath.Join(path, "old", "main.qmd"))
	if err != nil {
		return "", "", err
	}
	newFile, err := ioutil.ReadFile(filepath.Join(path, "new", "main.qmd"))
	if err != nil {
		return "", "", err
	}

	return string(oldFile), string(newFile), nil
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
