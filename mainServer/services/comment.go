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
func (serv CommentService) SaveComment(comment entities.Comment, tid int64) (int64, error) {
	// TODO: check if user is authenticated

	var err error
	id, err := serv.CommentRepository.SaveComment(
		entities.Comment{
			AuthorId:     comment.AuthorId,
			ThreadId:     tid,
			Content:      comment.Content,
			CreationDate: comment.CreationDate,
		})
	if err != nil {
		fmt.Println(err)
		return int64(0), err
	}
	return id, err
}
