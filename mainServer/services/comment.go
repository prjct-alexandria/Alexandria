package services

import (
	"fmt"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type CommentService struct {
	CommentRepository interfaces.CommentRepository
}

// SaveComment saves list of comments to the db
// returns the id's of the saved comments
func (serv CommentService) SaveComment(comment []entities.Comment, tid int64) ([]int64, error) {
	// TODO: check if user is authenticated
	var err error
	var ids []int64

	for _, v := range comment {
		id, err := serv.CommentRepository.SaveComment(
			models.CommentNoId{
				AuthorId:     v.AuthorId,
				ThreadId:     tid,
				Content:      v.Content,
				CreationDate: v.CreationDate,
			})
		if err != nil {
			fmt.Println(err)
			return []int64{0}, err
		}
		ids = append(ids, id)
	}
	return ids, err

}
