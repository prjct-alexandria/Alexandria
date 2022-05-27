package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type ReviewThreadService struct {
	ReviewThreadRepository interfaces.ReviewThreadRepository
}

func (serv ReviewThreadService) StartReviewThread(thread models.ThreadNoId, tid int64) (int64, error) {
	// create reviewThread
	id, err := serv.ReviewThreadRepository.CreateReviewThread(thread, strconv.FormatInt(tid, 10))
	if err != nil {
		return 0, err
	}

	return id, err
}
