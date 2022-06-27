package services

import (
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	gitUtils "mainServer/utils/git"
)

type CommitSelectionThreadService struct {
	CommitSelectionThreadRepository interfaces.CommitSelectionThreadRepository
}

func (serv CommitSelectionThreadService) StartCommitSelectionThread(cid string, tid int64, section string) (int64, error) {
	// check if the specific commit ID string can actually be a commit ID
	if !gitUtils.IsCommitHash(cid) {
		return -1, fmt.Errorf("invalid commit ID, got %s", cid)
	}

	// create commitThread
	id, err := serv.CommitSelectionThreadRepository.CreateCommitSelectionThread(cid, tid, section)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// GetCommitSelectionThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (serv CommitSelectionThreadService) GetCommitSelectionThreads(cid string, aid int64) ([]models.SelectionThread, error) {
	if !gitUtils.IsCommitHash(cid) {
		return nil, fmt.Errorf("invalid commit ID, got %s", cid)
	}

	threads, err := serv.CommitSelectionThreadRepository.GetCommitSelectionThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, nil
}
