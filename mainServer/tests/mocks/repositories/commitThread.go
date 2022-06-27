package repositories

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type CommitThreadRepositoryMock struct {
	Mock *mocks.Mock
}

func NewCommitThreadRepositoryMock() CommitThreadRepositoryMock {
	return CommitThreadRepositoryMock{Mock: mocks.NewMock()}
}

var GetCommitThreadsMock func(aid int64, cid string) ([]models.Thread, error)

func (m CommitThreadRepositoryMock) GetCommitThreads(aid int64, cid string) ([]models.Thread, error) {
	m.Mock.CallFunc("GetCommitThreads", &map[string]interface{}{
		"aid": aid,
		"cid": cid,
	})
	return GetCommitThreadsMock(aid, cid)
}

var CreateCommitThreadMock func(cid string, tid int64) (int64, error)

func (m CommitThreadRepositoryMock) CreateCommitThread(cid string, tid int64) (int64, error) {
	m.Mock.CallFunc("CreateCommitThread", &map[string]interface{}{
		"cid": cid,
		"tid": tid,
	})
	return CreateCommitThreadMock(cid, tid)
}
