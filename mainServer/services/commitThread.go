package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type CommitThreadService struct {
	CommitThreadRepository interfaces.CommitThreadRepository
}

func (serv CommitThreadService) StartCommitThread(cid string, tid int64) (int64, error) {

	// create commitThread
	id, err := serv.CommitThreadRepository.CreateCommitThread(cid, tid)
	if err != nil {
		return 0, err
	}

	return id, err
}

// GetCommitThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (serv CommitThreadService) GetCommitThreads(aid int64, cid string) ([]models.Thread, error) {
	threads, err := serv.CommitThreadRepository.GetCommitThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
