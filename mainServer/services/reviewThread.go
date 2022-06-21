package services

import (
	"mainServer/repositories/interfaces"
)

type ReviewThreadService struct {
	ReviewThreadRepository interfaces.ReviewThreadRepository
}

func (serv ReviewThreadService) StartReviewThread(rid int64, tid int64) (int64, error) {
	// create reviewThread
	id, err := serv.ReviewThreadRepository.CreateReviewThread(rid, tid)
	if err != nil {
		return 0, err
	}

	return id, err
}
