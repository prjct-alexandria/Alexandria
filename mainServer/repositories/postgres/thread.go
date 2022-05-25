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
		fmt.Println(err)
		return PgThreadRepository{}
	}
	//err = repo.createCommitThreadTable()
	//if err != nil {
	//	return PgCommitThreadRepository{}
	//}
	return repo
}

func (r PgThreadRepository) CreateThread(aid string) (int64, error) {
	var tid int64
	_, err := r.Db.Exec("INSERT INTO thread (articleId) " +
		"VALUES ('" + aid + "')" +
		"RETURNING threadId")

	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("CreateThread: %v", err)
	}

	return tid, err
}

//func (r PgCommitThreadRepository) CreateCommitThread(thread models.CommitThreadNoId) (models.CommitThread, error) {
//	// TODO
//	return models.CommitThread{}, errors.New("")
//}

func (r PgThreadRepository) createThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS thread (
    	threadId SERIAL,
    	articleId BIGINT NOT NULL,
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
