package services

import (
	"errors"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type ThreadService struct {
	ThreadRepository interfaces.ThreadRepository
}

// creates thread entity in db
// returns thread id
// TODO: divide in two parts: create thread and create specific thread (like review/request/commit)
func (serv ThreadService) StartThread(thread models.ThreadNoId, aid string, cid string) (int64, error) {
	// TODO: check if user is authenticated

	// check model has same aid and cid as params
	intCid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return 0, err
	}
	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return 0, err
	}
	if thread.CommitId != intCid || thread.ArticleId != intAid {
		return 0, errors.New("parameters in url not equal to the thread object")
	}

	// create thread
	tid, err := serv.ThreadRepository.CreateThread(aid)
	if err != nil {
		return 0, err
	}

	// return threadid and possible error
	return tid, err
}
