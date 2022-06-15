package services

import (
	"mainServer/repositories/interfaces"
)

type RequestThreadService struct {
	RequestThreadRepository interfaces.RequestThreadRepository
}

func (serv RequestThreadService) StartRequestThread(rid int64, tid int64) (int64, error) {
	// create requestThread
	id, err := serv.RequestThreadRepository.CreateRequestThread(rid, tid)
	if err != nil {
		return 0, err
	}

	return id, err
}
