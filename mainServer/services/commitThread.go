package services

import (
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	gitUtils "mainServer/utils/git"
)

type CommitThreadService struct {
	CommitThreadRepository interfaces.CommitThreadRepository
}

func (serv CommitThreadService) StartCommitThread(cid string, tid int64) (int64, error) {
	if !gitUtils.IsCommitHash(cid) {
		return int64(-1), fmt.Errorf("invalid commit ID, got %s", cid)
	}

	// create commitThread
	id, err := serv.CommitThreadRepository.CreateCommitThread(cid, tid)
	if err != nil {
		return int64(-1), err
	}

	return id, err
}

// GetCommitThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (serv CommitThreadService) GetCommitThreads(aid int64, cid string) ([]models.Thread, error) {
	if !gitUtils.IsCommitHash(cid) {
		return nil, fmt.Errorf("invalid commit ID, got %s", cid)
	}

	threads, err := serv.CommitThreadRepository.GetCommitThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
