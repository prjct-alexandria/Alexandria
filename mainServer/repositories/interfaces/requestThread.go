package interfaces

import "mainServer/models"

type RequestThreadRepository interface {
	CreateRequestThread(thread models.Thread, tid int64) (int64, error)
}
