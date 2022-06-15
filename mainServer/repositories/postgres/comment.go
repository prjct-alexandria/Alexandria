package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"mainServer/utils/clock"
)

type PgCommentRepository struct {
	Db    *sql.DB
	Clock clock.Clock
}

func NewPgCommentRepository(db *sql.DB) PgCommentRepository {
	repo := PgCommentRepository{Db: db, Clock: clock.RealClock{}}
	err := repo.createCommentTable()
	if err != nil {
		fmt.Println(err)
		return PgCommentRepository{}
	}
	return repo
}

func (r PgCommentRepository) createCommentTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS comment (
    	commentId SERIAL,
    	authorId varchar(64) NOT NULL,
    	content varchar(64) NOT NULL,
        threadId BIGINT NOT NULL,
        creationDate BIGINT NOT NULL,
    	PRIMARY KEY (commentId)
    )`)
	return err
}

func (r PgCommentRepository) SaveComment(comment entities.Comment) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO comment (commentid, authorid, content, threadid, creationdate) 
		VALUES (DEFAULT, $1, $2, $3, $4) RETURNING commentid`)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("SaveComment: %v", err)
	}
	row := stmt.QueryRow(comment.AuthorId, comment.Content, comment.ThreadId, r.Clock.Now().Unix())

	var id int64
	err = row.Scan(&id)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("SaveComment: %v", err)
	}

	return id, nil
}
