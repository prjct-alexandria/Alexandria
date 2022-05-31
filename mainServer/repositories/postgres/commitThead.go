package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
	"strconv"
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

func (r PgCommitThreadRepository) CreateCommitThread(thread models.Thread, tid string) (int64, error) {
	row := r.Db.QueryRow("INSERT INTO commitThread (commitId, threadId) " +
		"VALUES ('" +
		strconv.FormatInt(thread.SpecificId, 10) + "', '" +
		tid + "')" +
		"RETURNING commitThreadId")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateThread: %v", err)
	}

	return id, err
}
