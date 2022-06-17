package interfaces

import "mainServer/models"

type CommitSectionThreadService interface {
	StartCommitSectionThread(cid string, tid int64, section string) (int64, error)
	GetCommitSectionThreads(sid string, tid int64) ([]models.Thread, error)
}
