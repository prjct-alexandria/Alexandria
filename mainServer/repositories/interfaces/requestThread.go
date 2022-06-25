package interfaces

import "mainServer/models"

type RequestThreadRepository interface {
	GetRequestThreads(aid int64, rid int64) ([]models.Thread, error)
	CreateRequestThread(rid int64, tid int64) (int64, error)
}
