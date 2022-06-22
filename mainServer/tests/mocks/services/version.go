package services

import (
	"github.com/gin-gonic/gin"
	"mainServer/models"
	"mime/multipart"
)

// UpdateVersionMock is a declared function whose behaviour can be modified by individual tests
var UpdateVersionMock func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error

// VersionServiceMock mocks class using publicly modifiable mock functions
type VersionServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

func (m VersionServiceMock) ListVersions(article int64) ([]models.Version, error) {
	// added to solve merge conflicts after the testing issue was finished
	//TODO implement me
	panic("implement me")
}

func (m VersionServiceMock) GetVersion(article int64, version int64) (models.Version, error) {
	// added to solve merge conflicts after the testing issue was finished
	// TODO: implement for future testing
	panic("implement me")
}

func (m VersionServiceMock) GetVersionByCommitID(article int64, version int64, commit string) (models.Version, error) {
	// added to solve merge conflicts after the testing issue was finished
	// TODO: implement for future testing
	panic("implement me")
}

func (m VersionServiceMock) CreateVersionFrom(article int64, source int64, title string, owners []string) (models.Version, error) {
	//TODO implement me
	panic("implement me")
}

func (m VersionServiceMock) GetVersionFiles(aid int64, vid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

// NewVersionServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewVersionServiceMock() VersionServiceMock {
	return VersionServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

// UpdateVersion implements the corresponding version from the VersionService interface.
// Stores in the mock that it was called, including the arguments, and executes the custom UpdateVersionMock function
func (m VersionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
	(*m.Called)["UpdateVersion"] = true
	(*m.Params)["UpdateVersion"] = map[string]interface{}{
		"article": article,
		"version": version,
	}
	return UpdateVersionMock(c, file, article, version)
}
