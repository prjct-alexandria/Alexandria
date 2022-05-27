package interfaces

import "mainServer/models"

type RequestThreadRepository interface {
	CreateRequestThread(thread models.ThreadNoId, tid string) (int64, error)
}
