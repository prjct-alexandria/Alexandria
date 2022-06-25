package interfaces

import (
	"mainServer/entities"
)

type CommentRepository interface {
	SaveComment(id entities.Comment) (int64, error)
}
