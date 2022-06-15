package postgres

import (
	"database/sql"
	"fmt"
)

type PgCommitThreadRepository struct {
	Db *sql.DB
}

func NewPgCommitThreadRepository(db *sql.DB) PgCommitThreadRepository {
	repo := PgCommitThreadRepository{db}
	err := repo.createCommitThreadTable()
	if err != nil {
		return PgCommitThreadRepository{}
	}
	return repo
}

func (r PgCommitThreadRepository) createCommitThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitThread (
    	commitThreadId SERIAL,
    	commitId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (commitThreadId)
    )`)
	return err
}

func (r PgCommitThreadRepository) CreateCommitThread(cid string, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO committhread (committhreadid, commitid, threadid) 
	VALUES (DEFAULT, $1, $2) RETURNING committhreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(cid, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, nil
}
