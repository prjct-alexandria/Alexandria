package services

import (
	"errors"
	"fmt"
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type RequestThreadService struct {
	RequestThreadRepository interfaces.RequestThreadRepository
	VersionRepository       interfaces.VersionRepository
	RequestRepository       interfaces.RequestRepository
}

func (serv RequestThreadService) StartRequestThread(rid int64, tid int64, loggedInAs string) (int64, error) {
	req, err := serv.RequestRepository.GetRequest(rid)

	//Check if user is allowed to create request
	isSourceOwner, err := serv.VersionRepository.CheckIfOwner(req.SourceVersionID, loggedInAs)
	if err != nil {
		return -1, errors.New("could not verify version ownership")
	}
	isTargetOwner, err := serv.VersionRepository.CheckIfOwner(req.TargetVersionID, loggedInAs)
	if err != nil {
		return -1, errors.New("could not verify version ownership")
	}
	// TODO make endpoint return 403 Forbidden after this error
	if !(isSourceOwner || isTargetOwner) {
		return -1,
			errors.New(fmt.Sprintf(`request creation forbidden: %v does not own source or target version`, loggedInAs))
	}

	// create requestThread
	id, err := serv.RequestThreadRepository.CreateRequestThread(rid, tid)
	if err != nil {
		return -1, err
	}

	return id, err
}

//GetRequestThreads gets the request comment threads from the database, using the article id (aid) and request id (rid)
func (serv RequestThreadService) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	threads, err := serv.RequestThreadRepository.GetRequestThreads(aid, rid)
	if err != nil {
		return nil, err
	}

	return threads, err
}
