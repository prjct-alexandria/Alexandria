package services

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type CommentServiceMock struct {
	Mock *mocks.Mock
}

func NewCommentServiceMock() CommentServiceMock {
	return CommentServiceMock{Mock: mocks.NewMock()}
}

var SaveCommentMock func(comment entities.Comment, tid int64, loggedInAs string) (int64, error)

func (m CommentServiceMock) SaveComment(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
	m.Mock.CallFunc("SaveComment", &map[string]interface{}{
		"comment":    comment,
		"tid":        tid,
		"loggedInAs": loggedInAs,
	})
	return SaveCommentMock(comment, tid, loggedInAs)
}
