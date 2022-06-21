package storer

import (
	"github.com/gin-gonic/gin"
	"mainServer/repositories/storer/git"
	"mime/multipart"
)

// GetVersion returns the contents of the specified article version file as string
func (s *Storer) GetVersion(article int64, version int64) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path to the article repository
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout version
	repo := git.NewRepo(path)
	err = repo.CheckoutBranch(version)
	if err != nil {
		return "", err
	}

	// Return the article file contents as string
	return s.fs.ReadArticleFile(path)
}

// GetVersionByCommit returns the contents of the specified article at the specified commit
func (s *Storer) GetVersionByCommit(article int64, commit [20]byte) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path to the article repository
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout commit
	repo := git.NewRepo(path)
	err = repo.CheckoutCommit(commit)
	if err != nil {
		return "", err
	}

	// Return the article file contents as string
	return s.fs.ReadArticleFile(path)
}

// GetVersionZipped creates a .zip, with specified name, with all files of an article version at the returned path
func (s *Storer) GetVersionZipped(article int64, version int64, filename string) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path to the article repository
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout version
	repo := git.NewRepo(path)
	err = repo.CheckoutBranch(version)
	if err != nil {
		return "", err
	}

	// Add a .zip with the files to a cache folder and return the path
	path, err = s.fs.MakeDownloadZip(filename, path)
	if err != nil {
		return "", nil
	}
	return path, nil
}

// UpdateAndCommit overwrites the file content of the specified article version and committing changes.
// Returns the ID of the created commit.
func (s *Storer) UpdateAndCommit(c *gin.Context, file *multipart.FileHeader, article int64, version int64) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path to the article repository
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout version
	repo := git.NewRepo(path)
	err = repo.CheckoutBranch(version)
	if err != nil {
		return "", err
	}

	// Save file
	err = s.fs.SaveArticleFile(c, file, path)
	if err != nil {
		return "", err
	}

	// Commit
	err = repo.Commit(s.clock.Now())
	if err != nil {
		return "", err
	}

	// Return the id of the created commit
	return repo.GetLatestCommit(version)
}

// CreateVersionFrom creates a new target version based on a source version, does not modify any file contents
// Returns the ID of the created commit, which is the same as that of the source version
func (s *Storer) CreateVersionFrom(article int64, source int64, target int64) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Get the path to the article repository
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}

	// Checkout source version
	repo := git.NewRepo(path)
	err = repo.CheckoutBranch(source)
	if err != nil {
		return "", err
	}

	// Create the new branch
	err = repo.CreateBranch(source, target)
	if err != nil {
		return "", err
	}

	// Return the commit id of the created version
	return repo.GetLatestCommit(target)
}

// InitMainVersion creates a new git repository for the article, with an initial branch named after the given version
// Returns ID of the initial commit that added a default file.
func (s *Storer) InitMainVersion(article int64, mainVersion int64) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Create article repository folder
	path, err := s.fs.CreateArticlePath(article)
	if err != nil {
		return "", err
	}

	// Initialize git repository
	repo := git.NewRepo(path)
	err = repo.Init(mainVersion)
	if err != nil {
		return "", err
	}

	// Put the default file there and commit
	err = s.fs.PlaceDefaultFile(path)
	if err != nil {
		return "", err
	}
	err = repo.Commit(s.clock.Now())
	if err != nil {
		return "", err
	}

	// Return the commit id of the created version
	return repo.GetLatestCommit(mainVersion)
}
