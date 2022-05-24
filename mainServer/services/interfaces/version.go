package interfaces

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type VersionService interface {

	// UpdateVersion overwrites file of specified article version and commits
	UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error
}
