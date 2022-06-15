package postgres

import (
	"database/sql"
	"fmt"
)

type PgRequestThreadRepository struct {
	Db *sql.DB
}

func NewPgRequestThreadRepository(db *sql.DB) PgRequestThreadRepository {
	repo := PgRequestThreadRepository{db}
	err := repo.createRequestThreadTable()
	if err != nil {
		return PgRequestThreadRepository{}
	}
	return repo
}

func (r PgRequestThreadRepository) createRequestThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS RequestThread (
    	RequestThreadId SERIAL,
    	RequestId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (RequestThreadId)
    )`)
	return err
}

func (r PgRequestThreadRepository) CreateRequestThread(rid int64, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO requestthread (requestthreadid, requestid, threadid)
	VALUES (DEFAULT, $1, $2) RETURNING requestthreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(rid, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, err
}
