package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type RequestThreadService struct {
	RequestThreadRepository interfaces.RequestThreadRepository
}

func (serv RequestThreadService) StartRequestThread(thread models.ThreadNoId, tid int64) (int64, error) {
	// create requestThread
	id, err := serv.RequestThreadRepository.CreateRequestThread(thread, strconv.FormatInt(tid, 10))
	if err != nil {
		return 0, err
	}

	return id, err
}
