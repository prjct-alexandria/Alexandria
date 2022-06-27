package repositories

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type RequestThreadRepositoryMock struct {
	Mock *mocks.Mock
}

func NewRequestThreadRepositoryMock() RequestThreadRepositoryMock {
	return RequestThreadRepositoryMock{Mock: mocks.NewMock()}
}

var GetRequestThreadsMock func(aid int64, rid int64) ([]models.Thread, error)

func (m RequestThreadRepositoryMock) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	m.Mock.CallFunc("GetRequestThreads", &map[string]interface{}{
		"aid": aid,
		"rid": rid,
	})
	return GetRequestThreadsMock(aid, rid)
}

var CreateRequestThreadMock func(rid int64, tid int64) (int64, error)

func (m RequestThreadRepositoryMock) CreateRequestThread(rid int64, tid int64) (int64, error) {
	m.Mock.CallFunc("CreateRequestThread", &map[string]interface{}{
		"rid": rid,
		"tid": tid,
	})
	return CreateRequestThreadMock(rid, tid)
}
