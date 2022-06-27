package repositories

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type CommitSelectionThreadRepositoryMock struct {
	Mock *mocks.Mock
}

func NewCommitSelectionThreadRepositoryMock() CommitSelectionThreadRepositoryMock {
	return CommitSelectionThreadRepositoryMock{Mock: mocks.NewMock()}
}

var GetCommitSelectionThreadsMock func(aid int64, cid string) ([]models.SelectionThread, error)

func (m CommitSelectionThreadRepositoryMock) GetCommitSelectionThreads(aid int64, cid string) ([]models.SelectionThread, error) {
	m.Mock.CallFunc("GetCommitSelectionThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
	return GetCommitSelectionThreadsMock(aid, cid)
}

var CreateCommitSelectionThreadMock func(cid string, tid int64, section string) (int64, error)

func (m CommitSelectionThreadRepositoryMock) CreateCommitSelectionThread(cid string, tid int64, section string) (int64, error) {
	m.Mock.CallFunc("CreateCommitSelectionThread", &map[string]interface{}{
		"cid":     cid,
		"tid":     tid,
		"section": section,
	})
	return CreateCommitSelectionThreadMock(cid, tid, section)
}
