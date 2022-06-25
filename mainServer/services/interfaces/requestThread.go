package interfaces

import "mainServer/models"

type RequestThreadService interface {
	StartRequestThread(rid int64, tid int64, loggedInAs string) (int64, error)
	GetRequestThreads(aid int64, rid int64) ([]models.Thread, error)
}
