package interfaces

type CommitThreadRepository interface {
	CreateCommitThread(cid string, tid int64) (int64, error)
}
