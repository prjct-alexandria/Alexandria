package services

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type RequestThreadServiceMock struct {
	Mock *mocks.Mock
}

func NewRequestThreadServiceMock() RequestThreadServiceMock {
	return RequestThreadServiceMock{Mock: mocks.NewMock()}
}

var StartRequestThreadMock func(rid int64, tid int64, loggedInAs string) (int64, error)

func (m RequestThreadServiceMock) StartRequestThread(rid int64, tid int64, loggedInAs string) (int64, error) {
	m.Mock.CallFunc("StartRequestThread", &map[string]interface{}{
		"rid":        rid,
		"tid":        tid,
		"loggedInAs": loggedInAs,
	})
	return StartRequestThreadMock(rid, tid, loggedInAs)
}

var GetRequestThreadsMock func(aid int64, rid int64) ([]models.Thread, error)

func (m RequestThreadServiceMock) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	m.Mock.CallFunc("GetRequestThreads", &map[string]interface{}{
		"aid": aid,
		"rid": rid,
	})
	return GetRequestThreadsMock(aid, rid)
}
