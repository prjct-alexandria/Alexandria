package interfaces

import "mainServer/models"

type CommitSelectionThreadRepository interface {
	GetCommitSelectionThreads(aid int64, cid string) ([]models.SelectionThread, error)
	CreateCommitSelectionThread(cid string, tid int64, section string) (int64, error)
}
