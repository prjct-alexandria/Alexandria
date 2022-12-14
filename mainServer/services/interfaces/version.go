package interfaces

import (
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mime/multipart"
)

type VersionService interface {

	// UpdateVersion overwrites file of specified article version and commits
	UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error

	// GetVersion looks for a version in the filesystem and creates a version entity from it with the appropriate metadata.
	GetVersion(article int64, version int64) (models.Version, error)

	// GetVersionByCommitID does the same as GetVersion, but with a specific history/commit ID
	GetVersionByCommitID(article int64, version int64, commit string) (models.Version, error)

	// CreateVersionFrom makes a new version, based of an existing one. Version content is ignored in return value
	CreateVersionFrom(article int64, source int64, title string, owners []string, loggedInAs string) (models.Version, error)

	// ListVersions returns a list of models for all versions of the specified article
	ListVersions(article int64) ([]models.Version, error)

	// GetVersionFiles returns the path to a zip with all files of a version (except for the .git folder)
	// Includes cleanup function to delete temporary files when done
	GetVersionFiles(aid int64, vid int64) (string, func(), error)
}
