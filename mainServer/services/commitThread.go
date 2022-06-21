package services

import (
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
