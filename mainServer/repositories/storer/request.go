package storer

import (
	"errors"
	"fmt"
	"mainServer/repositories/storer/git"
)

// StoreRequestComparison performs a merge without committing,
// storing the before and after in the file system.
// Returns whether there were conflicts while merging
func (s *Storer) StoreRequestComparison(article int64, request int64, source int64, target int64) (bool, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the paths to the article repository and request before-and-after location
	articlePath, err := s.fs.GetArticlePath(article)
	if err != nil {
		return false, err
	}
	requestPath, err := s.fs.GetRequestPath(article, request)
	if err != nil {
		return false, err
	}

	// Checkout the target branch
	repo := git.NewRepo(articlePath)
	err = repo.CheckoutBranch(target)

	// Store the files as the "old" snapshot
	err = s.fs.CopyArticleToRequest(articlePath, requestPath, false)
	if err != nil {
		return false, err
	}

	// Merge without committing
	conflicted, err := repo.Merge(source)
	defer repo.Abort()

	// Store the files as the "new" snapshot
	err = s.fs.CopyArticleToRequest(articlePath, requestPath, true)
	if err != nil {
		return false, err
	}

	return conflicted, err
}

// GetRequestComparison returns the before and after main article file of a request
// returns (before, after, err)
func (s *Storer) GetRequestComparison(article int64, request int64) (string, string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path
	requestPath, err := s.fs.GetRequestPath(article, request)
	if err != nil {
		return "", "", err
	}

	// Read the contents of the path
	return s.fs.GetRequestComparison(requestPath)
}

// Merge executes a merge, returning the resulting merge commit id
func (s *Storer) Merge(article int64, source int64, target int64) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the paths to the article repository
	articlePath, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout the target branch
	repo := git.NewRepo(articlePath)
	err = repo.CheckoutBranch(target)

	// Merge without committing
	conflicted, err := repo.Merge(source)
	if err != nil {
		return "", err
	}

	// Abort when conflicted
	if conflicted {
		repo.Abort()
		return "", errors.New("trying to merge with conflicts")
	}

	// Commit
	err = repo.Commit(s.clock.Now(), fmt.Sprintf("Merge %v into %v", source, target))
	if err != nil {
		return "", err
	}

	// Return the id of the created commit
	return repo.GetLatestCommit(target)
}
