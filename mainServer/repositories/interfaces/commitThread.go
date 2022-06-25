package interfaces

import "mainServer/models"

type CommitThreadRepository interface {
	GetCommitThreads(aid int64, cid string) ([]models.Thread, error)
	CreateCommitThread(cid string, tid int64) (int64, error)
}
