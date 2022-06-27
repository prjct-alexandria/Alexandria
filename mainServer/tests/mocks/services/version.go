package services

import (
	"github.com/gin-gonic/gin"
	"mainServer/models"
	mocks "mainServer/tests/util"
	"mime/multipart"
)

type VersionServiceMock struct {
	Mock *mocks.Mock
}

func NewVersionServiceMock() VersionServiceMock {
	return VersionServiceMock{Mock: mocks.NewMock()}
}

var UpdateVersionMock func(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error

func (m VersionServiceMock) UpdateVersion(c *gin.Context, file *multipart.FileHeader, article int64, version int64, loggedInAs string) error {
	m.Mock.CallFunc("UpdateVersion", &map[string]interface{}{
		"c":          c,
		"file":       file,
		"article":    article,
		"version":    version,
		"loggedInAs": loggedInAs,
	})
	return UpdateVersionMock(c, file, article, version, loggedInAs)
}

var GetVersionMock func(article int64, version int64) (models.Version, error)

func (m VersionServiceMock) GetVersion(article int64, version int64) (models.Version, error) {
	m.Mock.CallFunc("GetVersion", &map[string]interface{}{
		"article": article,
		"version": version,
	})
	return GetVersionMock(article, version)
}

var GetVersionByCommitIDMock func(article int64, version int64, commit string) (models.Version, error)

func (m VersionServiceMock) GetVersionByCommitID(article int64, version int64, commit string) (models.Version, error) {
	m.Mock.CallFunc("GetVersionByCommitID", &map[string]interface{}{
		"article": article,
		"version": version,
		"commit":  commit,
	})
	return GetVersionByCommitIDMock(article, version, commit)
}

var CreateVersionFromMock func(article int64, source int64, title string, owners []string, loggedInAs string) (models.Version, error)

func (m VersionServiceMock) CreateVersionFrom(article int64, source int64, title string, owners []string, loggedInAs string) (models.Version, error) {
	m.Mock.CallFunc("CreateVersionFrom", &map[string]interface{}{
		"article":    article,
		"source":     source,
		"title":      title,
		"owners":     owners,
		"loggedInAs": loggedInAs,
	})
	return CreateVersionFromMock(article, source, title, owners, loggedInAs)
}

var ListVersionsMock func(article int64) ([]models.Version, error)

func (m VersionServiceMock) ListVersions(article int64) ([]models.Version, error) {
	m.Mock.CallFunc("ListVersions", &map[string]interface{}{
		"article": article,
	})
	return ListVersionsMock(article)
}

var GetVersionFilesMock func(aid int64, vid int64) (string, func(), error)

func (m VersionServiceMock) GetVersionFiles(aid int64, vid int64) (string, func(), error) {
	m.Mock.CallFunc("GetVersionFiles", &map[string]interface{}{
		"aid": aid,
		"vid": vid,
	})
	return GetVersionFilesMock(aid, vid)
}
