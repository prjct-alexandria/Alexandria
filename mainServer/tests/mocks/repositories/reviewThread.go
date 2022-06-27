package repositories

import mocks "mainServer/tests/util"

type ReviewThreadRepositoryMock struct {
	Mock *mocks.Mock
}

func NewReviewThreadRepositoryMock() ReviewThreadRepositoryMock {
	return ReviewThreadRepositoryMock{Mock: mocks.NewMock()}
}

var CreateReviewThreadMock func(rid int64, tid int64) (int64, error)

func (m ReviewThreadRepositoryMock) CreateReviewThread(rid int64, tid int64) (int64, error) {
	m.Mock.CallFunc("CreateReviewThread", &map[string]interface{}{
		"rid": rid,
		"tid": tid,
	})
	return CreateReviewThreadMock(rid, tid)
}
