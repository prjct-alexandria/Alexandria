package interfaces

import (
	"github.com/gin-gonic/gin"
	"mainServer/utils/clock"
	"mime/multipart"
)

type Storer interface {

	// GetVersion returns the contents of the specified article version file as string
	GetVersion(article int64, version int64) (string, error)

	// GetVersionByCommit returns the contents of the specified article at the specified commit
	GetVersionByCommit(article int64, commit string) (string, error)

	// GetVersionZipped creates a .zip, with specified name, with all files of an article version at the returned path
	// includes function that can be used to delete the zip after it's not needed anymore
	GetVersionZipped(article int64, version int64, filename string) (string, func(), error)

	// UpdateAndCommit overwrites the file content of the specified article version and committing changes.
	// Returns the ID of the created commit.
	UpdateAndCommit(c *gin.Context, file *multipart.FileHeader, article int64, version int64) (string, error)

	// CreateVersionFrom creates a new target version based on a source version, does not modify any file contents
	// Returns the ID of the created commit, which is the same as that of the source version
	CreateVersionFrom(article int64, source int64, target int64) (string, error)

	// InitMainVersion creates a new git repository for the article, with an initial branch named after the given version
	// Returns ID of the initial commit that added a default file.
	InitMainVersion(article int64, mainVersion int64) (string, error)

	// SetClock changes the clock that is used to set creation dates of commits
	// Can be used to supply a mock for testing purposes
	SetClock(clock clock.Clock)

	// StoreRequestComparison performs a merge without committing,
	// storing the before and after in the file system.
	// Returns whether there were conflicts while merging
	StoreRequestComparison(article int64, request int64, source int64, target int64) (bool, error)

	// GetRequestComparison returns the before and after main article file of a request
	// returns (before, after, err)
	GetRequestComparison(article int64, request int64) (string, string, error)

	// Merge executes a merge, returning the resulting merge commit id
	Merge(article int64, source int64, target int64) (string, error)
}
