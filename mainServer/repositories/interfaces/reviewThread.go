package interfaces

import "mainServer/models"

type ReviewThreadRepository interface {
	CreateReviewThread(thread models.ThreadNoId, tid string) (int64, error)
}
