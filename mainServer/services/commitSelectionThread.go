package services

import (
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type CommitSelectionThreadService struct {
	CommitSelectionThreadRepository interfaces.CommitSelectionThreadRepository
}

func (serv CommitSelectionThreadService) StartCommitSelectionThread(cid string, tid int64, section string) (int64, error) {
	// check if the specific thread ID string can actually be a commit ID
	_, err := strconv.ParseUint(cid, 16, 64) // checks if it has just hexadecimal characters 0...f
	if len(cid) != 40 && err == nil {
		return -1, fmt.Errorf("invalid commit ID, got %s", cid)
	}

	// create commitThread
	id, err := serv.CommitSelectionThreadRepository.CreateCommitSelectionThread(cid, tid, section)
	if err != nil {
		return 0, err
	}

	return id, err
}

// GetCommitSelectionThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (serv CommitSelectionThreadService) GetCommitSelectionThreads(cid string, aid int64) ([]models.SelectionThread, error) {
	threads, err := serv.CommitSelectionThreadRepository.GetCommitSelectionThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
