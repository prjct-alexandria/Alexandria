package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type CommitThreadService struct {
	CommitThreadRepository interfaces.CommitThreadRepository
}

func (serv CommitThreadService) StartCommitThread(thread models.ThreadNoId, tid int64) (int64, error) {
	// create commitThread
	id, err := serv.CommitThreadRepository.CreateCommitThread(thread, strconv.FormatInt(tid, 10))
	if err != nil {
		return 0, err
	}

	return id, err
}
