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
	"io/ioutil"
	"mainServer/entities"
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
func NewGitRepository(path string) (GitRepository, error) {

	// make folders for git files
	err := os.MkdirAll(filepath.Join(path, "persistent"), os.ModePerm)
	if err != nil {
		return GitRepository{}, err
	}
	err = os.MkdirAll(filepath.Join(path, "cache"), os.ModePerm)
	if err != nil {
		return GitRepository{}, err
	}

	return GitRepository{Path: path, Clock: clock.RealClock{}}, nil
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
	path, err := filepath.Abs(filepath.Join(r.Path, "persistent", idString))
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), err
}

// GetRequestCachePath returns the path to a cache folder for a request
// makes sure that the folder is created, including nested /old and /new folders
func (r GitRepository) GetRequestCachePath(article int64, sourceHistoryID string, targetHistoryID string) (string, error) {

	// get the path by generating a unique cache id
	id := fmt.Sprintf("%d-%s-%s", article, sourceHistoryID, targetHistoryID)
	path, err := filepath.Abs(filepath.Join(r.Path, "cache", "requests", id))
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
	commitHash, err := git2.RevParse(revparse.Args(versionStr), runGitIn(path))
	if err != nil {
		return "", err
	}
	return commitHash, nil
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

// RequestPreviewExists stores whether there is already a preview of the specified request,
// this requires the historyID's to be set in the req struct
func (r GitRepository) RequestPreviewExists(req entities.Request) (bool, error) {
	path, err := r.GetRequestCachePath(req.ArticleID, req.SourceHistoryID, req.TargetHistoryID)
	if err != nil {
		return false, err
	}

	// check if the file exists, using errors because there is no direct exists function from os
	// testing for file, because folder was automatically created by calling GetRequestCachePath
	_, err = os.Stat(filepath.Join(path, "new", "main.qmd"))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// StoreRequestPreview performs a merge without committing.
// Stores the before-and-after in a cache folder
// Requires the commit/history ID's to be specified in the req struct
// Might leave the repo behind with a detached HEAD
// Returns true iff there are no conflicts
func (r GitRepository) StoreRequestPreview(req entities.Request) (bool, error) {
	// get paths
	repo, err := r.GetArticlePath(req.ArticleID)
	if err != nil {
		return false, err
	}
	cache, err := r.GetRequestCachePath(req.ArticleID, req.SourceHistoryID, req.TargetHistoryID)
	if err != nil {
		return false, err
	}

	// checkout target commit, (possibly creating a detached head)
	res, err := git2.Checkout(checkout.Branch(req.TargetHistoryID), runGitIn(repo))
	if err != nil {
		return false, errors.New(res)
	}

	// copy files to "old" cache
	input, err := ioutil.ReadFile(filepath.Join(repo, "main.qmd"))
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(filepath.Join(cache, "old", "main.qmd"), input, 0644)
	if err != nil {
		return false, err
	}

	// merge source commit into target, without committing
	mergeRes, err := git2.Merge(merge.Commits(req.SourceHistoryID), merge.NoCommit, merge.NoFf, runGitIn(repo))
	conflicts := strings.Contains(mergeRes, "CONFLICT")
	if err != nil && !conflicts { // if err is just that there are conflicts, the execution can continue as normal
		return false, errors.New(mergeRes)
	}

	// copy merged files to "new" cache (possibly with conflicts)
	input, err = ioutil.ReadFile(filepath.Join(repo, "main.qmd"))
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(filepath.Join(cache, "new", "main.qmd"), input, 0644)
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
	return !conflicts, nil
}

// TODO: this helper func might be  redundant, because executing a merge will already inform about conflicts, remove func completely?
// hasConflicts checks if there are conflicts in an ongoing merge
// should only be used during a merge that has not been committed yet
func (r GitRepository) hasConflicts(article int64) (bool, error) {
	path, err := r.GetArticlePath(article)
	if err != nil {
		return false, err
	}

	// check for conflicts using raw execution, because go-git-cmd-wrapper does not support the diff command
	// should be secure even though it's raw, because the command doesn't take any user input
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	res := strings.TrimSuffix(string(output), "\n")
	if err != nil {
		return false, errors.New(res)
	}

	return res != "", nil
}

// GetRequestPreview returns the before and after main article file of a request
// requires the history ID's to be up-to-date in the req parameter
func (r GitRepository) GetRequestPreview(article int64, sourceHistoryID string, targetHistoryID string) (string, string, error) {
	// get paths
	cache, err := r.GetRequestCachePath(article, sourceHistoryID, targetHistoryID)
	if err != nil {
		return "", "", err
	}

	// read both old and new file from the cache
	oldFile, err := ioutil.ReadFile(filepath.Join(cache, "old", "main.qmd"))
	if err != nil {
		return "", "", err
	}
	newFile, err := ioutil.ReadFile(filepath.Join(cache, "new", "main.qmd"))
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
