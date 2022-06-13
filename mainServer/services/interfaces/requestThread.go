package interfaces

import "mainServer/models"

type RequestThreadService interface {
	StartRequestThread(thread models.Thread, tid int64) (int64, error)
	GetRequestThreads(aid int64, rid int64) ([]models.Thread, error)
}
