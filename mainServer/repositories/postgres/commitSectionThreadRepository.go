package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
)

type PgCommitSectionThreadRepository struct {
	Db *sql.DB
}

func NewPgCommitSectionThreadRepository(db *sql.DB) PgCommitSectionThreadRepository {
	repo := PgCommitSectionThreadRepository{db}
	err := repo.createCommitSectionThreadTable()
	if err != nil {
		return PgCommitSectionThreadRepository{}
	}
	return repo
}

func (r PgCommitSectionThreadRepository) createCommitSectionThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitSectionThread (
    	commitSectionThreadId SERIAL,
    	commitId NCHAR(40) NOT NULL,
        threadId BIGINT NOT NULL,
        section  NCHAR(255) NOT NULL,
    	PRIMARY KEY (commitSectionThreadId)
    )`)
	return err
}

func (r PgCommitSectionThreadRepository) CreateCommitSectionThread(cid string, tid int64, section string) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO commitsectionthread (commitsectionthreadid, commitid, threadid, section)
	VALUES (DEFAULT, $1, $2, $3) RETURNING commitsectionthreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(cid, tid, section)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, nil
}

func (r PgCommitSectionThreadRepository) GetCommitSectionThreads(aid int64, cid string) ([]models.SectionThread, error) {
	return nil, nil
}
