package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
)

type PgCommitThreadRepository struct {
	Db *sql.DB
}

func (r PgCommitThreadRepository) GetCommitThreads(aid int64, cid int64) (interface{}, error) {
	// get all thread id's that belong to the article
	// get all threads that belong to the thread id's
	// construct list of threads
	return nil, nil
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

func (r PgCommitThreadRepository) CreateCommitThread(thread models.Thread, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare("INSERT INTO committhread (committhreadid, commitid, threadid) VALUES (DEFAULT, $1, $2) RETURNING committhreadid")
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(thread.SpecificId, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, nil
}
