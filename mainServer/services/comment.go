package services

import (
	"fmt"
	"mainServer/entities"
	"mainServer/repositories/interfaces"
)

type CommentService struct {
	CommentRepository interfaces.CommentRepository
}

// SaveComment saves list of comments to the db
// returns the id's of the saved comments
func (serv CommentService) SaveComment(comment entities.Comment, tid int64, loggedInAs string) (int64, error) {
	if loggedInAs != comment.AuthorId {
		return int64(-1),
			fmt.Errorf("provided author %v does not match logged-in author %v", comment.AuthorId, loggedInAs)
	}

	id, err := serv.CommentRepository.SaveComment(
		entities.Comment{
			AuthorId:     comment.AuthorId,
			ThreadId:     tid,
			Content:      comment.Content,
			CreationDate: comment.CreationDate,
		})
	if err != nil {
		return int64(-1), err
	}
	return id, nil
}
