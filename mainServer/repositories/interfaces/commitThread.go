package interfaces

import "mainServer/models"

type CommitThreadRepository interface {
	CreateCommitThread(thread models.Thread, tid int64) (int64, error)
	GetCommitThreads(aid int64, cid int64) (interface{}, error)
}
