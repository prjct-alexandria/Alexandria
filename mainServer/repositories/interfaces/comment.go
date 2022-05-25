package interfaces

import "mainServer/models"

type CommentRepository interface {
	SaveComment(id models.CommentNoId) (int64, error)
}
