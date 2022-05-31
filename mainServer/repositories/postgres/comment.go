package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
	"strconv"
)

type PgCommentRepository struct {
	Db *sql.DB
}

func NewPgCommentRepository(db *sql.DB) PgCommentRepository {
	repo := PgCommentRepository{db}
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
        creationDate varchar(64) NOT NULL,
    	PRIMARY KEY (commentId)
    )`)
	return err
}

func (r PgCommentRepository) SaveComment(comment models.Comment) (int64, error) {
	row := r.Db.QueryRow("INSERT INTO comment (authorId, content, threadId, creationDate) " +
		"VALUES ('" + comment.AuthorId + "', '" +
		comment.Content + "', '" +
		strconv.FormatInt(comment.ThreadId, 10) + "', '" +
		comment.CreationDate +
		"')" +
		"RETURNING commentId")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("SaveComment: %v", err)
	}

	return id, err
}
