package git

import (
	"context"
	"errors"
	"github.com/ldez/go-git-cmd-wrapper/v2/add"
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	"github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	"github.com/ldez/go-git-cmd-wrapper/v2/types"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Repo struct {
	path string
}

// NewRepo creates a new GitRepo class. This references a git repository that represents an article.
// This function does not actually initialize a repository or verify that it exists
func NewRepo(path string) Repo {
	return Repo{path: path}
}

// Init initializes a repository with the specified id as main branch name
func (r Repo) Init(mainVersion int64) error {

	// Git init
	output, err := git.Init(runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}

	// Checkout a new branch immediately before committing, this renames the main branch
	branchName := strconv.FormatInt(mainVersion, 10)
	output, err = git.Checkout(checkout.NewBranch(branchName), runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}

	return nil
}

// Commit commits all changes in the specified article, with commit message
func (r Repo) Commit(timestamp time.Time, msg string) error {

	// stage all files
	output, err := git.Add(add.All, runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}

	// commit
	unixTimeStr := strconv.FormatInt(timestamp.Unix(), 10)
	output, err = git.Commit(
		commit.Message(msg),
		commit.Author("Alexandria Git Manager"),
		commit.Date(unixTimeStr),
		runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}

	return err
}

// CheckoutCommit checks out the specified commit
func (r Repo) CheckoutCommit(commit string) error {
	output, err := git.Checkout(checkout.Branch(commit), runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}
	return nil
}

// CheckoutBranch checks out the specified version
func (r Repo) CheckoutBranch(version int64) error {
	// with this git implementation, the same as specific commit
	name := strconv.FormatInt(version, 10)
	return r.CheckoutCommit(name)
}

// CreateBranch creates a new branch based on the source one, named as target. Will automatically check out source branch.
func (r Repo) CreateBranch(source int64, target int64) error {

	// checkout source branch
	err := r.CheckoutBranch(source)
	if err != nil {
		return err
	}

	// create new branch from here
	targetName := strconv.FormatInt(target, 10)
	output, err := git.Checkout(checkout.NewBranch(targetName), runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}
	return nil
}

// GetLatestCommit returns the commit ID of the latest commit on the specified article version
func (r Repo) GetLatestCommit(version int64) (string, error) {
	versionStr := strconv.FormatInt(version, 10)

	// call the git command rev-parse, which returns a commit hash when given a branch name
	output, err := git.RevParse(revparse.Args(versionStr), runGitIn(r.path))
	if err != nil {
		return "", errors.New(output)
	}

	return output, nil
}

// Merge merges the source branch (version) into the currently checked out branch without committing
// returns whether there are conflicts, no error is returned if there are conflicts
func (r Repo) Merge(source int64) (bool, error) {
	sourceStr := strconv.FormatInt(source, 10)

	// merge source into current branch
	output, err := git.Merge(merge.Commits(sourceStr), merge.NoCommit, merge.NoFf, runGitIn(r.path))
	conflicts := strings.Contains(output, "CONFLICT")
	if err != nil && !conflicts {
		return false, nil
	}
	return conflicts, nil
}

// Abort aborts an ongoing merge, does not verify if a merge is actually going on
func (r Repo) Abort() error {
	output, err := git.Merge(merge.Abort, runGitIn(r.path))
	if err != nil {
		return errors.New(output)
	}
	return nil
}

// custom option made for use with the go-git-cmd-wrapper library,
// enables execution in specific paths, without using os change dir, which possibly interferes with other operations
func runGitIn(path string) types.Option {
	return git.CmdExecutor(
		func(ctx context.Context, name string, debug bool, args ...string) (string, error) {
			cmd := exec.CommandContext(ctx, name, args...)
			cmd.Dir = path

			output, err := cmd.CombinedOutput()
			return strings.TrimSuffix(string(output), "\n"), err
		},
	)
}
