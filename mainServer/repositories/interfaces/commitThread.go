package interfaces

import (
	"mainServer/entities"
	"mainServer/models"
)

type CommitThreadRepository interface {
	CreateThread(aid string) (entities.Thread, error)
	CreateCommitThread(thread models.CommitThreadNoId) (models.CommitThread, error)
}
