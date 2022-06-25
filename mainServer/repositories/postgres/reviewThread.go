package postgres

import (
	"database/sql"
	"fmt"
)

type PgReviewThreadRepository struct {
	Db *sql.DB
}

func NewPgReviewThreadRepository(db *sql.DB) PgReviewThreadRepository {
	repo := PgReviewThreadRepository{db}
	err := repo.createReviewThreadTable()
	if err != nil {
		return PgReviewThreadRepository{}
	}
	return repo
}

func (r PgReviewThreadRepository) createReviewThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS ReviewThread (
    	ReviewThreadId SERIAL,
    	ReviewId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (ReviewThreadId)
    )`)
	return err
}

func (r PgReviewThreadRepository) CreateReviewThread(rid int64, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO reviewthread (reviewthreadid, reviewid, threadid)
		VALUES (DEFAULT, $1, $2) RETURNING reviewthreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateReviewThread: %v", err)
	}
	row := stmt.QueryRow(rid, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateReviewThread: %v", err)
	}

	return id, err
}
