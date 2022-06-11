package interfaces

import "mainServer/models"

type CommitThreadService interface {
	StartCommitThread(thread models.Thread, tid int64) (int64, error)
	GetCommitThreads(aid int64, cid int64) ([]models.Thread, error)
}
