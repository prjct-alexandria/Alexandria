package interfaces

import "mainServer/models"

type CommentRepository interface {
	SaveComment(id models.Comment) (int64, error)
}
