package interfaces

import "mainServer/models"

type RequestThreadRepository interface {
	CreateRequestThread(thread models.Thread, tid string) (int64, error)
}
