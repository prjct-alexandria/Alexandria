package services

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

// UpdateVersionMock is a declared function whose behaviour can be modified by individual tests
var UpdateVersionMock func(c *gin.Context, file *multipart.FileHeader, article string, version string) error

// VersionServiceMock mocks class using publicly modifiable mock functions
type VersionServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

// NewVersionServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewVersionServiceMock() VersionServiceMock {
	return VersionServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

func (m VersionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article string, version string) error {
	(*m.Called)["UpdateVersion"] = true
	(*m.Params)["UpdateVersion"] = map[string]interface{}{
		"article": article,
		"version": version,
	}
	return UpdateVersionMock(c, file, article, version)
}
