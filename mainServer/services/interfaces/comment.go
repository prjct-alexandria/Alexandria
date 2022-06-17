package interfaces

import "mainServer/entities"

type CommentService interface {
	// SaveComment saves list of comments to the db
	SaveComment(comment entities.Comment, tid int64) (int64, error)
}
