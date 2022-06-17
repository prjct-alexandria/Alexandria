package services

import (
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type CommitSectionThreadService struct {
	CommitSectionThreadRepository interfaces.CommitSectionThreadRepository
}

func (serv CommitSectionThreadService) StartCommitSectionThread(cid string, tid int64, section string) (int64, error) {
	// check if the specific thread ID string can actually be a commit ID
	_, err := strconv.ParseUint(cid, 16, 64) // checks if it has just hexadecimal characters 0...f
	if len(cid) != 40 && err == nil {
		return -1, fmt.Errorf("invalid commit ID, got %s", cid)
	}

	// create commitThread
	id, err := serv.CommitSectionThreadRepository.CreateCommitSectionThread(cid, tid, section)
	if err != nil {
		return 0, err
	}

	return id, err
}

// GetCommitSectionThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (serv CommitSectionThreadService) GetCommitSectionThreads(cid string, aid int64) ([]models.SectionThread, error) {
	threads, err := serv.CommitSectionThreadRepository.GetCommitSectionThreads(aid, cid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
