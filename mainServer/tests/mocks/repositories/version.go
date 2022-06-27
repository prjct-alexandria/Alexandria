package repositories

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type VersionRepositoryMock struct {
	Mock *mocks.Mock
}

func NewVersionRepositoryMock() VersionRepositoryMock {
	return VersionRepositoryMock{Mock: mocks.NewMock()}
}

var CreateVersionMock func(version entities.Version) (entities.Version, error)

func (m VersionRepositoryMock) CreateVersion(version entities.Version) (entities.Version, error) {
	m.Mock.CallFunc("CreateVersion", &map[string]interface{}{
		"version": version,
	})
	return CreateVersionMock(version)
}

var GetVersionMock func(version int64) (entities.Version, error)

func (m VersionRepositoryMock) GetVersion(version int64) (entities.Version, error) {
	m.Mock.CallFunc("GetVersion", &map[string]interface{}{
		"version": version,
	})
	return GetVersionMock(version)
}

var GetVersionsByArticleMock func(article int64) ([]entities.Version, error)

func (m VersionRepositoryMock) GetVersionsByArticle(article int64) ([]entities.Version, error) {
	m.Mock.CallFunc("GetVersionsByArticle", &map[string]interface{}{
		"article": article,
	})
	return GetVersionsByArticleMock(article)
}

var CheckIfOwnerMock func(version int64, email string) (bool, error)

func (m VersionRepositoryMock) CheckIfOwner(version int64, email string) (bool, error) {
	m.Mock.CallFunc("CheckIfOwner", &map[string]interface{}{
		"version": version,
		"email":   email,
	})
	return CheckIfOwnerMock(version, email)
}

var UpdateVersionLatestCommitMock func(version int64, commit string) error

func (m VersionRepositoryMock) UpdateVersionLatestCommit(version int64, commit string) error {
	m.Mock.CallFunc("UpdateVersionLatestCommit", &map[string]interface{}{
		"version": version,
		"commit":  commit,
	})
	return UpdateVersionLatestCommitMock(version, commit)
}
