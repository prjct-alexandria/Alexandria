package services

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type CommitSelectionThreadServiceMock struct {
	Mock *mocks.Mock
}

func NewCommitSelectionThreadServiceMock() CommitSelectionThreadServiceMock {
	return CommitSelectionThreadServiceMock{Mock: mocks.NewMock()}
}

var StartCommitSelectionThreadMock func(cid string, tid int64, section string) (int64, error)

func (m CommitSelectionThreadServiceMock) StartCommitSelectionThread(cid string, tid int64, section string) (int64, error) {
	m.Mock.CallFunc("StartCommitSelectionThread", &map[string]interface{}{
		"cid":     cid,
		"tid":     tid,
		"section": section,
	})
	return StartCommitSelectionThreadMock(cid, tid, section)
}

var GetCommitSelectionThreadsMock func(sid string, aid int64) ([]models.SelectionThread, error)

func (m CommitSelectionThreadServiceMock) GetCommitSelectionThreads(sid string, aid int64) ([]models.SelectionThread, error) {
	m.Mock.CallFunc("GetCommitSelectionThreads", &map[string]interface{}{
		"sid": sid,
		"aid": aid,
	})
	return GetCommitSelectionThreadsMock(sid, aid)
}
