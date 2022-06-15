package services

import (
	"mainServer/models"
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

//GetRequestThreads gets the request comment threads from the database, using the article id (aid) and request id (rid)
func (serv RequestThreadService) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	threads, err := serv.RequestThreadRepository.GetRequestThreads(aid, rid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
