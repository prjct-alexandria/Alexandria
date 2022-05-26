package interfaces

import (
	"github.com/gin-gonic/gin"
	"mainServer/entities"
	"mime/multipart"
)

type VersionService interface {

	// UpdateVersion overwrites file of specified article version and commits
	UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error

	// GetVersion looks for a version in the filesystem and creates a version entity from it with the appropriate metadata.
	GetVersion(c *gin.Context, article string, version string) (entities.Version, error)
}
