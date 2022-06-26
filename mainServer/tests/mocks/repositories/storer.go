package repositories

import (
	"github.com/gin-gonic/gin"
	mocks "mainServer/tests/util"
	"mainServer/utils/clock"
	"mime/multipart"
)

type StorerMock struct {
	Mock *mocks.Mock
}

var GetVersionMockStorer func(article int64, version int64) (string, error)

func (m StorerMock) GetVersion(article int64, version int64) (string, error) {
	m.Mock.CallFunc("GetVersion", &map[string]interface{}{
		"article": article,
		"version": version,
	})
	return GetVersionMockStorer(article, version)
}

var GetVersionByCommitMock func(article int64, commit string) (string, error)

func (m StorerMock) GetVersionByCommit(article int64, commit string) (string, error) {
	m.Mock.CallFunc("GetVersionByCommit", &map[string]interface{}{
		"article": article,
		"commit":  commit,
	})
	return GetVersionByCommitMock(article, commit)
}

var GetVersionZippedMock func(article int64, version int64, filename string) (string, func(), error)

func (m StorerMock) GetVersionZipped(article int64, version int64, filename string) (string, func(), error) {
	m.Mock.CallFunc("GetVersionZipped", &map[string]interface{}{
		"article":  article,
		"version":  version,
		"filename": filename,
	})
	return GetVersionZippedMock(article, version, filename)
}

var UpdateAndCommitMock func(c *gin.Context, file *multipart.FileHeader, article int64, version int64) (string, error)

func (m StorerMock) UpdateAndCommit(c *gin.Context, file *multipart.FileHeader, article int64, version int64) (string, error) {
	m.Mock.CallFunc("UpdateAndCommit", &map[string]interface{}{
		"c":       c,
		"file":    file,
		"article": article,
		"version": version,
	})
	return UpdateAndCommitMock(c, file, article, version)
}

var CreateVersionFromMock func(article int64, source int64, target int64) (string, error)

func (m StorerMock) CreateVersionFrom(article int64, source int64, target int64) (string, error) {
	m.Mock.CallFunc("CreateVersionFrom", &map[string]interface{}{
		"article": article,
		"source":  source,
		"target":  target,
	})
	return CreateVersionFromMock(article, source, target)
}

var InitMainVersionMock func(article int64, mainVersion int64) (string, error)

func (m StorerMock) InitMainVersion(article int64, mainVersion int64) (string, error) {
	m.Mock.CallFunc("InitMainVersion", &map[string]interface{}{
		"article":     article,
		"mainVersion": mainVersion,
	})
	return InitMainVersionMock(article, mainVersion)
}

func (m StorerMock) SetClock(clock clock.Clock) {
	m.Mock.CallFunc("SetClock", &map[string]interface{}{
		"clock": clock,
	})
}

var StoreRequestComparisonMock func(article int64, request int64, source int64, target int64) (bool, error)

func (m StorerMock) StoreRequestComparison(article int64, request int64, source int64, target int64) (bool, error) {
	m.Mock.CallFunc("StoreRequestComparison", &map[string]interface{}{
		"article": article,
		"request": request,
		"source":  source,
		"target":  target,
	})
	return StoreRequestComparisonMock(article, request, source, target)
}

var GetRequestComparisonMock func(article int64, request int64) (string, string, error)

func (m StorerMock) GetRequestComparison(article int64, request int64) (string, string, error) {
	m.Mock.CallFunc("GetRequestComparison", &map[string]interface{}{
		"article": article,
		"request": request,
	})
	return GetRequestComparisonMock(article, request)
}

var MergeMock func(article int64, source int64, target int64) (string, error)

func (m StorerMock) Merge(article int64, source int64, target int64) (string, error) {
	m.Mock.CallFunc("Merge", &map[string]interface{}{
		"article": article,
		"source":  source,
		"target":  target,
	})
	return MergeMock(article, source, target)
}
