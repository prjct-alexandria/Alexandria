package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type ReviewThreadService struct {
	ReviewThreadRepository interfaces.ReviewThreadRepository
}

func (serv ReviewThreadService) StartReviewThread(thread models.Thread, tid int64) (int64, error) {
	// create reviewThread
	id, err := serv.ReviewThreadRepository.CreateReviewThread(thread, tid)
	if err != nil {
		return 0, err
	}

	return id, err
}
