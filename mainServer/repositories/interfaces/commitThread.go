package interfaces

import "mainServer/models"

type CommitThreadRepository interface {
	CreateCommitThread(thread models.Thread, tid int64) (int64, error)
}
