package repositories

import (
	"mainServer/entities"
	mocks "mainServer/tests/util"
)

type CommentRepositoryMock struct {
	Mock *mocks.Mock
}

func NewCommentRepositoryMock() CommentRepositoryMock {
	return CommentRepositoryMock{Mock: mocks.NewMock()}
}

var SaveCommentMock func(id entities.Comment) (int64, error)

func (m CommentRepositoryMock) SaveComment(id entities.Comment) (int64, error) {
	m.Mock.CallFunc("SaveComment", &map[string]interface{}{
		"id": id,
	})
	return SaveCommentMock(id)
}
