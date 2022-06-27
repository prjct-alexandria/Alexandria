package services

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type CommitThreadServiceMock struct {
	Mock *mocks.Mock
}

func NewCommitThreadServiceMock() CommitThreadServiceMock {
	return CommitThreadServiceMock{Mock: mocks.NewMock()}
}

var StartCommitThreadMock func(cid string, tid int64) (int64, error)

func (m CommitThreadServiceMock) StartCommitThread(cid string, tid int64) (int64, error) {
	m.Mock.CallFunc("StartCommitThread", &map[string]interface{}{
		"cid": cid,
		"tid": tid,
	})
	return StartCommitThreadMock(cid, tid)
}

var GetCommitThreadsMock func(aid int64, cid string) ([]models.Thread, error)

func (m CommitThreadServiceMock) GetCommitThreads(aid int64, cid string) ([]models.Thread, error) {
	m.Mock.CallFunc("GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
	return GetCommitThreadsMock(aid, cid)
}
