package services

import (
	"errors"
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type ThreadService struct {
	ThreadRepository interfaces.ThreadRepository
}

// StartThread creates thread entity in db
// returns thread id
func (serv ThreadService) StartThread(thread models.Thread, aid int64) (int64, error) {
	// check model has same aid as params
	if thread.ArticleId != aid {
		return -1, errors.New("parameters in url not equal to the thread object")
	}

	// create thread
	tid, err := serv.ThreadRepository.CreateThread(aid)
	if err != nil {
		return -1, err
	}

	// return threadid
	return tid, nil
}
