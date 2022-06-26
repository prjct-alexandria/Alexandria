package services

import (
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type ThreadServiceMock struct {
	Mock *mocks.Mock
}

func NewThreadServiceMock() ThreadServiceMock {
	return ThreadServiceMock{Mock: mocks.NewMock()}
}

var StartThreadMock func(thread models.Thread, aid int64) (int64, error)

func (m ThreadServiceMock) StartThread(thread models.Thread, aid int64) (int64, error) {
	m.Mock.CallFunc("StartThread", &map[string]interface{}{
		"thread": thread,
		"aid":    aid,
	})
	return StartThreadMock(thread, aid)
}
