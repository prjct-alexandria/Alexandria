package services

import (
	"github.com/gin-gonic/gin"
	"mainServer/models"
	mocks "mainServer/tests/util"
	"mime/multipart"
)

// VersionServiceMock mocks class using publicly modifiable mock functions
type VersionServiceMock struct {
	Mock *mocks.Mock
}

// NewVersionServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewVersionServiceMock() VersionServiceMock {
	return VersionServiceMock{Mock: mocks.NewMock()}
}

// UpdateVersionMock is a declared function whose behaviour can be modified by individual tests
var UpdateVersionMock func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error

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

func (m VersionServiceMock) GetVersionFiles(aid int64, vid int64) (string, func(), error) {
	//TODO implement me
	panic("implement me")
}

// UpdateVersion implements the corresponding version from the VersionService interface.
// Stores in the mock that it was called, including the arguments, and executes the custom UpdateVersionMock function
func (m VersionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64) error {
	m.Mock.CallFunc("UpdateVersion", &map[string]interface{}{
		"article": article,
		"version": version,
	})
	return UpdateVersionMock(c, file, article, version)
}
