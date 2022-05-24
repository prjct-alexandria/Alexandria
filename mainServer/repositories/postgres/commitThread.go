package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
	"strconv"
)

type PgCommitThreadRepository struct {
	Db *sql.DB
}

func NewPgCommitThreadRepository(db *sql.DB) PgCommitThreadRepository {
	repo := PgCommitThreadRepository{db}
	err := repo.createThreadTable()
	if err != nil {
		return PgCommitThreadRepository{}
	}
	err = repo.createCommitThreadTable()
	if err != nil {
		return PgCommitThreadRepository{}
	}
	return repo
}

func (r PgCommitThreadRepository) CreateThread(aid string) (entities.Thread, error) {
	result, err := r.Db.Exec("INSERT INTO thread (articleId) " +
		"OUTPUT Inserted.threadId" +
		"VALUES ('" + aid + "')")
	if err != nil {
		return entities.Thread{}, fmt.Errorf("CreateThread: %v", err)
	}
	intAid, err1 := strconv.ParseInt(aid, 10, 64)
	if err1 != nil {
		return entities.Thread{}, err1
	}

	tid, err := result.LastInsertId()

	return entities.Thread{
		Id:        tid,
		ArticleId: intAid,
	}, err
}

func (r PgCommitThreadRepository) CreateCommitThread(thread models.CommitThreadNoId) (models.CommitThread, error) {
	// TODO
	return models.CommitThread{}, errors.New("")
}

func (r PgCommitThreadRepository) createThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS thread (
    	threadId int(64) NOT NULL AUTO_INCREMENT,
    	articleId int(64) NOT NULL,
    	PRIMARY KEY (threadId)
    )`)
	return err
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
