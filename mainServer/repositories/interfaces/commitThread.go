package interfaces

import "mainServer/models"

type CommitThreadRepository interface {
	CreateCommitThread(thread models.Thread, tid string) (int64, error)
}
