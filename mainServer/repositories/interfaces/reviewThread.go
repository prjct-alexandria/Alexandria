package interfaces

import "mainServer/models"

type ReviewThreadRepository interface {
	CreateReviewThread(thread models.Thread, tid int64) (int64, error)
}
