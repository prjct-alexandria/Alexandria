package repositories

import mocks "mainServer/tests/util"

type ThreadRepositoryMock struct {
	Mock *mocks.Mock
}

func NewThreadRepositoryMock() ThreadRepositoryMock {
	return ThreadRepositoryMock{Mock: mocks.NewMock()}
}

var CreateThreadMock func(aid int64) (int64, error)

func (m ThreadRepositoryMock) CreateThread(aid int64) (int64, error) {
	m.Mock.CallFunc("CreateThread", &map[string]interface{}{
		"aid": aid,
	})
	return CreateThreadMock(aid)
}
