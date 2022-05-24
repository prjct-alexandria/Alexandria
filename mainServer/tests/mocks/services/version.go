package services

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

// UpdateVersionMock is a declared function whose behaviour can be modified by individual tests
var UpdateVersionMock func(c *gin.Context, file *multipart.FileHeader, article string, version string) error

// VersionServiceMock mocks class using publicly variable mock functions
type VersionServiceMock struct {
	// mock tracks what functions were called
	UpdateVersionCalled bool
	UpdateVersionParams map[string]interface{}
}

func (m VersionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
	m.UpdateVersionCalled = true
	m.UpdateVersionParams = map[string]interface{}{
		"article": article,
		"version": version,
	}
	return nil //UpdateVersionMock(c, file, article, version)
}
