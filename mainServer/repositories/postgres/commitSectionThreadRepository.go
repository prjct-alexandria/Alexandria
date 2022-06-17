package postgres

import (
	"database/sql"
	"mainServer/models"
)

type PgCommitSectionThreadRepository struct {
	Db *sql.DB
}

func NewPgCommitSectionThreadRepository(db *sql.DB) PgCommitThreadRepository {
	repo := PgCommitThreadRepository{db}
	err := repo.createCommitThreadTable()
	if err != nil {
		return PgCommitThreadRepository{}
	}
	return repo
}

func (r PgCommitThreadRepository) createCommitSectionThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitSectionThread (
    	commitSectionThreadId SERIAL,
    	commitId NCHAR(40) NOT NULL,
        threadId BIGINT NOT NULL,
        content  NCHAR(255) NOT NULL,
    	PRIMARY KEY (commitSectionThreadId)
    )`)
	return err
}

func (r PgCommitSectionThreadRepository) GetCommitSectionThreads(aid int64, cid string) ([]models.SectionThread, error) {

	return nil, nil
}

func (r PgCommitSectionThreadRepository) CreateCommitSectionThread(cid string, tid int64, section string) (int64, error) {
	return 0, nil
}
