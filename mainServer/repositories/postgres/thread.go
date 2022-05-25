package postgres

import (
	"database/sql"
	"fmt"
)

type PgThreadRepository struct {
	Db *sql.DB
}

func NewPgThreadRepository(db *sql.DB) PgThreadRepository {
	repo := PgThreadRepository{db}
	err := repo.createThreadTable()
	if err != nil {
		return PgThreadRepository{}
	}
	//err = repo.createCommitThreadTable()
	//if err != nil {
	//	return PgCommitThreadRepository{}
	//}
	return repo
}

func (r PgThreadRepository) CreateThread(aid string) (int64, error) {
	result, err := r.Db.Exec("INSERT INTO thread (articleId) " +
		"OUTPUT Inserted.threadId" +
		"VALUES ('" + aid + "')")
	if err != nil {
		return 0, fmt.Errorf("CreateThread: %v", err)
	}
	tid, err := result.LastInsertId()

	return tid, err
}

//func (r PgCommitThreadRepository) CreateCommitThread(thread models.CommitThreadNoId) (models.CommitThread, error) {
//	// TODO
//	return models.CommitThread{}, errors.New("")
//}

func (r PgThreadRepository) createThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS thread (
    	threadId int(64) NOT NULL AUTO_INCREMENT,
    	articleId int(64) NOT NULL,
    	PRIMARY KEY (threadId)
    )`)
	return err
}

//func (r PgCommitThreadRepository) createCommitThreadTable() error {
//	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitThread (
//    	commitThreadId int(64) NOT NULL AUTO_INCREMENT,
//    	commitId int(64) NOT NULL,
//	    threadId int(64) NOT NULL,
//    	PRIMARY KEY (commitThreadId)
//    )`)
//	return err
//}
