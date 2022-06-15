package interfaces

import "mainServer/models"

type ThreadService interface {
	// StartThread creates thread entity in db
	StartThread(thread models.Thread, aid int64) (int64, error)
}
