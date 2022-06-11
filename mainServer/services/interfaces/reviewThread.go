package interfaces

import "mainServer/models"

type ReviewThreadService interface {
	StartReviewThread(thread models.Thread, tid int64) (int64, error)
}
