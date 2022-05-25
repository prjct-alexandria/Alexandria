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
    	commitThreadId int(64) NOT NULL AUTO_INCREMENT,
    	commitId int(64) NOT NULL,
        threadId int(64) NOT NULL,
    	PRIMARY KEY (commitThreadId)
    )`)
	return err
}

func (r PgCommitThreadRepository) CreateCommitThread(thread models.ThreadNoId, tid string) (int64, error) {
	result, err := r.Db.Exec("INSERT INTO commitThread (commitId, threadId) " +
		"OUTPUT Inserted.threadId" +
		"VALUES ('" +
		strconv.FormatInt(thread.CommitId, 10) + "', '" +
		tid + "')")
	if err != nil {
		return 0, fmt.Errorf("CreateThread: %v", err)
	}
	id, err := result.LastInsertId()

	return id, err
}
