package interfaces

import (
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mime/multipart"
)

type VersionService interface {

	// UpdateVersion overwrites file of specified article version and commits
	UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error

	// GetVersion looks for a version in the filesystem and creates a version entity from it with the appropriate metadata.
	GetVersion(article int64, version int64) (models.Version, error)

	// ListVersions returns a list of models for all versions of the specified article
	ListVersions(article int64) ([]models.Version, error)
}
