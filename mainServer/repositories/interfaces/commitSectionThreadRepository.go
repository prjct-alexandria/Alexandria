package interfaces

import "mainServer/models"

type CommitSectionThreadRepository interface {
	GetCommitSectionThreads(aid int64, cid string) ([]models.SectionThread, error)
	CreateCommitSectionThread(cid string, tid int64, section string) (int64, error)
}
