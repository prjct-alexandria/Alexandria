package interfaces

import "mainServer/models"

type ReviewThreadRepository interface {
	CreateReviewThread(thread models.Thread, tid string) (int64, error)
}
