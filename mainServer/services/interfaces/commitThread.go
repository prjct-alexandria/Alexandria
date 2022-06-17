package interfaces

import "mainServer/models"

type CommitThreadService interface {
	StartCommitThread(cid string, tid int64) (int64, error)
	GetCommitThreads(aid int64, cid string) ([]models.Thread, error)
}
