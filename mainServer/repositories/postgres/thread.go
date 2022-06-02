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

	return repo
}

func (r PgThreadRepository) CreateThread(aid int64) (int64, error) {
	stmt, err := r.Db.Prepare("INSERT INTO thread (threadid, articleid) VALUES (DEFAULT, $1) RETURNING threadid")
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(aid)

	//row := r.Db.QueryRow("INSERT INTO thread (articleId) " +
	//	"VALUES ('" + aid + "')" +
	//	"RETURNING threadId")
	var tid int64
	err = row.Scan(&tid)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return tid, err
}

func (r PgThreadRepository) createThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS thread (
    	threadId SERIAL,
    	articleId BIGINT NOT NULL,
    	PRIMARY KEY (threadId)
    )`)
	return err
}
