package services

import (
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type CommitThreadService struct {
	CommitThreadRepository interfaces.CommitThreadRepository
}

func (serv CommitThreadService) StartCommitThread(cid string, tid int64) (int64, error) {

	// check if the specific thread ID string can actually be a commit ID
	//_, err := strconv.ParseUint(cid, 16, 64) // checks if it has just hexadecimal characters 0...f
	//if len(cid) != 40 && err == nil {
	//	return -1, fmt.Errorf("invalid commit ID, got %s", cid)
	//}

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
