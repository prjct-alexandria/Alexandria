package services

import (
	"errors"
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"strconv"
)

type ThreadService struct {
	ThreadRepository interfaces.ThreadRepository
}

// StartThread creates thread entity in db
// returns thread id
func (serv ThreadService) StartThread(thread models.Thread, aid string, sid string) (int64, error) {
	// TODO: check if user is authenticated

	// check model has same aid and cid as params
	intSid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return 0, err
	}
	intAid, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return 0, err
	}
	if thread.SpecificId != intSid || thread.ArticleId != intAid {
		return 0, errors.New("parameters in url not equal to the thread object")
	}

	// create thread
	tid, err := serv.ThreadRepository.CreateThread(aid)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// return threadid and possible error
	return tid, err
}
