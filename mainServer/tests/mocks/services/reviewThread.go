package services

import mocks "mainServer/tests/util"

type ReviewThreadServiceMock struct {
	Mock *mocks.Mock
}

func NewReviewThreadServiceMock() ReviewThreadServiceMock {
	return ReviewThreadServiceMock{Mock: mocks.NewMock()}
}

var StartReviewThreadMock func(rid int64, tid int64) (int64, error)

func (m ReviewThreadServiceMock) StartReviewThread(rid int64, tid int64) (int64, error) {
	m.Mock.CallFunc("StartReviewThread", &map[string]interface{}{
		"rid": rid,
		"tid": tid,
	})
	return StartReviewThreadMock(rid, tid)
}
