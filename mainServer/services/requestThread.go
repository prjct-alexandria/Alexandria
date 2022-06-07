package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type RequestThreadService struct {
	RequestThreadRepository interfaces.RequestThreadRepository
}

func (serv RequestThreadService) StartRequestThread(thread models.Thread, tid int64) (int64, error) {
	// create requestThread
	id, err := serv.RequestThreadRepository.CreateRequestThread(thread, tid)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (serv RequestThreadService) GetRequestThreads(aid int64, cid int64) ([]models.Thread, error) {
	threads, err := serv.RequestThreadRepository.GetRequestThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
