package interfaces

import "mainServer/models"

type CommitThreadRepository interface {
	CreateCommitThread(thread models.ThreadNoId, tid string) (int64, error)
}
