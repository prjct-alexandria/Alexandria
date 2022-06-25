package interfaces

import "mainServer/models"

type CommitSelectionThreadService interface {
	StartCommitSelectionThread(cid string, tid int64, section string) (int64, error)
	GetCommitSelectionThreads(sid string, aid int64) ([]models.SelectionThread, error)
}
