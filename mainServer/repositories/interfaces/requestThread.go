package interfaces

import "mainServer/models"

type RequestThreadRepository interface {
	CreateRequestThread(thread models.Thread, tid int64) (int64, error)
	GetRequestThreads(aid int64, rid int64) ([]models.Thread, error)
}
